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
		fmt.Println("Starting goroutine in GetChannelOverview")
		channelMetaInfo := <-inChannel
		call := service.Channels.List("contentDetails,snippet,statistics").ForUsername(channelMetaInfo.CustomURL)

		response, responseError := call.Do()
		if responseError != nil {
			formattdErrorMessage := GetFormattedErrorMessage(responseError, "GetChannelOverview Response error!")
			if formattdErrorMessage != "" {
				log.Println(formattdErrorMessage)
			}
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
		fmt.Println("Starting goroutine in GetVideoIDsOverview")
		channelMetaInfo := <-inChannel
		//Printfln("Input channelMetaInfo: %+v", channelMetaInfo)
		call := service.PlaylistItems.List("contentDetails,snippet").PlaylistId(channelMetaInfo.Playlists["uploads"].PlaylistID).MaxResults(50)

		response, responseError := call.Do()
		if responseError != nil {
			formattdErrorMessage := GetFormattedErrorMessage(responseError, "GetVideoIDsOverview Response error!")
			if formattdErrorMessage != "" {
				log.Println(formattdErrorMessage)
			}
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
		for i, video := range channelMetaInfo.Playlists["uploads"].PlaylistItems {
			go func(inputVideo *Video) {
				Printfln("Starting goroutine GetCommentsOverview#%d", i)

				//Printfln("Input channelMetaInfo: %+v", channelMetaInfo)
				call := service.CommentThreads.List("snippet").VideoId(inputVideo.VideoID)

				response, responseError := call.Do()

				if responseError != nil {
					formattdErrorMessage := GetFormattedErrorMessage(responseError, fmt.Sprintf("GetCommentsOverview#%d Response error! (videoId: %s)", i, inputVideo.VideoID))
					if formattdErrorMessage != "" {
						log.Println(formattdErrorMessage)
					}
					return
				}

				var comments []*Comment
				for _, item := range response.Items {
					comment := new(Comment)
					comment.CommentID = item.Snippet.TopLevelComment.Id
					comment.AuthorChannelID = item.Snippet.TopLevelComment.Snippet.AuthorChannelId.(map[string]interface{})["value"].(string)
					channelMetaInfo.CommentAuthorChannelIDs = append(channelMetaInfo.CommentAuthorChannelIDs, comment.AuthorChannelID)

					comments = append(comments, comment)
					Printfln("Appended comment: %v to video: %v", comment, video)
					Printfln("video.Comments now: %+v", video.Comments)
				}

				video.Comments = comments
				Printfln("GetCommentsOverview#%d filling complete. Result: %+v", i, channelMetaInfo)
				fmt.Println("Sending result to getCommentsOverviewOutChannel...")
				outChannel <- channelMetaInfo
				fmt.Println("Result successfully sent to getCommentsOverviewOutChannel")
			}(video)
			Printfln("Ending goroutine in GetCommentsOverview#%d", i)
		}
	}()
	fmt.Println("End GetCommentsOverview. Returning getCommentsOverviewOutChannel.")
	return outChannel
}

// GetObviouslyRelatedChannelsOverview foo
func GetObviouslyRelatedChannelsOverview(service *youtube.Service, inChannel <-chan ChannelMetaInfo) <-chan ChannelMetaInfo {
	fmt.Println("Begin GetObviouslyRelatedChannelsOverview")
	outChannel := make(chan ChannelMetaInfo)
	channelMetaInfo := <-inChannel
	go func() {
		for i, commentatorChannelID := range channelMetaInfo.CommentAuthorChannelIDs {
			go func(inputCommentatorChannelID string) {
				Printfln("Starting goroutine GetObviouslyRelatedChannelsOverview#%d", i)

				//Printfln("Input channelMetaInfo: %+v", channelMetaInfo)
				getChannelCall := service.Channels.List("snippet,contentDetails").Id(inputCommentatorChannelID)
				getChannelResponse, getChannelResponseError := getChannelCall.Do()

				if getChannelResponseError != nil {
					formattdErrorMessage := GetFormattedErrorMessage(getChannelResponseError, fmt.Sprintf("GetObviouslyRelatedChannelsOverview#%d Response error!", i))
					if formattdErrorMessage != "" {
						log.Println(formattdErrorMessage)
					}
					return
				}

				getPlaylistItemsCall := service.PlaylistItems.List("snippet,contentDetails").PlaylistId(getChannelResponse.Items[0].ContentDetails.RelatedPlaylists.Favorites)
				getPlaylistItemsResponse, getPlaylistItemsResponseError := getPlaylistItemsCall.Do()

				if getPlaylistItemsResponseError != nil {
					formattedErrorMesage := GetFormattedErrorMessage(getChannelResponseError, fmt.Sprintf("GetObviouslyRelatedChannelsOverview#%d Response error!", i))
					if formattedErrorMesage != "" {
						log.Println(formattedErrorMesage)
					}
					return
				}

				// ObviouslyRelatedChannelIDs

				channelMetaInfo.ObviouslyRelatedChannelIDs = comments
				Printfln("GetCommentsOverview#%d filling complete. Result: %+v", i, channelMetaInfo)
				fmt.Println("Sending result to getCommentsOverviewOutChannel...")
				outChannel <- channelMetaInfo
				fmt.Println("Result successfully sent to getCommentsOverviewOutChannel")
			}(commentatorChannelID)
			Printfln("Ending goroutine in GetCommentsOverview#%d", i)
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
	getObviouslyRelatedChannelsOverviewChannel := GetObviouslyRelatedChannelsOverview(service, getCommentsOverviewOutChannel)
	timeout := time.After(10 * time.Second)
	// time.Sleep(time.Second)
	select {
	case channelMetaInfo = <-getObviouslyRelatedChannelsOverviewChannel:
		Printfln("Got %+v from getCommentsOverviewOutChannel", channelMetaInfo)
	case <-timeout:
		fmt.Println("Request timed out...")
		return channelMetaInfo
	}

	return channelMetaInfo
}
