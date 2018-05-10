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
func GetChannelOverview(service *youtube.Service, inChannel chan ChannelMetaInfo) <-chan ChannelMetaInfo {

	fmt.Println("GetChannelInfo Call")
	outChannel := make(chan ChannelMetaInfo)
	go func() {
		channelMetaInfo := <-inChannel
		fmt.Println(fmt.Sprintf("GetChannelInfo gets: %v", channelMetaInfo))
		call := service.Channels.List("contentDetails,snippet,statistics")
		call = call.ForUsername(channelMetaInfo.CustomURL)

		response, responseError := call.Do()
		HandleError(responseError, "GetChannelInfo Response error!")

		firstItem := response.Items[0]

		// Fill ChannelID
		if channelMetaInfo.ChannelID == "" {
			channelMetaInfo.ChannelID = firstItem.Id
			fmt.Println(channelMetaInfo.ChannelID)
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
		fmt.Println(channelMetaInfo)
		outChannel <- channelMetaInfo
		fmt.Println("would return now.")
	}()
	return outChannel
}

// Exfoliator exfoliates
func Exfoliator(service *youtube.Service, channelMetaInfo ChannelMetaInfo) ChannelMetaInfo {
	initialInChannel := make(chan ChannelMetaInfo)

	go func() {
		initialInChannel <- channelMetaInfo
	}()
	getChannelInfoOutChannel := GetChannelOverview(service, initialInChannel)
	timeout := time.After(5 * time.Second)
	// time.Sleep(time.Second)
	select {
	case channelMetaInfo = <-getChannelInfoOutChannel:
		fmt.Println(fmt.Sprintf("Exfoliator gets: %v", channelMetaInfo))
	case <-timeout:
		fmt.Println("timed out")
		return channelMetaInfo
	}

	return channelMetaInfo
}
