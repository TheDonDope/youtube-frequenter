package api

import (
	"fmt"
	"log"
	"os"
)

// Printfln prints a line with a formatted string.
func Printfln(format string, a ...interface{}) {
	log.Println(fmt.Sprintf(format, a))
}

// WriteToJSON writes the given jsonBytes as a json for the given path
func WriteToJSON(path string, jsonBytes []byte) {
	if _, err := os.Stat(path); err == nil {
		os.Remove(path)
	}

	jsonFile, jsonFileError := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0644)
	defer jsonFile.Close()
	HandleError(jsonFileError, "Error opening JSON file for path: "+path)
	jsonFile.Write(jsonBytes)
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
