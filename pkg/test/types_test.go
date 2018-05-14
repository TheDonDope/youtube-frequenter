package test

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/TheDonDope/youtube-frequenter/pkg/types"
	ourErrors "github.com/TheDonDope/youtube-frequenter/pkg/util/errors"
)

func TestNextOperationStringer(t *testing.T) {
	if getChannelOverviewOperation := fmt.Sprint(types.GetChannelOverviewOperation); getChannelOverviewOperation != "GetChannelOverviewOperation" {
		t.Errorf(ourErrors.GetFormattedFailMessage("GetChannelOverviewOperation", "GetChannelOverviewOperation", getChannelOverviewOperation))
	}
	if getVideoIDsOverviewOperation := fmt.Sprint(types.GetVideoIDsOverviewOperation); getVideoIDsOverviewOperation != "GetVideoIDsOverviewOperation" {
		t.Errorf(ourErrors.GetFormattedFailMessage("GetVideoIDsOverviewOperation", "GetVideoIDsOverviewOperation", getVideoIDsOverviewOperation))
	}
	if getCommentsOverviewOperation := fmt.Sprint(types.GetCommentsOverviewOperation); getCommentsOverviewOperation != "GetCommentsOverviewOperation" {
		t.Errorf(ourErrors.GetFormattedFailMessage("GetCommentsOverviewOperation", "GetCommentsOverviewOperation", getCommentsOverviewOperation))
	}
	if getObviouslyRelatedChannelsOverviewOperation := fmt.Sprint(types.GetObviouslyRelatedChannelsOverviewOperation); getObviouslyRelatedChannelsOverviewOperation != "GetObviouslyRelatedChannelsOverviewOperation" {
		t.Errorf(ourErrors.GetFormattedFailMessage("GetObviouslyRelatedChannelsOverviewOperation", "GetObviouslyRelatedChannelsOverviewOperation", getObviouslyRelatedChannelsOverviewOperation))
	}
	if noOperation := fmt.Sprint(types.NoOperation); noOperation != "NoOperation" {
		t.Errorf(ourErrors.GetFormattedFailMessage("NoOperation", "NoOperation", noOperation))
	}
	var unknownNegativeOperation types.NextOperation = -1
	if negativeUnknownString := fmt.Sprint(unknownNegativeOperation); negativeUnknownString != "Unknown" {
		t.Errorf(ourErrors.GetFormattedFailMessage("UnknownNegativeOperation", "Unknow", negativeUnknownString))
	}
	var positiveUnknownOperation types.NextOperation = 5
	if positiveUnknownString := fmt.Sprint(positiveUnknownOperation); positiveUnknownString != "Unknown" {
		t.Errorf(ourErrors.GetFormattedFailMessage("PositiveUnknownOperation", "Unknow", positiveUnknownString))
	}
}

func TestMapEntryStringer(t *testing.T) {
	keyInput := "Foo"
	valueInput := 42
	inputMapEntry := &types.MapEntry{Key: keyInput, Value: valueInput}
	expectedResult := "Related ChannelID: " + keyInput + ", Number of Occurrences: " + strconv.Itoa(valueInput)
	if actualResult := fmt.Sprint(inputMapEntry); actualResult != expectedResult {
		t.Errorf(ourErrors.GetFormattedFailMessage("MapEntry#Stringer", expectedResult, actualResult))
	}
}

func TestMapEntryListLen(t *testing.T) {
	inputMapEntryList := types.MapEntryList{}
	expectedResult := 0
	if actualResult := inputMapEntryList.Len(); actualResult != expectedResult {
		t.Errorf(ourErrors.GetFormattedFailMessage("MapEntryList#Len", strconv.Itoa(expectedResult), strconv.Itoa(actualResult)))
	}
}

func TestMapEntryListLess(t *testing.T) {
	inputMapEntryList := types.MapEntryList{}
	inputMapEntryList = append(inputMapEntryList, types.MapEntry{Key: "Hans", Value: 1}, types.MapEntry{Key: "Peter", Value: 5})
	inputI := 0
	inputJ := 1
	expectedResult := true
	if actualResult := inputMapEntryList.Less(inputI, inputJ); actualResult != expectedResult {
		t.Errorf(ourErrors.GetFormattedFailMessage("MapEntryList#Less", strconv.FormatBool(expectedResult), strconv.FormatBool(actualResult)))
	}

	expectedResult = false
	if actualResult := inputMapEntryList.Less(inputJ, inputI); actualResult != expectedResult {
		t.Errorf(ourErrors.GetFormattedFailMessage("MapEntryList#Less", strconv.FormatBool(expectedResult), strconv.FormatBool(actualResult)))
	}
}

func TestMapEntrySwap(t *testing.T) {
	inputMapEntryList := types.MapEntryList{}
	oneEntry := types.MapEntry{Key: "Hans", Value: 1}
	anotherEntry := types.MapEntry{Key: "Peter", Value: 5}

	inputMapEntryList = append(inputMapEntryList, oneEntry, anotherEntry)
	inputI := 0
	inputJ := 1
	inputMapEntryList.Swap(inputI, inputJ)

	if actualFirstElement := inputMapEntryList[inputI]; actualFirstElement != anotherEntry {
		t.Errorf(ourErrors.GetFormattedFailMessage("MapEntryList#swap", anotherEntry.Key, actualFirstElement.Key))
	}

	if actualSecondElement := inputMapEntryList[inputJ]; actualSecondElement != oneEntry {
		t.Errorf(ourErrors.GetFormattedFailMessage("MapEntryList#swap", oneEntry.Key, actualSecondElement.Key))
	}
}

func TestRankByWordCount(t *testing.T) {
	inputMap := make(map[string]int)
	oneEntry := types.MapEntry{Key: "wwwKenFMde", Value: 420}
	anotherEntry := types.MapEntry{Key: "nuovisoTV", Value: 62}

	inputMap[oneEntry.Key] = oneEntry.Value
	inputMap[anotherEntry.Key] = anotherEntry.Value

	expectedResult := types.MapEntryList{anotherEntry, oneEntry}
	expectedFirstElement := expectedResult[1]
	actualFirstElement := types.MapEntryList{}.RankByWordCount(inputMap)[0]

	if actualFirstElement != expectedFirstElement {
		t.Errorf(ourErrors.GetFormattedFailMessage("MapEntryList#RankByWordCount", expectedFirstElement.Key, actualFirstElement.Key))
	}
}
