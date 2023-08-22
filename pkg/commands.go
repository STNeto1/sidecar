package pkg

import (
	"bufio"
	"encoding/json"
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
		return ErrNoProfileSpecified
	}

	profiles, err := getContext()
	if err != nil {
		return ErrFailedToGetProfiles
	}

	if hasProfile(profiles, profileName) {
		return ErrProfileExists
	}

	profilePath, err := getProfilePath(profileName)
	if err != nil {
		return ErrFailedToGetProfilePath
	}

	var profileData ProfileData
	jsonData, err := json.Marshal(profileData)
	if err != nil {
		return ErrFailedToSerializeData
	}

	if err := os.WriteFile(profilePath, jsonData, 0666); err != nil {
		return ErrFailedToWriteToFile
	}

	return nil
}

func ListProfiles() (Profile, error) {
	profiles, err := getContext()
	if err != nil {
		return nil, ErrFailedToGetProfiles
	}

	return profiles, nil
}

func ShowProfile(profileName string) (ProfileData, error) {
	if profileName == "" {
		return nil, ErrNoProfileSpecified
	}

	profiles, err := getContext()
	if err != nil {
		return nil, ErrFailedToGetProfiles
	}

	if !hasProfile(profiles, profileName) {
		return nil, ErrProfileDoesNotExist
	}

	data := profiles[profileName]
	return data, nil
}

func AddToProfile(profileName string, entries ...string) error {
	if profileName == "" {
		return ErrNoProfileSpecified
	}

	profiles, err := getContext()
	if err != nil {
		return ErrFailedToGetProfiles
	}

	if !hasProfile(profiles, profileName) {
		return ErrProfileDoesNotExist
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
		return ErrFailedToSerializeData
	}

	profilePath, err := getProfilePath(profileName)
	if err != nil {
		return ErrFailedToGetProfilePath
	}

	if err := os.WriteFile(profilePath, jsonData, 0666); err != nil {
		return ErrFailedToWriteToFile
	}

	return nil
}

func RemoveFromProfile(profileName string, entries ...string) error {
	if profileName == "" {
		return ErrNoProfileSpecified
	}

	profiles, err := getContext()
	if err != nil {
		return ErrFailedToGetProfiles
	}

	if !hasProfile(profiles, profileName) {
		return ErrProfileDoesNotExist
	}

	profileData := profiles[profileName]

	for _, value := range entries {
		delete(profileData, value)
	}

	jsonData, err := json.MarshalIndent(profileData, "", " ")
	if err != nil {
		return ErrFailedToSerializeData
	}

	profilePath, err := getProfilePath(profileName)
	if err != nil {
		return ErrFailedToGetProfilePath
	}

	if err := os.WriteFile(profilePath, jsonData, 0666); err != nil {
		return ErrFailedToWriteToFile
	}

	return nil
}

func DeleteProfile(profileName string) error {
	if profileName == "" {
		return ErrNoProfileSpecified
	}

	profiles, err := getContext()
	if err != nil {
		return ErrFailedToGetProfiles
	}

	if !hasProfile(profiles, profileName) {
		return ErrProfileDoesNotExist
	}

	profilePath, err := getProfilePath(profileName)
	if err != nil {
		return ErrFailedToGetProfilePath
	}

	if err := os.Remove(profilePath); err != nil {
		return ErrFailedToDeleteProfile
	}

	return nil
}

func Execute(profileName string, command string) error {
	if profileName == "" {
		return ErrNoProfileSpecified
	}

	if command == "" {
		return ErrNoCommandSpecified
	}

	profiles, err := getContext()
	if err != nil {
		return ErrFailedToGetProfiles
	}

	if !hasProfile(profiles, profileName) {
		return ErrProfileDoesNotExist
	}

	tokens := strings.Split(command, " ")

	profileData := profiles[profileName]
	for key, value := range profileData {
		if err := os.Setenv(key, value); err != nil {
			return ErrFailedToSetEnv
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
