package service

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"time"

	youtube "google.golang.org/api/youtube/v3"
)

var (
	// Web1 Replicae
	PlaylistID1 = MetaSearch("PlaylistID1")
	// Web2 Replicae
	PlaylistID2 = MetaSearch("PlaylistID2")
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

func HandleError(errorToHandle error, errorMessage string) {
	if errorMessage == "" {
		errorMessage = "Error making API call"
	}
	if errorToHandle != nil {
		log.Fatalf(errorMessage+": %v", errorToHandle.Error())
	}
}

func getPlaylistIDByChannelIDOrNameAndPlaylistName(service *youtube.Service, query Query) Result {
	call := service.Channels.List("contentDetails")

	if query.ChannelID != "" {
		call = call.Id(query.ChannelID)
	} else if query.CustomURL != "" {
		call = call.ForUsername(query.CustomURL)
	}

	response, responseError := call.Do()
	HandleError(responseError, "Response error!")
	firstItem := response.Items[0]
	if query.PlaylistName == "uploads" {
		return Result{firstItem.ContentDetails.RelatedPlaylists.Uploads, nil}
	} else if query.PlaylistName == "favorites" {
		return Result{firstItem.ContentDetails.RelatedPlaylists.Favorites, nil}
	}

	return Result{"", errors.New("Unknown playlist. Available playlists are: uploads, favorites")}
}

func MetaSearch(searchName string) Search {
	return func(service *youtube.Service, query Query) Result {
		getPlaylistIDByChannelIDOrNameAndPlaylistName(service, query)
		time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
		return Result{fmt.Sprintf("%s result for %q\n", searchName, query.ChannelID), nil}
	}
}

func FirstResponder(service *youtube.Service, query Query, replicas ...Search) Result {
	resultChannel := make(chan Result)
	searchReplica := func(index int) { resultChannel <- replicas[index](service, query) }
	for replicaIndex := range replicas {
		go searchReplica(replicaIndex)
	}
	return <-resultChannel
}

func Exfoliator(service *youtube.Service, query Query) (resultsFromChannel []Result) {
	resultChannel := make(chan Result)

	go func() { resultChannel <- FirstResponder(service, query, PlaylistID1, PlaylistID2) }()
	/*
		go func() { resultChannel <- FirstResponder(query, Image1, Image2) }()
		go func() { resultChannel <- FirstResponder(query, Video1, Video2) }()
	*/
	timeout := time.After(10 * time.Second)
	for i := 0; i < 3; i++ {
		select {
		case resultFromChannel := <-resultChannel:
			resultsFromChannel = append(resultsFromChannel, resultFromChannel)
		case <-timeout:
			fmt.Println("timed out")
			return
		}
	}
	return
}
