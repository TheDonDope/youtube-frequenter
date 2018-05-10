package service

import youtube "google.golang.org/api/youtube/v3"

// Comment information for a YouTube video
type Comment struct {
	CommentID       string
	AuthorChannelID string
}

// Video information for a YouTube video
type Video struct {
	VideoID           string
	UploaderChannelID string
	Comments          []Comment
}

// Playlist information of a YouTube playlist
type Playlist struct {
	PlaylistID    string
	PlaylistName  string
	PlaylistItems []Video
}

// ChannelMetaInfo of a YouTube channel
type ChannelMetaInfo struct {
	ChannelID       string
	ChannelName     string
	CustomURL       string
	SubscriberCount uint64
	ViewCount       uint64
	Playlists       []Playlist
}

// Search calls the YouTube API
type Search func(service *youtube.Service, channelMetaInfoChannel chan ChannelMetaInfo)
