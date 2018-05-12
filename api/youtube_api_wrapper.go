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
