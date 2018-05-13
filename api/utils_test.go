package api

import (
	"encoding/json"
	"os"
	"strconv"
	"testing"
)

func TestWriteToJSON(t *testing.T) {
	os.Create("some.json")
	jsonBytes, jsonBytesError := json.Marshal(MapEntryList{})
	if jsonBytesError != nil {
		t.Fatal(jsonBytesError)
	}
	WriteToJSON("some.json", jsonBytes)
	if _, err := os.Stat("some.json"); os.IsNotExist(err) {
		t.Fatal(err)
	}
	os.Remove("some.json")
}

func TestCountOccurences(t *testing.T) {
	inputSlice := []string{"Penis", "Penis", "Hans Meiser"}

	expectedFirstEntry := MapEntry{Key: "Penis", Value: 2}
	expectedSecondEntry := MapEntry{Key: "Hans Meiser", Value: 1}

	actualResult := CountOccurrences(inputSlice)
	actualFirstEntry := actualResult["Penis"]
	actualSecondEntry := actualResult["Hans Meiser"]

	if actualFirstEntry != expectedFirstEntry.Value {
		t.Fatalf(GetFormattedFailMessage("CountOccurences", strconv.Itoa(expectedFirstEntry.Value), strconv.Itoa(actualFirstEntry)))
	}

	if actualSecondEntry != expectedSecondEntry.Value {
		t.Fatalf(GetFormattedFailMessage("CountOccurences", strconv.Itoa(expectedSecondEntry.Value), strconv.Itoa(actualSecondEntry)))
	}
}
