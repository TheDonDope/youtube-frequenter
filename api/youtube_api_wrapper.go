package api

import (
	"google.golang.org/api/youtube/v3"
)

// ChannelsList returns the result of the Channels.list API call
func ChannelsList(service *youtube.Service, channelID string, customURL string) (*youtube.ChannelListResponse, error) {
	call := service.Channels.List("contentDetails,snippet,statistics")
	if channelID != "" {
		call = call.Id(channelID)
	} else if customURL != "" {
		call = call.ForUsername(customURL)
	}
	return call.Do()
}

// PlaylistItemsList returns the result of the PlaylistItems.list API call
func PlaylistItemsList(service *youtube.Service, playlistID string, playlistName string) (*youtube.PlaylistItemListResponse, error) {
	call := service.PlaylistItems.List("contentDetails,snippet").PlaylistId(playlistID)
	if playlistName == "uploads" {
		call = call.MaxResults(Opts.MaxResultsUploadedVideos)
	} else if playlistName == "favorites" {
		call = call.MaxResults(Opts.MaxResultsFavouritedVideos)
	}
	return call.Do()
}

// CommentThreadsList returns the result of the CommentThreads.list API call
func CommentThreadsList(service *youtube.Service, videoID string) (*youtube.CommentThreadListResponse, error) {
	return service.CommentThreads.List("snippet").VideoId(videoID).MaxResults(Opts.MaxResultsCommentPerVideo).Do()
}

// VideosList returns the result of the Videos.list API call
func VideosList(service *youtube.Service, videoIDs string) (*youtube.VideoListResponse, error) {
	return service.Videos.List("snippet").Id(videoIDs).Do()
}
