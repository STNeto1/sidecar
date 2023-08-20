package pkg

import (
	"encoding/json"
	"fmt"
	"os"
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
