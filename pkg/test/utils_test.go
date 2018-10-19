package test

import (
	"encoding/json"
	"os"
	"strconv"
	"testing"

	"gitlab.com/TheDonDope/youtube-frequenter/pkg/types"
	"gitlab.com/TheDonDope/youtube-frequenter/pkg/util/collections"
	ourErrors "gitlab.com/TheDonDope/youtube-frequenter/pkg/util/errors"
	"gitlab.com/TheDonDope/youtube-frequenter/pkg/util/files"
)

func TestWriteToJSON(t *testing.T) {
	os.Create("some.json")
	jsonBytes, jsonBytesError := json.Marshal(types.MapEntryList{})
	if jsonBytesError != nil {
		t.Fatal(jsonBytesError)
	}
	files.WriteToJSON("some.json", jsonBytes)
	if _, err := os.Stat("some.json"); os.IsNotExist(err) {
		t.Fatal(err)
	}
	os.Remove("some.json")
}

func TestCountOccurences(t *testing.T) {
	inputSlice := []string{"Penis", "Penis", "Hans Meiser"}

	expectedFirstEntry := types.MapEntry{Key: "Penis", Value: 2}
	expectedSecondEntry := types.MapEntry{Key: "Hans Meiser", Value: 1}

	actualResult := collections.CountOccurrences(inputSlice)
	actualFirstEntry := actualResult["Penis"]
	actualSecondEntry := actualResult["Hans Meiser"]

	if actualFirstEntry != expectedFirstEntry.Value {
		t.Fatalf(ourErrors.GetFormattedFailMessage("CountOccurences", strconv.Itoa(expectedFirstEntry.Value), strconv.Itoa(actualFirstEntry)))
	}

	if actualSecondEntry != expectedSecondEntry.Value {
		t.Fatalf(ourErrors.GetFormattedFailMessage("CountOccurences", strconv.Itoa(expectedSecondEntry.Value), strconv.Itoa(actualSecondEntry)))
	}
}
