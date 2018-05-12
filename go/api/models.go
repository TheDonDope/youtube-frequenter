package api

// Comment information for a YouTube video
type Comment struct {
	CommentID       string
	AuthorChannelID string
}

// Video information for a YouTube video
type Video struct {
	VideoID           string
	UploaderChannelID string
	Comments          []*Comment
}

// Playlist information of a YouTube playlist
type Playlist struct {
	PlaylistID    string
	PlaylistItems []*Video
}

// ChannelMetaInfo of a YouTube channel
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
