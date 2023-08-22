package pkg

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/rs/zerolog/log"
)

const folderSuffix string = ".sidecar"

type ProfileData = map[string]string
type Profile = map[string]ProfileData

func CreateProfile(profileName string) error {
	if profileName == "" {
		return errors.New("no profile was specified")
	}

	profiles, err := getContext()
	if err != nil {
		return errors.New("failed to get profiles")
	}

	if hasProfile(profiles, profileName) {
		return errors.New("profile already exists")
	}

	profilePath, err := getProfilePath(profileName)
	if err != nil {
		return errors.New("failed to get profile path")
	}

	var profileData ProfileData
	jsonData, err := json.Marshal(profileData)
	if err != nil {
		return errors.New("failed to serialize new data")
	}

	if err := os.WriteFile(profilePath, jsonData, 0666); err != nil {
		return errors.New("failed to write to file")
	}

	return nil
}

func ListProfiles() (Profile, error) {
	profiles, err := getContext()
	if err != nil {
		return nil, errors.New("failed to get profiles")
	}

	return profiles, nil
}

func ShowProfile(profileName string) (ProfileData, error) {
	if profileName == "" {
		return nil, errors.New("no profile was specified")
	}

	profiles, err := getContext()
	if err != nil {
		return nil, errors.New("failed to get profiles")
	}

	if !hasProfile(profiles, profileName) {
		return nil, errors.New("profile does not exists")
	}

	data := profiles[profileName]
	return data, nil
}

func AddToProfile(profileName string, entries ...string) error {
	if profileName == "" {
		return errors.New("no profile was specified")
	}

	profiles, err := getContext()
	if err != nil {
		return errors.New("failed to get profiles")
	}

	if !hasProfile(profiles, profileName) {
		return errors.New("profile does not exists")
	}

	profileData := profiles[profileName]

	for _, value := range entries {
		tokens := strings.Split(value, "=")
		if len(tokens) != 2 {
			continue
		}

		profileData[tokens[0]] = tokens[1]
	}

	jsonData, err := json.MarshalIndent(profileData, "", " ")
	if err != nil {
		return errors.New("failed to serialize new data")
	}

	profilePath, err := getProfilePath(profileName)
	if err != nil {
		return errors.New("failed to get profile path")
	}

	if err := os.WriteFile(profilePath, jsonData, 0666); err != nil {
		return errors.New("failed to write to file")
	}

	return nil
}

func RemoveFromProfile(profileName string, entries ...string) error {
	if profileName == "" {
		return errors.New("no profile was specified")
	}

	profiles, err := getContext()
	if err != nil {
		return errors.New("failed to get profiles")
	}

	if !hasProfile(profiles, profileName) {
		return errors.New("profile does not exists")
	}

	profileData := profiles[profileName]

	for _, value := range entries {
		delete(profileData, value)
	}

	jsonData, err := json.MarshalIndent(profileData, "", " ")
	if err != nil {
		return errors.New("failed to serialize new data")
	}

	profilePath, err := getProfilePath(profileName)
	if err != nil {
		return errors.New("failed to get profile path")
	}

	if err := os.WriteFile(profilePath, jsonData, 0666); err != nil {
		return errors.New("failed to write to file")
	}

	return nil
}

func DeleteProfile(profileName string) error {
	if profileName == "" {
		return errors.New("no profile was specified")
	}

	profiles, err := getContext()
	if err != nil {
		return errors.New("failed to get profiles")
	}

	if !hasProfile(profiles, profileName) {
		return errors.New("profile does not exists")
	}

	profilePath, err := getProfilePath(profileName)
	if err != nil {
		return errors.New("failed to get profile path")
	}

	if err := os.Remove(profilePath); err != nil {
		return errors.New("failed to delete profile")
	}

	return nil
}

func Execute(profileName string, command string) error {
	if profileName == "" {
		return errors.New("no profile was specified")
	}

	if command == "" {
		return errors.New("no command was specified")
	}

	profiles, err := getContext()
	if err != nil {
		return errors.New("failed to get profiles")
	}

	if !hasProfile(profiles, profileName) {
		return errors.New("profile does not exists")
	}

	tokens := strings.Split(command, " ")

	profileData := profiles[profileName]
	for key, value := range profileData {
		if err := os.Setenv(key, value); err != nil {
			return errors.New("error setting env")
		}
	}

	head, tail := getHeadTail[string](tokens)
	cmd := exec.Command(head, tail...)

	pipe, _ := cmd.StdoutPipe()

	cmd.Start()

	scanner := bufio.NewScanner(pipe)
	scanner.Split(bufio.ScanWords)
	for scanner.Scan() {
		m := scanner.Text()

		log.Info().Msg(m)
	}

	cmd.Wait()

	return nil
}

func getContext() (Profile, error) {
	dirname, err := os.UserHomeDir()

	if err != nil {
		return nil, err
	}

	rootFolder := path.Join(dirname, folderSuffix)

	entries, err := os.ReadDir(rootFolder)
	if err != nil {
		return nil, err
	}

	profiles := make(Profile)

	for _, e := range entries {
		if e.IsDir() {
			continue
		}

		fileBytes, err := os.ReadFile(path.Join(rootFolder, e.Name()))
		if err != nil {
			continue
		}

		var profileData ProfileData
		err = json.Unmarshal(fileBytes, &profileData)
		if err != nil {
			continue
		}

		tokens := strings.Split(e.Name(), ".")
		if len(tokens) < 1 {
			continue
		}

		profiles[tokens[0]] = profileData
	}

	return profiles, nil
}

func hasProfile(profile Profile, profileName string) bool {
	for key := range profile {

		if key == profileName {
			return true
		}
	}

	return false
}

func getProfilePath(profileName string) (string, error) {

	dirname, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	rootFolder := path.Join(dirname, folderSuffix)
	profileName = fmt.Sprintf("%s.json", profileName)
	profilePath := path.Join(rootFolder, profileName)

	return profilePath, nil
}

func getKey(profileName string) string {
	return fmt.Sprintf("%s.json", profileName)
}

// assumes that the list has at least one element
func getHeadTail[K any](data []K) (K, []K) {
	var head K
	var tail []K

	for idx, value := range data {
		if idx == 0 {
			head = value
		} else {
			tail = append(tail, value)
		}
	}

	return head, tail
}
