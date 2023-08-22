package pkg

import (
	"github.com/rs/zerolog/log"
)

func DisplayProfiles(profiles Profile) {
	for profileName := range profiles {
		log.Info().Msgf("> %s", profileName)
	}
}

func DisplayProfile(profileName string, profileData ProfileData) {
	log.Info().Msgf("> %s", profileName)
	for key, value := range profileData {
		log.Info().Msgf(">> %s=%s", key, value)
	}
}
