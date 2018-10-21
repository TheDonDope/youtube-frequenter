package mocks

import (
	youtubeV3 "google.golang.org/api/youtube/v3"
)

// YouTubeService is a mock implementation for the YouTube Service interface
type YouTubeService struct{}

// ChannelsList returns the result of the Channels.list API call
func (m *YouTubeService) ChannelsList(channelID string, customURL string) (*youtubeV3.ChannelListResponse, error) {
	result := &youtubeV3.ChannelListResponse{}

	return result, nil
}

// PlaylistItemsList returns the result of the PlaylistItems.list API call
func (m *YouTubeService) PlaylistItemsList(playlistID string, playlistName string) (*youtubeV3.PlaylistItemListResponse, error) {
	result := &youtubeV3.PlaylistItemListResponse{}

	return result, nil
}

// CommentThreadsList returns the result of the CommentThreads.list API call
func (m *YouTubeService) CommentThreadsList(videoID string) (*youtubeV3.CommentThreadListResponse, error) {
	result := &youtubeV3.CommentThreadListResponse{}

	return result, nil
}

// VideosList returns the result of the Videos.list API call
func (m *YouTubeService) VideosList(videoIDs string) (*youtubeV3.VideoListResponse, error) {
	result := &youtubeV3.VideoListResponse{}

	return result, nil
}

// NewYouTube returns a service to interact with the YouTube API
func (m *YouTubeService) NewYouTube() (*youtubeV3.Service, error) {
	result := &youtubeV3.Service{}

	return result, nil
}
