package api

import (
	"fmt"
	"log"
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
		Printfln("Starting goroutine in GetChannelOverview")
		channelMetaInfo := <-inChannel
		call := service.Channels.List("contentDetails,snippet,statistics").ForUsername(channelMetaInfo.CustomURL)

		response, responseError := call.Do()
		formattdErrorMessage := GetFormattedErrorMessage(responseError, "GetChannelOverview Response error!")
		if formattdErrorMessage != "" {
			log.Fatal(formattdErrorMessage)
		}

		firstItem := response.Items[0]

		// Fill ChannelID
		if channelMetaInfo.ChannelID == "" {
			channelMetaInfo.ChannelID = firstItem.Id
			//Printfln("channelMetaInfo.ChannelID: %s", channelMetaInfo.ChannelID)
		}

		// Fill CustomURL
		if channelMetaInfo.CustomURL == "" {
			channelMetaInfo.CustomURL = firstItem.Snippet.CustomUrl
			//Printfln("channelMetaInfo.CustomURL: %s", channelMetaInfo.CustomURL)
		}

		// Fill ChannelName
		if channelMetaInfo.ChannelName == "" {
			channelMetaInfo.ChannelName = firstItem.Snippet.Title
			//Printfln("channelMetaInfo.ChannelName: %s", channelMetaInfo.ChannelName)
		}

		// Fill Playlists
		if channelMetaInfo.Playlists == nil {
			channelMetaInfo.Playlists = make(map[string]*Playlist)
			var videos []*Video
			channelMetaInfo.Playlists["uploads"] = &Playlist{firstItem.ContentDetails.RelatedPlaylists.Uploads, videos}
			channelMetaInfo.Playlists["favorites"] = &Playlist{firstItem.ContentDetails.RelatedPlaylists.Favorites, videos}
			//Printfln("channelMetaInfo.Playlists: %+v", channelMetaInfo.Playlists)
		}

		// Fill SubscriberCount
		if channelMetaInfo.SubscriberCount == 0 {
			channelMetaInfo.SubscriberCount = firstItem.Statistics.SubscriberCount
			//Printfln("channelMetaInfo.SubscriberCount: %d", channelMetaInfo.SubscriberCount)
		}

		// Fill ViewCount
		if channelMetaInfo.ViewCount == 0 {
			channelMetaInfo.ViewCount = firstItem.Statistics.ViewCount
			//Printfln("channelMetaInfo.ViewCount: %d", channelMetaInfo.ViewCount)
		}
		Printfln("GetChannelOverview filling complete. Result: %+v", channelMetaInfo)
		Printfln("Sending result to getChannelOverviewOutChannel...")
		outChannel <- channelMetaInfo
		Printfln("Result successfully sent to getChannelOverviewOutChannel")
		Printfln("Ending goroutine in GetChannelOverview")
	}()
	Printfln("End GetChannelOverview. Returning getChannelOverviewOutChannel.")
	return outChannel
}

// GetVideoIDsOverview bla
func GetVideoIDsOverview(service *youtube.Service, inChannel <-chan ChannelMetaInfo) <-chan ChannelMetaInfo {
	Printfln("Begin GetVideoIDsOverview")
	outChannel := make(chan ChannelMetaInfo)

	go func() {
		fmt.Println(fmt.Sprintf("Starting goroutine in GetVideoIDsOverview"))
		channelMetaInfo := <-inChannel
		//Printfln("Input channelMetaInfo: %+v", channelMetaInfo)
		call := service.PlaylistItems.List("contentDetails,snippet").PlaylistId(channelMetaInfo.Playlists["uploads"].PlaylistID).MaxResults(10)

		response, responseError := call.Do()
		formattdErrorMessage := GetFormattedErrorMessage(responseError, "GetVideoIDsOverview Response error!")
		if formattdErrorMessage != "" {
			log.Fatal(formattdErrorMessage)
		}

		var videos []*Video
		for _, item := range response.Items {
			video := &Video{VideoID: item.Snippet.ResourceId.VideoId}
			videos = append(videos, video)
			//Printfln("Appended video %s to playlist uploads", video)
		}

		uploadPlaylist := channelMetaInfo.Playlists["uploads"]
		Printfln("playlistID: %s", uploadPlaylist.PlaylistID)
		uploadPlaylist.PlaylistItems = videos
		Printfln("PlaylistItems: %s", uploadPlaylist.PlaylistItems)

		//Printfln("uploadPlaylist.PlaylistItems now: %+v", uploadPlaylist.PlaylistItems)
		//Printfln("GetVideoIDsOverview filling complete. Result: %+v", channelMetaInfo)
		Printfln("Sending result to getVideoIDsOverviewOutChannel...")
		outChannel <- channelMetaInfo
		Printfln("Result successfully sent to getVideoIDsOverviewOutChannel")
		Printfln("Ending goroutine in GetVideoIDsOverview")
	}()
	Printfln("End GetVideoIDsOverview. Returning getVideoIDsOverviewOutChannel.")
	return outChannel
}

// GetCommentsOverview foo
func GetCommentsOverview(service *youtube.Service, inChannel <-chan ChannelMetaInfo) <-chan ChannelMetaInfo {
	Printfln("Begin GetCommentsOverview")
	outChannel := make(chan ChannelMetaInfo)
	channelMetaInfo := <-inChannel
	go func() {
		for i, video := range channelMetaInfo.Playlists["uploads"].PlaylistItems {
			go func(inputVideo *Video) {
				Printfln("Starting goroutine GetCommentsOverview#%d", i)

				//Printfln("Input channelMetaInfo: %+v", channelMetaInfo)
				call := service.CommentThreads.List("snippet").VideoId(video.VideoID)

				response, responseError := call.Do()

				formattdErrorMessage := GetFormattedErrorMessage(responseError, fmt.Sprintf("GetCommentsOverview#%d Response error!"+video.VideoID, i))
				if formattdErrorMessage != "" {
					log.Fatal(formattdErrorMessage)
				}

				var comments []*Comment
				for _, item := range response.Items {
					comment := new(Comment)
					comment.CommentID = item.Snippet.TopLevelComment.Id
					comment.AuthorChannelID = item.Snippet.TopLevelComment.Snippet.AuthorChannelId.(string)

					comments = append(comments, comment)
					Printfln("Appended comment: %v to video: %v", comment, video)
					Printfln("video.Comments now: %+v", video.Comments)
				}

				video.Comments = comments
				Printfln("GetCommentsOverview#%d filling complete. Result: %+v", i, channelMetaInfo)
				Printfln("Sending result to getCommentsOverviewOutChannel...")
				outChannel <- channelMetaInfo
				Printfln("Result successfully sent to getCommentsOverviewOutChannel")
			}(video)
			Printfln("Ending goroutine in GetCommentsOverview#%d", i)
		}
	}()
	Printfln("End GetCommentsOverview. Returning getCommentsOverviewOutChannel.")
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
		Printfln("Got %+v from getCommentsOverviewOutChannel", channelMetaInfo)
	case <-timeout:
		Printfln("Request timed out...")
		return channelMetaInfo
	}

	return channelMetaInfo
}
