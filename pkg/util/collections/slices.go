package collections

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
