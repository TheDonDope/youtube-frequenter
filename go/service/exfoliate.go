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
	return func(service *youtube.Service, query Query) Result {
		fmt.Println(fmt.Sprintf("%v starts running with query: %v", searchName, query))
		queryMethod := query.QueryMethod

		result := GetPlaylistIDByQueryParameters(service, query)

		return result
	}
}

// FirstResponder returns the first replica search result, whichever was the faster
func FirstResponder(service *youtube.Service, query Query, replicas ...Search) Result {
	resultChannel := make(chan Result)
	searchReplica := func(index int) { resultChannel <- replicas[index](service, query) }
	for replicaIndex := range replicas {
		go searchReplica(replicaIndex)
	}
	return <-resultChannel
}

// Exfoliator exfoliates
func Exfoliator(service *youtube.Service, query Query) (resultsFromChannel []Result) {
	resultChannel := make(chan Result)

	go func() {
		resultChannel <- FirstResponder(service, query, MetaSearch("PlaylistID1"), MetaSearch("PlaylistID2"))
	}()

	timeout := time.After(10 * time.Second)
	for i := 0; i < 1; i++ {
		select {
		case resultFromChannel := <-resultChannel:
			resultsFromChannel = append(resultsFromChannel, resultFromChannel)
		case <-timeout:
			fmt.Println("timed out")
			return
		}
	}
	return resultsFromChannel
}
