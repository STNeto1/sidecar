package pkg

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"
)

const folderSuffix string = ".sidecar"

type ProfileData = map[string]string
type Profile = map[string]ProfileData

func CreateProfile(profileName string) {
	if profileName == "" {
		fmt.Println("-> no profile was specified")
		return
	}

	profiles, err := getContext()
	if err != nil {
		fmt.Println("-> failed to get profiles", err.Error())
		return
	}

	if hasProfile(profiles, profileName) {
		fmt.Println("-> profile already exists")
		return
	}

	profilePath, err := getProfilePath(profileName)
	if err != nil {
		fmt.Println("-> failed to get profile path -", err.Error())
		return
	}

	var profileData ProfileData
	jsonData, err := json.Marshal(profileData)
	if err != nil {
		fmt.Println("-> failed to serialize new data -", err.Error())
		return
	}

	if err := os.WriteFile(profilePath, jsonData, 0666); err != nil {
		fmt.Println("-> failed to write to file -", err.Error())
		return
	}

	fmt.Println("-> created new profile")
}

func ListProfiles() {
	profiles, err := getContext()
	if err != nil {
		fmt.Println("-> failed to get profiles -", err.Error())
		return
	}

	for key := range profiles {
		fmt.Println(fmt.Sprintf("-> %s", key))
	}
}

func ShowProfile(profileName string) {
	if profileName == "" {
		fmt.Println("-> no profile was specified")
		return
	}

	profiles, err := getContext()
	if err != nil {
		fmt.Println("-> failed to get profiles", err.Error())
		return
	}

	if !hasProfile(profiles, profileName) {
		fmt.Println("-> profile does not exists")
		return
	}

	data := profiles[profileName]
	fmt.Println(fmt.Sprintf("===== %s =====", profileName))
	for key, value := range data {
		fmt.Println(fmt.Sprintf("-> %s: %s", key, value))
	}
}

func AddToProfile(profileName string, entries ...string) {
	if profileName == "" {
		fmt.Println("-> no profile was specified")
		return
	}

	profiles, err := getContext()
	if err != nil {
		fmt.Println("-> failed to get profiles", err.Error())
		return
	}

	if !hasProfile(profiles, profileName) {
		fmt.Println("-> profile does not exists")
		return
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
		fmt.Println("-> failed to serialize new data -", err.Error())
		return
	}

	profilePath, err := getProfilePath(profileName)
	if err != nil {
		fmt.Println("-> failed to get profile path -", err.Error())
		return
	}

	if err := os.WriteFile(profilePath, jsonData, 0666); err != nil {
		fmt.Println("-> failed to write to file -", err.Error())
		return
	}

	fmt.Println("-> addeded profile")
}

func RemoveFromProfile(profileName string, entries ...string) {
	if profileName == "" {
		fmt.Println("-> no profile was specified")
		return
	}

	profiles, err := getContext()
	if err != nil {
		fmt.Println("-> failed to get profiles", err.Error())
		return
	}

	if !hasProfile(profiles, profileName) {
		fmt.Println("-> profile does not exists")
		return
	}

	profileData := profiles[profileName]

	for _, value := range entries {
		delete(profileData, value)
	}

	jsonData, err := json.MarshalIndent(profileData, "", " ")
	if err != nil {
		fmt.Println("-> failed to serialize new data -", err.Error())
		return
	}

	profilePath, err := getProfilePath(profileName)
	if err != nil {
		fmt.Println("-> failed to get profile path -", err.Error())
		return
	}

	if err := os.WriteFile(profilePath, jsonData, 0666); err != nil {
		fmt.Println("-> failed to write to file -", err.Error())
		return
	}

	fmt.Println("-> addeded profile")
}

func DeleteProfile(profileName string) {
	if profileName == "" {
		fmt.Println("-> no profile was specified")
		return
	}

	profiles, err := getContext()
	if err != nil {
		fmt.Println("-> failed to get profiles", err.Error())
		return
	}

	if !hasProfile(profiles, profileName) {
		fmt.Println("-> profile does not exists")
		return
	}

	profilePath, err := getProfilePath(profileName)
	if err != nil {
		fmt.Println("-> failed to get profile path -", err.Error())
		return
	}

	if err := os.Remove(profilePath); err != nil {
		fmt.Println("-> failed to delete profile -", err.Error())
		return
	}

	fmt.Println("-> deleted profile")

}

func Execute(profileName string, command string) {
	if profileName == "" {
		fmt.Println("-> no profile was specified")
		return
	}

	if command == "" {
		fmt.Println("-> no command was specified")
		return
	}

	profiles, err := getContext()
	if err != nil {
		fmt.Println("-> failed to get profiles", err.Error())
		return
	}

	if !hasProfile(profiles, profileName) {
		fmt.Println("-> profile does not exists")
		return
	}

	tokens := strings.Split(command, " ")

	profileData := profiles[profileName]
	for key, value := range profileData {
		if err := os.Setenv(key, value); err != nil {
			fmt.Println("-> error setting env", err)
			return
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

		fmt.Println(m)
	}

	cmd.Wait()
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
			fmt.Println("error reading file", e.Name(), " - ", err.Error())
			continue
		}

		var profileData ProfileData
		err = json.Unmarshal(fileBytes, &profileData)
		if err != nil {
			fmt.Println("error parsing file", e.Name(), " - ", err.Error())
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
