package types

import (
	"log"
	"sort"
)

// MapEntryList is a type for an array of MapEntrys
type MapEntryList []MapEntry

// Len implements the Sort interface method Len for the MapEntryList type
func (p MapEntryList) Len() int { return len(p) }

// Less implements the Sort interface method Less for the MapEntryList type
func (p MapEntryList) Less(i, j int) bool { return p[i].Value < p[j].Value }

// Swap implements the Sort interface method Swap for the MapEntryList type
func (p MapEntryList) Swap(i, j int) { p[i], p[j] = p[j], p[i] }

// RankByWordCount returns a list of sorted MapEntrys
func (p MapEntryList) RankByWordCount(wordFrequencies map[string]int) MapEntryList {
	mapEntryList := make(MapEntryList, len(wordFrequencies))
	i := 0
	for key, value := range wordFrequencies {
		mapEntryList[i] = MapEntry{key, value}
		i++
	}
	sort.Sort(sort.Reverse(mapEntryList))
	return mapEntryList
}

// PrintResults prints the results
func (p MapEntryList) PrintResults() {
	for _, item := range p {
		log.Println(item)
	}
}
