package api

import (
	"fmt"
	"sort"
)

// NextOperation declares an enum type for operation names
type NextOperation int

const (
	// GetChannelOverviewOperation is a NextOperation enum value
	GetChannelOverviewOperation NextOperation = iota
	// GetVideoIDsOverviewOperation is a NextOperation enum value
	GetVideoIDsOverviewOperation
	// GetCommentsOverviewOperation is a NextOperation enum value
	GetCommentsOverviewOperation
	// GetObviouslyRelatedChannelsOverviewOperation is a NextOperation enum value
	GetObviouslyRelatedChannelsOverviewOperation
	// NoOperation is a NextOperation enum value
	NoOperation
)

// This function is the Stringer for the NextOperation enum
func (operation NextOperation) String() string {
	operationNames := [...]string{
		"GetChannelOverviewOperation",
		"GetVideoIDsOverviewOperation",
		"GetCommentsOverviewOperation",
		"GetObviouslyRelatedChannelsOverviewOperation",
		"NoOperation"}

	if operation < GetChannelOverviewOperation || operation > NoOperation {
		return "Unknown"
	}

	return operationNames[operation]
}

// String implements the String interface method String for the MapEntry type
func (p MapEntry) String() string {
	return fmt.Sprintf("Related ChannelID: %v, Number of Occurrences: %v", p.Key, p.Value)
}

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

// Comment is the struct for information on comments
type Comment struct {
	CommentID       string
	AuthorChannelID string
}

// Video is the struct for information on videos
type Video struct {
	VideoID           string
	UploaderChannelID string
	Comments          []*Comment
}

// Playlist is the struct for information on playlists
type Playlist struct {
	PlaylistID    string
	PlaylistItems []*Video
}

// ChannelMetaInfo is the struct for information YouTube channels
type ChannelMetaInfo struct {
	ChannelID                  string
	ChannelName                string
	CustomURL                  string
	SubscriberCount            uint64
	ViewCount                  uint64
	Playlists                  map[string]*Playlist
	CommentAuthorChannelIDs    []string
	ObviouslyRelatedChannelIDs []string
	NextOperation              NextOperation
}

// MapEntry is a stract for a map entry
type MapEntry struct {
	Key   string `json:"key"`
	Value int    `json:"value"`
}

// MapEntryList is a type for an array of MapEntrys
type MapEntryList []MapEntry
