package api

import youtube "google.golang.org/api/youtube/v3"

// YouTuber describes methods to interact with the YouTube API
type YouTuber interface {
	// ChannelsList returns the result of the Channels.list API call
	ChannelsList(service *youtube.Service, channelID string, customURL string) (*youtube.ChannelListResponse, error)

	// PlaylistItemsList returns the result of the PlaylistItems.list API call
	PlaylistItemsList(service *youtube.Service, playlistID string, playlistName string) (*youtube.PlaylistItemListResponse, error)

	// CommentThreadsList returns the result of the CommentThreads.list API call
	CommentThreadsList(service *youtube.Service, videoID string) (*youtube.CommentThreadListResponse, error)

	// VideosList returns the result of the Videos.list API call
	VideosList(service *youtube.Service, videoIDs string) (*youtube.VideoListResponse, error)

	// GetServicecd returns a service to interact with the YouTube API
	GetService() (*youtube.Service, error)
}
