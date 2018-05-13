package api

import (
	"fmt"
	"strconv"
	"testing"
)

func TestNextOperationStringer(t *testing.T) {
	if getChannelOverviewOperation := fmt.Sprint(GetChannelOverviewOperation); getChannelOverviewOperation != "GetChannelOverviewOperation" {
		t.Errorf(GetFormattedFailMessage("GetChannelOverviewOperation", "GetChannelOverviewOperation", getChannelOverviewOperation))
	}
	if getVideoIDsOverviewOperation := fmt.Sprint(GetVideoIDsOverviewOperation); getVideoIDsOverviewOperation != "GetVideoIDsOverviewOperation" {
		t.Errorf(GetFormattedFailMessage("GetVideoIDsOverviewOperation", "GetVideoIDsOverviewOperation", getVideoIDsOverviewOperation))
	}
	if getCommentsOverviewOperation := fmt.Sprint(GetCommentsOverviewOperation); getCommentsOverviewOperation != "GetCommentsOverviewOperation" {
		t.Errorf(GetFormattedFailMessage("GetCommentsOverviewOperation", "GetCommentsOverviewOperation", getCommentsOverviewOperation))
	}
	if getObviouslyRelatedChannelsOverviewOperation := fmt.Sprint(GetObviouslyRelatedChannelsOverviewOperation); getObviouslyRelatedChannelsOverviewOperation != "GetObviouslyRelatedChannelsOverviewOperation" {
		t.Errorf(GetFormattedFailMessage("GetObviouslyRelatedChannelsOverviewOperation", "GetObviouslyRelatedChannelsOverviewOperation", getObviouslyRelatedChannelsOverviewOperation))
	}
	if noOperation := fmt.Sprint(NoOperation); noOperation != "NoOperation" {
		t.Errorf(GetFormattedFailMessage("NoOperation", "NoOperation", noOperation))
	}
	var unknownNegativeOperation NextOperation = -1
	if negativeUnknownString := fmt.Sprint(unknownNegativeOperation); negativeUnknownString != "Unknown" {
		t.Errorf(GetFormattedFailMessage("UnknownNegativeOperation", "Unknow", negativeUnknownString))
	}
	var positiveUnknownOperation NextOperation = 5
	if positiveUnknownString := fmt.Sprint(positiveUnknownOperation); positiveUnknownString != "Unknown" {
		t.Errorf(GetFormattedFailMessage("PositiveUnknownOperation", "Unknow", positiveUnknownString))
	}
}

func TestMapEntryStringer(t *testing.T) {
	keyInput := "Foo"
	valueInput := 42
	inputMapEntry := &MapEntry{Key: keyInput, Value: valueInput}
	expectedResult := "Related ChannelID: " + keyInput + ", Number of Occurrences: " + strconv.Itoa(valueInput)
	if actualResult := fmt.Sprint(inputMapEntry); actualResult != expectedResult {
		t.Errorf(GetFormattedFailMessage("MapEntry#Stringer", expectedResult, actualResult))
	}
}

func TestMapEntryListLen(t *testing.T) {
	inputMapEntryList := MapEntryList{}
	expectedResult := 0
	if actualResult := inputMapEntryList.Len(); actualResult != expectedResult {
		t.Errorf(GetFormattedFailMessage("MapEntryList#Len", strconv.Itoa(expectedResult), strconv.Itoa(actualResult)))
	}
}

func TestMapEntryListLess(t *testing.T) {
	inputMapEntryList := MapEntryList{}
	inputMapEntryList = append(inputMapEntryList, MapEntry{Key: "Hans", Value: 1}, MapEntry{Key: "Peter", Value: 5})
	inputI := 0
	inputJ := 1
	expectedResult := true
	if actualResult := inputMapEntryList.Less(inputI, inputJ); actualResult != expectedResult {
		t.Errorf(GetFormattedFailMessage("MapEntryList#Less", strconv.FormatBool(expectedResult), strconv.FormatBool(actualResult)))
	}

	expectedResult = false
	if actualResult := inputMapEntryList.Less(inputJ, inputI); actualResult != expectedResult {
		t.Errorf(GetFormattedFailMessage("MapEntryList#Less", strconv.FormatBool(expectedResult), strconv.FormatBool(actualResult)))
	}
}

func TestMapEntrySwap(t *testing.T) {
	inputMapEntryList := MapEntryList{}
	oneEntry := MapEntry{Key: "Hans", Value: 1}
	anotherEntry := MapEntry{Key: "Peter", Value: 5}

	inputMapEntryList = append(inputMapEntryList, oneEntry, anotherEntry)
	inputI := 0
	inputJ := 1
	inputMapEntryList.Swap(inputI, inputJ)

	if actualFirstElement := inputMapEntryList[inputI]; actualFirstElement != anotherEntry {
		t.Errorf(GetFormattedFailMessage("MapEntryList#swap", anotherEntry.Key, actualFirstElement.Key))
	}

	if actualSecondElement := inputMapEntryList[inputJ]; actualSecondElement != oneEntry {
		t.Errorf(GetFormattedFailMessage("MapEntryList#swap", oneEntry.Key, actualSecondElement.Key))
	}
}
