package files

import (
	"os"

	"gitlab.com/TheDonDope/youtube-frequenter/pkg/util/errors"
)

// WriteToJSON writes the given jsonBytes as a json for the given path
func WriteToJSON(path string, jsonBytes []byte) {
	if _, err := os.Stat(path); err == nil {
		os.Remove(path)
	}

	jsonFile, jsonFileError := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0644)
	defer jsonFile.Close()
	errors.HandleError(jsonFileError, "Error opening JSON file for path: "+path)
	jsonFile.Write(jsonBytes)
}
