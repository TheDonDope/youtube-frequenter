package api

import (
	"fmt"
	"log"
)

// HandleError handles errors
func HandleError(errorToHandle error, errorMessage string) {
	if errorMessage == "" {
		errorMessage = "Error making API call"
	}
	if errorToHandle != nil {
		log.Fatalf(errorMessage+": %v", errorToHandle.Error())
	}
}

// Printfln ...
func Printfln(format string, a ...interface{}) {
	fmt.Println(fmt.Sprintf(format, a))
}
