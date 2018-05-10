package service

import (
	"fmt"
	"time"

	"google.golang.org/api/youtube/v3"
)

// GetPlaylistIDByChannelMetaInfo returns the id of the playlist for the given
// search query. You can either query by channelID or customURL.
func GetPlaylistIDByChannelMetaInfo(service *youtube.Service, channelMetaInfo ChannelMetaInfo) ChannelMetaInfo {
	call := service.Channels.List("contentDetails")

	if channelMetaInfo.ChannelID != "" {
		call = call.Id(channelMetaInfo.ChannelID)
	} else if channelMetaInfo.CustomURL != "" {
		call = call.ForUsername(channelMetaInfo.CustomURL)
	}

	response, responseError := call.Do()
	HandleError(responseError, "Response error!")

	firstItem := response.Items[0]
	if channelMetaInfo.Playlists == nil {
		channelMetaInfo.Playlists = append(channelMetaInfo.Playlists, Playlist{firstItem.ContentDetails.RelatedPlaylists.Uploads, "uploads", nil})
		channelMetaInfo.Playlists = append(channelMetaInfo.Playlists, Playlist{firstItem.ContentDetails.RelatedPlaylists.Favorites, "favorites", nil})
	}

	return channelMetaInfo
}

// MetaSearch searches meta
func MetaSearch(searchName string) Search {
	return func(service *youtube.Service, channelMetaInfo ChannelMetaInfo) ChannelMetaInfo {
		fmt.Println(fmt.Sprintf("%v starts running with query: %v", searchName, channelMetaInfo))

		// priority 1: ....
		if channelMetaInfo.ChannelID == "" {
			channelMetaInfo = GetPlaylistIDByChannelMetaInfo(service, channelMetaInfo)
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
		channelMetaInfoChannel <- FirstResponder(service, channelMetaInfo, MetaSearch("PlaylistID1"), MetaSearch("PlaylistID2"))
	}()
	go func() {
		channelMetaInfoChannel <- FirstResponder(service, channelMetaInfo, MetaSearch("PlaylistID1"), MetaSearch("PlaylistID2"))
	}()

	timeout := time.After(10 * time.Second)
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
