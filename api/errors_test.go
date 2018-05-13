package api

import (
	"errors"
	"testing"
)

func TestGetFormattedFailMessage(t *testing.T) {
	whatInput := "Something"
	expectedInput := "expected"
	actualInput := "actual"
	expectedResult := whatInput + " was incorrect, expected: <" + expectedInput + ">, actual: <" + actualInput + ">"
	if actualResult := GetFormattedFailMessage(whatInput, expectedInput, actualInput); actualResult != expectedResult {
		t.Errorf(GetFormattedFailMessage("GetFormattedFailMessage", expectedResult, actualResult))
	}
}

func TestGetFormattedErrorMessage(t *testing.T) {
	someError := errors.New("Some error")
	errorMessage := ""
	expectedResult := "Error making API call: Some error"
	if actualResult := GetFormattedErrorMessage(someError, errorMessage); actualResult != expectedResult {
		t.Errorf(GetFormattedFailMessage("GetFormattedErrorMessage", expectedResult, actualResult))
	}
}
