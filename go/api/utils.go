package api

import (
	"fmt"
)

// GetFormattedErrorMessage returns a formatted error message for the given
// errorToHandle and errorMessage.
func GetFormattedErrorMessage(errorToHandle error, errorMessage string) string {
	if errorMessage == "" {
		errorMessage = "Error making API call"
	}
	if errorToHandle != nil {
		return fmt.Sprintf(errorMessage+": %v", errorToHandle.Error())
	}
	return ""
}

// Printfln prints a line with a formatted string.
func Printfln(format string, a ...interface{}) {
	fmt.Println(fmt.Sprintf(format, a))
}
