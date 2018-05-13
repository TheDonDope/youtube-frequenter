package api

import youtube "google.golang.org/api/youtube/v3"

// Exfoliator describes methods to interact with the exfoliator api
type Exfoliator interface {

	// GetChannelOverview implements a Search which fills the most basic info about a channel.
	// To use the method simply pass a ChannelMetaInfo with either the ChannelID or CustomURL set.
	// On search completion the following attributes will be filled, if not already set:
	// - ChannelID
	// - CustomURL
	// - ChannelName
	// - Playlists (With their PlaylistID and PlaylistName)
	// - SubscriberCount
	// - ViewCount
	GetChannelOverview(service *youtube.Service, serviceImpl YouTuber, monoChannel chan ChannelMetaInfo)

	// GetVideoIDsOverview gets all videos.
	GetVideoIDsOverview(service *youtube.Service, serviceImpl YouTuber, monoChannel chan ChannelMetaInfo)

	// GetCommentsOverview returns the overview over the comments
	GetCommentsOverview(service *youtube.Service, serviceImpl YouTuber, monoChannel chan ChannelMetaInfo)

	// GetObviouslyRelatedChannelsOverview gets the related channels for a YouTube channel
	GetObviouslyRelatedChannelsOverview(service *youtube.Service, serviceImpl YouTuber, monoChannel chan ChannelMetaInfo, lastButNotLeastChannel chan ChannelMetaInfo)

	// CreateInitialChannelMetaInfo creates the initial request context
	CreateInitialChannelMetaInfo() ChannelMetaInfo

	// Exfoliate returns the result
	Exfoliate(service *youtube.Service, serviceImpl YouTuber, channelMetaInfo ChannelMetaInfo) ChannelMetaInfo

	// AnalyseChannelMetaInfo prints additional information for a given channelMetaInfo.
	AnalyseChannelMetaInfo(channelMetaInfo *ChannelMetaInfo)
}
