package errors

import (
	"fmt"
	"log"
)

// GetFormattedFailMessage returns a formatted fail message for unit tests
func GetFormattedFailMessage(what string, expected string, actual string) string {
	return fmt.Sprintf("%s was incorrect, expected: <%s>, actual: <%s>", what, expected, actual)
}

// GetFormattedErrorMessage returns a formatted error message for the given
// errorToHandle and errorMessage.
func GetFormattedErrorMessage(errorToHandle error, errorMessage string) string {
	result := ""
	if errorMessage == "" {
		errorMessage = "Error making API call"
	}
	if errorToHandle != nil {
		result = fmt.Sprintf(errorMessage+": %v", errorToHandle.Error())
	}
	return result
}

// HandleError handles the given error
func HandleError(errorToHandle error, errorMessage string) {
	if errorToHandle != nil {
		log.Println(GetFormattedErrorMessage(errorToHandle, errorMessage))
	}
}
