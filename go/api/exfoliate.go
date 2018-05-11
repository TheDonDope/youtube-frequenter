package api

import (
	"fmt"
	"log"
	"strings"
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
func GetChannelOverview(service *youtube.Service, monoChannel chan ChannelMetaInfo) {
	go func() {
		log.Println("<<<<<Begin GetChannelOverview Go Routine")
		defer log.Println("End GetChannelOverview Go Routine>>>>>")
		for {
			channelMetaInfo := <-monoChannel
			if channelMetaInfo.NextOperation == "GetChannelOverview" {

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
					channelMetaInfo.Playlists = make(map[string]*Playlist)
					var videos []*Video
					channelMetaInfo.Playlists["uploads"] = &Playlist{firstItem.ContentDetails.RelatedPlaylists.Uploads, videos}
					channelMetaInfo.Playlists["favorites"] = &Playlist{firstItem.ContentDetails.RelatedPlaylists.Favorites, videos}
				}

				// Fill SubscriberCount
				if channelMetaInfo.SubscriberCount == 0 {
					channelMetaInfo.SubscriberCount = firstItem.Statistics.SubscriberCount
				}

				// Fill ViewCount
				if channelMetaInfo.ViewCount == 0 {
					channelMetaInfo.ViewCount = firstItem.Statistics.ViewCount
				}

				channelMetaInfo.NextOperation = "GetVideoIDsOverview"
			}
			monoChannel <- channelMetaInfo
		}
	}()

}

// GetVideoIDsOverview gets all videos.
func GetVideoIDsOverview(service *youtube.Service, monoChannel chan ChannelMetaInfo) {
	go func() {
		log.Println("<<<<<Begin GetVideoIDsOverview Go Routine")
		defer log.Println("End GetVideoIDsOverview Go Routine>>>>>")
		for {
			channelMetaInfo := <-monoChannel
			if channelMetaInfo.NextOperation == "GetVideoIDsOverview" {
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
				}

				uploadPlaylist := channelMetaInfo.Playlists["uploads"]
				uploadPlaylist.PlaylistItems = videos
				channelMetaInfo.NextOperation = "GetCommentsOverview"
			}
			monoChannel <- channelMetaInfo
		}
	}()
}

// GetCommentsOverview foo
func GetCommentsOverview(service *youtube.Service, monoChannel chan ChannelMetaInfo) {
	go func() {
		log.Println("<<<<<Begin GetCommentsOverview Go Routine")
		defer log.Println("End GetCommentsOverview Go Routine>>>>>")
		for {
			channelMetaInfo := <-monoChannel
			if channelMetaInfo.NextOperation == "GetCommentsOverview" {

				for i, video := range channelMetaInfo.Playlists["uploads"].PlaylistItems {
					go func(index int, inputVideo *Video) {
						call := service.CommentThreads.List("snippet").VideoId(inputVideo.VideoID).MaxResults(100)

						response, responseError := call.Do()

						if responseError != nil {
							formattdErrorMessage := GetFormattedErrorMessage(responseError, fmt.Sprintf("GetCommentsOverview#%d Response error! (videoId: %s)", index, inputVideo.VideoID))
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
						}

						inputVideo.Comments = comments
						channelMetaInfo.NextOperation = "GetObviouslyRelatedChannelsOverview"
					}(i, video)
				}
			}
			monoChannel <- channelMetaInfo
		}
	}()
}

// GetObviouslyRelatedChannelsOverview gets the related channels for a YouTube channel
func GetObviouslyRelatedChannelsOverview(service *youtube.Service, monoChannel chan ChannelMetaInfo) {
	go func() {
		log.Println("<<<<<Begin GetObviouslyRelatedChannelsOverview Go Routine")
		defer log.Println("End GetObviouslyRelatedChannelsOverview Go Routine>>>>>")
		for {
			channelMetaInfo := <-monoChannel
			if channelMetaInfo.NextOperation == "GetObviouslyRelatedChannelsOverview" {
				for i, commentatorChannelID := range channelMetaInfo.CommentAuthorChannelIDs {
					go func(index int, inputCommentatorChannelID string) {
						getChannelCall := service.Channels.List("snippet,contentDetails").Id(inputCommentatorChannelID)
						getChannelResponse, getChannelResponseError := getChannelCall.Do()

						if getChannelResponseError != nil {
							formattdErrorMessage := GetFormattedErrorMessage(getChannelResponseError, fmt.Sprintf("GetObviouslyRelatedChannelsOverview#%d Response error!", index))
							if formattdErrorMessage != "" {
								log.Println(formattdErrorMessage)
							}
							return
						}

						getPlaylistItemsCall := service.PlaylistItems.List("contentDetails").PlaylistId(getChannelResponse.Items[0].ContentDetails.RelatedPlaylists.Favorites).MaxResults(50)
						getPlaylistItemsResponse, getPlaylistItemsResponseError := getPlaylistItemsCall.Do()

						if getPlaylistItemsResponseError != nil {
							formattedErrorMesage := GetFormattedErrorMessage(getChannelResponseError, fmt.Sprintf("GetObviouslyRelatedChannelsOverview#%d Response error!", index))
							if formattedErrorMesage != "" {
								log.Println(formattedErrorMesage)
							}
							return
						}

						var favoritedVideoIDs []string
						for _, item := range getPlaylistItemsResponse.Items {
							favoritedVideoIDs = append(favoritedVideoIDs, item.ContentDetails.VideoId)
						}

						getRelatedChannelCall := service.Videos.List("snippet").Id(strings.Join(favoritedVideoIDs, ","))
						getRelatedChannelResponse, getRelatedChannelResponseError := getRelatedChannelCall.Do()

						if getRelatedChannelResponseError != nil {
							formattedErrorMesage := GetFormattedErrorMessage(getRelatedChannelResponseError, fmt.Sprintf("GetObviouslyRelatedChannelsOverview#%d Response error!", index))
							if formattedErrorMesage != "" {
								log.Println(formattedErrorMesage)
							}
							return
						}

						var obviouslyRelatedChannelNames []string
						// ObviouslyRelatedChannelIDs
						for _, item := range getRelatedChannelResponse.Items {
							obviouslyRelatedChannelNames = append(obviouslyRelatedChannelNames, item.Snippet.ChannelTitle)
						}

						channelMetaInfo.ObviouslyRelatedChannelIDs = obviouslyRelatedChannelNames
						channelMetaInfo.NextOperation = "None"

					}(i, commentatorChannelID)
				}
			}
			monoChannel <- channelMetaInfo
		}
	}()
}

// Exfoliator exfoliates
func Exfoliator(service *youtube.Service, channelMetaInfo ChannelMetaInfo) ChannelMetaInfo {
	monoChannel := make(chan ChannelMetaInfo)
	channelMetaInfo.NextOperation = "GetChannelOverview"
	go func() {
		monoChannel <- channelMetaInfo
	}()
	GetChannelOverview(service, monoChannel)
	GetVideoIDsOverview(service, monoChannel)
	GetCommentsOverview(service, monoChannel)
	GetObviouslyRelatedChannelsOverview(service, monoChannel)

	timeout := time.After(20 * time.Second)
	// time.Sleep(time.Second)

	for {
		select {
		case channelMetaInfo = <-monoChannel:
			if channelMetaInfo.NextOperation == "None" {
				Printfln("Done #relatedChannels=%v ", len(channelMetaInfo.ObviouslyRelatedChannelIDs))
				return channelMetaInfo
			}
			monoChannel <- channelMetaInfo
		case <-timeout:
			log.Println("Initial Request timed out (10sec)")
			return channelMetaInfo
		}
	}

	// // check for unused channelMetaInfo
	// timeoutAfter := time.After(5 * time.Second)
	// // time.Sleep(time.Second)
	// for i := 0; i < 10000; i++ {
	// 	select {
	// 	case channelMetaInfo = <-getChannelOverviewOutChannel:
	// 		Printfln("!!!getChannelOverviewOutChannel%s", "")
	// 	case channelMetaInfo = <-getVideoIDsOverviewOutChannel:
	// 		Printfln("!!!getVideoIDsOverviewOutChannel%s", "")
	// 	case channelMetaInfo = <-getCommentsOverviewOutChannel:
	// 		Printfln("!!getCommentsOverviewOutChannel%s", "")
	// 	case channelMetaInfo = <-getObviouslyRelatedChannelsOverviewChannel:
	// 		Printfln("After #relatedChannels=%v ", len(channelMetaInfo.ObviouslyRelatedChannelIDs))
	// 	case <-timeoutAfter:
	// 		Printfln("After Request timed out (5sec): %d", i)
	// 		return channelMetaInfo
	// 	}
	// }
	return channelMetaInfo
}
