package test

import (
	youtube "google.golang.org/api/youtube/v3"
)

// YouTuberMock is a mock implementation for the YouTuber interface
type YouTuberMock struct{}

// ChannelsList returns the result of the Channels.list API call
func (mock YouTuberMock) ChannelsList(service *youtube.Service, channelID string, customURL string) (*youtube.ChannelListResponse, error) {
	result := &youtube.ChannelListResponse{}

	return result, nil
}

// PlaylistItemsList returns the result of the PlaylistItems.list API call
func (mock YouTuberMock) PlaylistItemsList(service *youtube.Service, playlistID string, playlistName string) (*youtube.PlaylistItemListResponse, error) {
	result := &youtube.PlaylistItemListResponse{}

	return result, nil
}

// CommentThreadsList returns the result of the CommentThreads.list API call
func (mock YouTuberMock) CommentThreadsList(service *youtube.Service, videoID string) (*youtube.CommentThreadListResponse, error) {
	result := &youtube.CommentThreadListResponse{}

	return result, nil
}

// VideosList returns the result of the Videos.list API call
func (mock YouTuberMock) VideosList(service *youtube.Service, videoIDs string) (*youtube.VideoListResponse, error) {
	result := &youtube.VideoListResponse{}

	return result, nil
}

// GetService returns a service to interact with the YouTube API
func (mock YouTuberMock) GetService() (*youtube.Service, error) {
	result := &youtube.Service{}

	return result, nil
}
