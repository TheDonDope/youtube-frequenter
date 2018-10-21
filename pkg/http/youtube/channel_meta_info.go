package youtube

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
