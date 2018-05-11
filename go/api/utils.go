package api

import "log"

// HandleError handles errors
func HandleError(errorToHandle error, errorMessage string) {
	if errorMessage == "" {
		errorMessage = "Error making API call"
	}
	if errorToHandle != nil {
		log.Fatalf(errorMessage+": %v", errorToHandle.Error())
	}
}
