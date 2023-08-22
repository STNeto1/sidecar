package pkg

import "errors"

var (
	ErrNoProfileSpecified     = errors.New("no profile was specified")
	ErrNoCommandSpecified     = errors.New("no command was specified")
	ErrProfileExists          = errors.New("profile already exists")
	ErrProfileDoesNotExist    = errors.New("profile does not exists")
	ErrFailedToGetProfiles    = errors.New("failed to get profiles")
	ErrFailedToGetProfilePath = errors.New("failed to get profile path")
	ErrFailedToSerializeData  = errors.New("failed to serialize new data")
	ErrFailedToWriteToFile    = errors.New("failed to write to file")
	ErrFailedToDeleteProfile  = errors.New("failed to delete profile")
	ErrFailedToSetEnv         = errors.New("error setting env")
)
