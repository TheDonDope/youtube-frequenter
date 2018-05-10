package service

import (
	"fmt"
	"time"

	"google.golang.org/api/youtube/v3"
)

// GetChannelOverview implements a Search which fills the most basic info about a channel.
// To use the method simply pass a ChannelMetaInfo with either the ChannelID or CustomURL set.
// On search completion the following attributes will be filled, if not already set:
// - ChannelID
// - CustomURL
// - ChannelName
// - Playlists (With their PlaylistID and PlaylistName)
// - SubscriberCount
// - ViewCount
func GetChannelOverview(service *youtube.Service, inChannel chan ChannelMetaInfo) <-chan ChannelMetaInfo {
	fmt.Println("Begin GetChannelOverview")
	outChannel := make(chan ChannelMetaInfo)

	go func() {
		fmt.Println(fmt.Sprintf("Starting goroutine in GetChannelOverview"))
		channelMetaInfo := <-inChannel
		fmt.Println(fmt.Sprintf("Input channelMetaInfo: %+v", channelMetaInfo))
		call := service.Channels.List("contentDetails,snippet,statistics")
		call = call.ForUsername(channelMetaInfo.CustomURL)

		response, responseError := call.Do()
		HandleError(responseError, "GetChannelOverview Response error!")

		firstItem := response.Items[0]

		// Fill ChannelID
		if channelMetaInfo.ChannelID == "" {
			channelMetaInfo.ChannelID = firstItem.Id
			fmt.Println(fmt.Sprintf("channelMetaInfo.ChannelID: %s", channelMetaInfo.ChannelID))
		}

		// Fill CustomURL
		if channelMetaInfo.CustomURL == "" {
			channelMetaInfo.CustomURL = firstItem.Snippet.CustomUrl
			fmt.Println(fmt.Sprintf("channelMetaInfo.CustomURL: %s", channelMetaInfo.CustomURL))
		}

		// Fill ChannelName
		if channelMetaInfo.ChannelName == "" {
			channelMetaInfo.ChannelName = firstItem.Snippet.Title
			fmt.Println(fmt.Sprintf("channelMetaInfo.ChannelName: %s", channelMetaInfo.ChannelName))
		}

		// Fill Playlists
		if channelMetaInfo.Playlists == nil {
			channelMetaInfo.Playlists = append(channelMetaInfo.Playlists, Playlist{firstItem.ContentDetails.RelatedPlaylists.Uploads, "uploads", nil})
			channelMetaInfo.Playlists = append(channelMetaInfo.Playlists, Playlist{firstItem.ContentDetails.RelatedPlaylists.Favorites, "favorites", nil})
			fmt.Println(fmt.Sprintf("channelMetaInfo.Playlists: %+v", channelMetaInfo.Playlists))
		}

		// Fill SubscriberCount
		if channelMetaInfo.SubscriberCount == 0 {
			channelMetaInfo.SubscriberCount = firstItem.Statistics.SubscriberCount
			fmt.Println(fmt.Sprintf("channelMetaInfo.SubscriberCount: %d", channelMetaInfo.SubscriberCount))
		}

		// Fill ViewCount
		if channelMetaInfo.ViewCount == 0 {
			channelMetaInfo.ViewCount = firstItem.Statistics.ViewCount
			fmt.Println(fmt.Sprintf("channelMetaInfo.ViewCount: %d", channelMetaInfo.ViewCount))
		}
		fmt.Println(fmt.Sprintf("Filling complete. Result: %+v", channelMetaInfo))
		fmt.Println("Sending result to outChannel...")
		outChannel <- channelMetaInfo
		fmt.Println("Result successfully sent to outChannel")
		fmt.Println("Ending goroutine in GetChannelOverview")
	}()
	fmt.Println("End GetChannelOverview. Returning outChannel.")
	return outChannel
}

// asdasdasd
func GetVideoIDsOverview(service *youtube.Service, inChannel chan ChannelMetaInfo) <-chan ChannelMetaInfo {
	fmt.Println("Begin GetVideoIDsOverview")
	outChannel := make(chan ChannelMetaInfo)

	go func() {
		fmt.Println(fmt.Sprintf("Starting goroutine in GetVideoIDsOverview"))
		channelMetaInfo := <-inChannel
		//fmt.Println(fmt.Sprintf("Input channelMetaInfo: %+v", channelMetaInfo))
		call := service.PlaylistItems.List("contentDetails,snippet")
		call = call.PlaylistId(channelMetaInfo.Playlists[0].PlaylistID)

		response, responseError := call.Do()
		HandleError(responseError, "GetChannelOverview Response error!")

		firstItem := response.Items[0]
		fmt.Println(firstItem)
		// Fill ChannelID
		fmt.Println(fmt.Sprintf("Filling complete. Result: %+v", channelMetaInfo))
		fmt.Println("Sending result to outChannel...")
		outChannel <- channelMetaInfo
		fmt.Println("Result successfully sent to outChannel")
		fmt.Println("Ending goroutine in GetChannelOverview")
	}()
	fmt.Println("End GetChannelOverview. Returning outChannel.")
	return outChannel
}

// Exfoliator exfoliates
func Exfoliator(service *youtube.Service, channelMetaInfo ChannelMetaInfo) ChannelMetaInfo {
	initialInChannel := make(chan ChannelMetaInfo)

	go func() {
		initialInChannel <- channelMetaInfo
	}()
	getChannelOverviewOutChannel := GetChannelOverview(service, initialInChannel)
	timeout := time.After(5 * time.Second)
	// time.Sleep(time.Second)
	select {
	case channelMetaInfo = <-getChannelOverviewOutChannel:
		fmt.Println(fmt.Sprintf("Got %+v from getChannelOverviewOutChannel", channelMetaInfo))
	case <-timeout:
		fmt.Println("Request timed out...")
		return channelMetaInfo
	}

	return channelMetaInfo
}
