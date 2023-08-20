package pkg

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path"
	"strings"
)

const folderSuffix string = ".sidecar"

type ProfileData = map[string]string
type Profile = map[string]ProfileData

func ListProfiles() {
	dirname, err := os.UserHomeDir()
	if err != nil {
		log.Println("failed to get user home dir -", err.Error())
		return
	}

	rootFolder := path.Join(dirname, folderSuffix)

	entries, err := os.ReadDir(rootFolder)
	if err != nil {
		log.Println("failed to read sidecar profiles - ", err.Error())
		return
	}

	profiles := make(Profile)

	for _, e := range entries {
		if e.IsDir() {
			fmt.Println(e.Name(), "is a directory")
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

		profiles[e.Name()] = profileData

	}

	for key := range profiles {
		tokens := strings.Split(key, ".")
		if len(tokens) < 2 {
			continue
		}

		fmt.Println(tokens[0])
	}
}

func DeleteProfile(profileName string) {
	if profileName == "" {
		fmt.Println("no profile as specified")
		return
	}

	dirname, err := os.UserHomeDir()
	if err != nil {
		log.Println("failed to get user home dir -", err.Error())
		return
	}

	rootFolder := path.Join(dirname, folderSuffix)
	profileName = fmt.Sprintf("%s.json", profileName)

	profilePath := path.Join(rootFolder, profileName)

	if err := os.Remove(profilePath); err != nil {
		fmt.Println("error deleting profile", err.Error())
		return
	}

}
