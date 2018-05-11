package api

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
		call := service.Channels.List("contentDetails,snippet,statistics").ForUsername(channelMetaInfo.CustomURL)

		response, responseError := call.Do()
		HandleError(responseError, "GetChannelOverview Response error!")

		firstItem := response.Items[0]

		// Fill ChannelID
		if channelMetaInfo.ChannelID == "" {
			channelMetaInfo.ChannelID = firstItem.Id
			//fmt.Println(fmt.Sprintf("channelMetaInfo.ChannelID: %s", channelMetaInfo.ChannelID))
		}

		// Fill CustomURL
		if channelMetaInfo.CustomURL == "" {
			channelMetaInfo.CustomURL = firstItem.Snippet.CustomUrl
			//fmt.Println(fmt.Sprintf("channelMetaInfo.CustomURL: %s", channelMetaInfo.CustomURL))
		}

		// Fill ChannelName
		if channelMetaInfo.ChannelName == "" {
			channelMetaInfo.ChannelName = firstItem.Snippet.Title
			//fmt.Println(fmt.Sprintf("channelMetaInfo.ChannelName: %s", channelMetaInfo.ChannelName))
		}

		// Fill Playlists
		if channelMetaInfo.Playlists == nil {
			channelMetaInfo.Playlists = make(map[string]*Playlist)
			var videos []*Video
			channelMetaInfo.Playlists["uploads"] = &Playlist{firstItem.ContentDetails.RelatedPlaylists.Uploads, videos}
			channelMetaInfo.Playlists["favorites"] = &Playlist{firstItem.ContentDetails.RelatedPlaylists.Favorites, videos}
			//fmt.Println(fmt.Sprintf("channelMetaInfo.Playlists: %+v", channelMetaInfo.Playlists))
		}

		// Fill SubscriberCount
		if channelMetaInfo.SubscriberCount == 0 {
			channelMetaInfo.SubscriberCount = firstItem.Statistics.SubscriberCount
			//fmt.Println(fmt.Sprintf("channelMetaInfo.SubscriberCount: %d", channelMetaInfo.SubscriberCount))
		}

		// Fill ViewCount
		if channelMetaInfo.ViewCount == 0 {
			channelMetaInfo.ViewCount = firstItem.Statistics.ViewCount
			//fmt.Println(fmt.Sprintf("channelMetaInfo.ViewCount: %d", channelMetaInfo.ViewCount))
		}
		fmt.Println(fmt.Sprintf("GetChannelOverview filling complete. Result: %+v", channelMetaInfo))
		fmt.Println("Sending result to getChannelOverviewOutChannel...")
		outChannel <- channelMetaInfo
		fmt.Println("Result successfully sent to getChannelOverviewOutChannel")
		fmt.Println("Ending goroutine in GetChannelOverview")
	}()
	fmt.Println("End GetChannelOverview. Returning getChannelOverviewOutChannel.")
	return outChannel
}

// GetVideoIDsOverview bla
func GetVideoIDsOverview(service *youtube.Service, inChannel <-chan ChannelMetaInfo) <-chan ChannelMetaInfo {
	fmt.Println("Begin GetVideoIDsOverview")
	outChannel := make(chan ChannelMetaInfo)

	go func() {
		fmt.Println(fmt.Sprintf("Starting goroutine in GetVideoIDsOverview"))
		channelMetaInfo := <-inChannel
		//fmt.Println(fmt.Sprintf("Input channelMetaInfo: %+v", channelMetaInfo))
		call := service.PlaylistItems.List("contentDetails,snippet").PlaylistId(channelMetaInfo.Playlists["uploads"].PlaylistID).MaxResults(10)

		response, responseError := call.Do()
		HandleError(responseError, "GetVideoIDsOverview Response error!")

		var videos []*Video
		for _, item := range response.Items {
			video := &Video{VideoID: item.Id}
			videos = append(videos, video)
			//fmt.Println(fmt.Sprintf("Appended video %s to playlist uploads", video))
		}
		uploadPlaylist := channelMetaInfo.Playlists["uploads"]
		uploadPlaylist.PlaylistItems = videos
		//fmt.Println(fmt.Sprintf("uploadPlaylist.PlaylistItems now: %+v", uploadPlaylist.PlaylistItems))
		//fmt.Println(fmt.Sprintf("GetVideoIDsOverview filling complete. Result: %+v", channelMetaInfo))
		fmt.Println("Sending result to getVideoIDsOverviewOutChannel...")
		outChannel <- channelMetaInfo
		fmt.Println("Result successfully sent to getVideoIDsOverviewOutChannel")
		fmt.Println("Ending goroutine in GetVideoIDsOverview")
	}()
	fmt.Println("End GetVideoIDsOverview. Returning getVideoIDsOverviewOutChannel.")
	return outChannel
}

// GetCommentsOverview foo
func GetCommentsOverview(service *youtube.Service, inChannel <-chan ChannelMetaInfo) <-chan ChannelMetaInfo {
	fmt.Println("Begin GetCommentsOverview")
	outChannel := make(chan ChannelMetaInfo)
	channelMetaInfo := <-inChannel
	go func() {
		for _, video := range channelMetaInfo.Playlists["uploads"].PlaylistItems {
			go func(inputVideo *Video) {
				fmt.Println(fmt.Sprintf("Starting goroutine in GetCommentsOverview"))

				//fmt.Println(fmt.Sprintf("Input channelMetaInfo: %+v", channelMetaInfo))
				call := service.CommentThreads.List("snippet").VideoId(video.VideoID)

				response, responseError := call.Do()
				HandleError(responseError, "GetCommentsOverview Response error!")

				var comments []*Comment
				fmt.Println(fmt.Sprintf("!!!!!!!!!!!!!!!!!!!!!!"))
				for _, item := range response.Items {
					comment := &Comment{CommentID: item.Snippet.TopLevelComment.Id, AuthorChannelID: item.Snippet.TopLevelComment.Snippet.AuthorChannelId.(string)}
					comments = append(comments, comment)
					fmt.Println(fmt.Sprintf("Appended comment: %v to video: %v", comment, video))
					fmt.Println(fmt.Sprintf("video.Comments now: %+v", video.Comments))
				}

				video.Comments = comments
				fmt.Println(fmt.Sprintf("GetCommentsOverview filling complete. Result: %+v", channelMetaInfo))
				fmt.Println("Sending result to getCommentsOverviewOutChannel...")
				outChannel <- channelMetaInfo
				fmt.Println("Result successfully sent to getCommentsOverviewOutChannel")
			}(video)
			fmt.Println("Ending goroutine in GetCommentsOverview")
		}
	}()
	fmt.Println("End GetCommentsOverview. Returning getCommentsOverviewOutChannel.")
	return outChannel
}

// Exfoliator exfoliates
func Exfoliator(service *youtube.Service, channelMetaInfo ChannelMetaInfo) ChannelMetaInfo {
	initialInChannel := make(chan ChannelMetaInfo)

	go func() {
		initialInChannel <- channelMetaInfo
	}()
	getChannelOverviewOutChannel := GetChannelOverview(service, initialInChannel)
	getVideoIDsOverviewOutChannel := GetVideoIDsOverview(service, getChannelOverviewOutChannel)
	getCommentsOverviewOutChannel := GetCommentsOverview(service, getVideoIDsOverviewOutChannel)
	timeout := time.After(20 * time.Second)
	// time.Sleep(time.Second)
	select {
	case channelMetaInfo = <-getCommentsOverviewOutChannel:
		fmt.Println(fmt.Sprintf("Got %+v from getCommentsOverviewOutChannel", channelMetaInfo))
	case <-timeout:
		fmt.Println("Request timed out...")
		return channelMetaInfo
	}

	return channelMetaInfo
}
