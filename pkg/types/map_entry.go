package types

import "fmt"

// MapEntry is a struct for a map entry
type MapEntry struct {
	Key   string `json:"key"`
	Value int    `json:"value"`
}

// String implements the String interface method String for the MapEntry type
func (p MapEntry) String() string {
	return fmt.Sprintf("Related ChannelID: %v, Number of Occurrences: %v", p.Key, p.Value)
}
