package service

import (
	"errors"
	"fmt"
	"log"
	"time"

	"google.golang.org/api/youtube/v3"
)

// Result datatype
type Result struct {
	ResultValue string
	ResultError error
}

// Query datatype
type Query struct {
	ChannelID    string
	CustomURL    string
	PlaylistName string
}

// Search function declaration
type Search func(serviceservice *youtube.Service, query Query) Result

// HandleError handles errors
func HandleError(errorToHandle error, errorMessage string) {
	if errorMessage == "" {
		errorMessage = "Error making API call"
	}
	if errorToHandle != nil {
		log.Fatalf(errorMessage+": %v", errorToHandle.Error())
	}
}

// GetPlaylistIDByQueryParameters returns the id of the playlist for the given
// search query. You can either query by channelID or customURL.
func GetPlaylistIDByQueryParameters(service *youtube.Service, query Query) Result {
	call := service.Channels.List("contentDetails")

	if query.ChannelID != "" {
		call = call.Id(query.ChannelID)
	} else if query.CustomURL != "" {
		call = call.ForUsername(query.CustomURL)
	}

	response, responseError := call.Do()
	HandleError(responseError, "Response error!")

	result := Result{}
	firstItem := response.Items[0]
	if query.PlaylistName == "uploads" {
		result.ResultValue = firstItem.ContentDetails.RelatedPlaylists.Uploads
		return result
	} else if query.PlaylistName == "favorites" {
		result.ResultValue = firstItem.ContentDetails.RelatedPlaylists.Favorites
		return result
	}
	result.ResultError = errors.New("Unknown playlist. Available playlists are: uploads, favorites")
	return result
}

// MetaSearch searches meta
func MetaSearch(searchName string) Search {
	return func(service *youtube.Service, query Query) Result {
		fmt.Println(fmt.Sprintf("%v starts running with query: %v", searchName, query))
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
