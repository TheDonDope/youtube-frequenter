package api

import (
	"fmt"
	"log"
	"sort"
)

// GetFormattedFailMessage returns a formatted fail message for unit tests
func GetFormattedFailMessage(what string, expected string, actual string) string {
	return fmt.Sprintf("%s was incorrect, expected: <%s>, actual: <%s>", what, expected, actual)
}

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

// HandleError handles the given error
func HandleError(errorToHandle error, errorMessage string) {
	if errorToHandle != nil {
		log.Println(GetFormattedErrorMessage(errorToHandle, errorMessage))
	}
}

// Printfln prints a line with a formatted string.
func Printfln(format string, a ...interface{}) {
	log.Println(fmt.Sprintf(format, a))
}

// CountOccurrences returns a map with the string value as key, and the number of occurences as value.
func CountOccurrences(stringSlice []string) map[string]int {
	elementToFrequencyMap := make(map[string]int)

	for _, item := range stringSlice {
		// check if the item/element exist in the duplicate_frequency map
		_, exist := elementToFrequencyMap[item]

		if exist {
			elementToFrequencyMap[item]++ // increase counter by 1 if already in the map
		} else {
			elementToFrequencyMap[item] = 1 // else start counting from 1
		}
	}
	return elementToFrequencyMap
}

// RankByWordCount returns a list of sorted MapEntrys
func RankByWordCount(wordFrequencies map[string]int) MapEntryList {
	mapEntryList := make(MapEntryList, len(wordFrequencies))
	i := 0
	for key, value := range wordFrequencies {
		mapEntryList[i] = MapEntry{key, value}
		i++
	}
	sort.Sort(sort.Reverse(mapEntryList))
	return mapEntryList
}

// Len implements the Sort interface method Len for the MapEntryList type
func (p MapEntryList) Len() int { return len(p) }

// Less implements the Sort interface method Less for the MapEntryList type
func (p MapEntryList) Less(i, j int) bool { return p[i].Value < p[j].Value }

// Swap implements the Sort interface method Swap for the MapEntryList type
func (p MapEntryList) Swap(i, j int) { p[i], p[j] = p[j], p[i] }
