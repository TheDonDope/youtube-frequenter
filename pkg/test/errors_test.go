package test

import (
	"errors"
	"testing"

	ourErrors "gitlab.com/TheDonDope/youtube-frequenter/pkg/util/errors"
)

func TestGetFormattedFailMessage(t *testing.T) {
	whatInput := "Something"
	expectedInput := "expected"
	actualInput := "actual"
	expectedResult := whatInput + " was incorrect, expected: <" + expectedInput + ">, actual: <" + actualInput + ">"
	if actualResult := ourErrors.GetFormattedFailMessage(whatInput, expectedInput, actualInput); actualResult != expectedResult {
		t.Errorf(ourErrors.GetFormattedFailMessage("GetFormattedFailMessage", expectedResult, actualResult))
	}
}

func TestGetFormattedErrorMessage(t *testing.T) {
	someError := errors.New("Some error")
	errorMessage := ""
	expectedResult := "Error making API call: Some error"
	if actualResult := ourErrors.GetFormattedErrorMessage(someError, errorMessage); actualResult != expectedResult {
		t.Errorf(ourErrors.GetFormattedFailMessage("GetFormattedErrorMessage", expectedResult, actualResult))
	}
}
