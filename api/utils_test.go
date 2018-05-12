package api

import (
	"errors"
	"testing"
)

func TestGetFormattedErrorMessage(t *testing.T) {
	someError := errors.New("Some error")
	errorMessage := ""
	expectedResult := "Error making API call: Some error"
	if actualResult := GetFormattedErrorMessage(someError, errorMessage); actualResult != expectedResult {
		t.Errorf(GetFormattedFailMessage("GetFormattedErrorMessage", expectedResult, actualResult))
	}
}
