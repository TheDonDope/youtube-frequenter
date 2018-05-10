package service

import (
	"fmt"
	"time"

	"google.golang.org/api/youtube/v3"
)

// GetChannelInfo implements a Search which fills the most basic info about a channel.
// To use the method simply pass a ChannelMetaInfo with either the ChannelID or CustomURL set.
// On search completion the following attributes will be filled, if not already set:
// - ChannelID
// - CustomURL
// - ChannelName
// - Playlists (With their PlaylistID and PlaylistName)
// - SubscriberCount
// - ViewCount
func GetChannelInfo() Search {
	return func(service *youtube.Service, channelMetaInfo ChannelMetaInfo) ChannelMetaInfo {
		fmt.Println(fmt.Sprintf("Channel Info search starts running with meta info: %v", channelMetaInfo))
		call := service.Channels.List("contentDetails,snippet,statistics")
		call = call.ForUsername(channelMetaInfo.CustomURL)

		response, responseError := call.Do()
		HandleError(responseError, "GetChannelInfo Response error!")

		firstItem := response.Items[0]

		// Fill ChannelID
		if channelMetaInfo.ChannelID == "" {
			channelMetaInfo.ChannelID = firstItem.Id
		}

		// Fill CustomURL
		if channelMetaInfo.CustomURL == "" {
			channelMetaInfo.CustomURL = firstItem.Snippet.CustomUrl
		}

		// Fill ChannelName
		if channelMetaInfo.ChannelName == "" {
			channelMetaInfo.ChannelName = firstItem.Snippet.Title
		}

		// Fill Playlists
		if channelMetaInfo.Playlists == nil {
			channelMetaInfo.Playlists = append(channelMetaInfo.Playlists, Playlist{firstItem.ContentDetails.RelatedPlaylists.Uploads, "uploads", nil})
			channelMetaInfo.Playlists = append(channelMetaInfo.Playlists, Playlist{firstItem.ContentDetails.RelatedPlaylists.Favorites, "favorites", nil})
		}

		// Fill SubscriberCount
		if channelMetaInfo.SubscriberCount == 0 {
			channelMetaInfo.SubscriberCount = firstItem.Statistics.SubscriberCount
		}

		// Fill ViewCount
		if channelMetaInfo.ViewCount == 0 {
			channelMetaInfo.ViewCount = firstItem.Statistics.ViewCount
		}

		return channelMetaInfo
	}
}

// FirstResponder returns the first replica search result, whichever was the faster
func FirstResponder(service *youtube.Service, channelMetaInfo ChannelMetaInfo, replicas ...Search) ChannelMetaInfo {
	channelMetaInfoChannel := make(chan ChannelMetaInfo)
	searchReplica := func(index int) { channelMetaInfoChannel <- replicas[index](service, channelMetaInfo) }
	for replicaIndex := range replicas {
		go searchReplica(replicaIndex)
	}
	return <-channelMetaInfoChannel
}

// Exfoliator exfoliates
func Exfoliator(service *youtube.Service, channelMetaInfo ChannelMetaInfo) ChannelMetaInfo {
	channelMetaInfoChannel := make(chan ChannelMetaInfo)

	go func() {
		channelMetaInfoChannel <- FirstResponder(service, channelMetaInfo, GetChannelInfo(), GetChannelInfo())
	}()

	timeout := time.After(3 * time.Second)
	for i := 0; i < 2; i++ {
		select {
		case channelMetaInfoFromChannel := <-channelMetaInfoChannel:
			channelMetaInfo = channelMetaInfoFromChannel
		case <-timeout:
			fmt.Println("timed out")
			return channelMetaInfo
		}
	}
	return channelMetaInfo
}
