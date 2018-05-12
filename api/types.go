package api

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
