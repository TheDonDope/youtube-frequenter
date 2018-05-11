package api

import (
	"fmt"
	"log"
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

// AnaliseChannelMetaInfo prints additional information for a given channelMetaInfo.
func AnaliseChannelMetaInfo(channelMetaInfo *ChannelMetaInfo) {

	relatedChannelIDToNumberOfOccurrences := CountOccurrences(channelMetaInfo.ObviouslyRelatedChannelIDs)

	for key, value := range relatedChannelIDToNumberOfOccurrences {
		log.Println(fmt.Sprintf("Related ChannelID: %v, Number of Occurrences: %v", key, value))
	}
}
