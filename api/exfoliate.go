package api

import (
	"fmt"
	"log"
	"math/rand"
	"strings"
	"time"

	"google.golang.org/api/youtube/v3"
)

var (
	// SleepTime is the time of sleep between each things
	SleepTime = 10
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
			log.Println("<-- (1/5): GetChannelOverview")
			if channelMetaInfo.NextOperation == GetChannelOverviewOperation {
				log.Println("<-> (1/5): Working on GetChannelOverview")
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

				channelMetaInfo.NextOperation = GetVideoIDsOverviewOperation
				monoChannel <- channelMetaInfo
				log.Println("--> (1/5): GetChannelOverview")
			} else {
				log.Println("x-x (1/5): NOT Working on GetChannelOverview")
				monoChannel <- channelMetaInfo
				log.Println("--> (1/5): GetChannelOverview")
			}
			time.Sleep(time.Duration(rand.Intn(5*SleepTime)) * time.Millisecond)
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
			log.Println("<-- (2/5): GetVideoIDsOverview")
			if channelMetaInfo.NextOperation == GetVideoIDsOverviewOperation {
				log.Println("<-> (2/5): Working on GetVideoIDsOverview")
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
				channelMetaInfo.NextOperation = GetCommentsOverviewOperation
				monoChannel <- channelMetaInfo
				log.Println("--> (2/5): GetVideoIDsOverview")
			} else {
				log.Println("x-x (2/5): NOT Working on GetVideoIDsOverview")
				monoChannel <- channelMetaInfo
				log.Println("--> (2/5): GetVideoIDsOverview")
			}
			time.Sleep(time.Duration(rand.Intn(2*SleepTime)) * time.Millisecond)
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
			log.Println("<-- (3/5): GetCommentsOverview")
			if channelMetaInfo.NextOperation == GetCommentsOverviewOperation {
				log.Println("<-> (3/5): Working on GetCommentsOverview")
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
						channelMetaInfo.NextOperation = GetObviouslyRelatedChannelsOverviewOperation
						monoChannel <- channelMetaInfo
						log.Println("--> (3/5): GetCommentsOverview")
					}(i, video)
				}
			} else {
				log.Println("x-x (3/5): NOT Working on GetCommentsOverview")
				monoChannel <- channelMetaInfo
				log.Println("--> (3/5): GetCommentsOverview")
			}
			time.Sleep(time.Duration(rand.Intn(2*SleepTime)) * time.Millisecond)
		}
	}()
}

// GetObviouslyRelatedChannelsOverview gets the related channels for a YouTube channel
func GetObviouslyRelatedChannelsOverview(service *youtube.Service, monoChannel chan ChannelMetaInfo, lastButNotLeastChannel chan ChannelMetaInfo) {
	go func() {
		log.Println("<<<<<Begin GetObviouslyRelatedChannelsOverview Go Routine")
		defer log.Println("End GetObviouslyRelatedChannelsOverview Go Routine>>>>>")
		for {
			channelMetaInfo := <-monoChannel
			log.Println("<-- (4/5): GetObviouslyRelatedChannelsOverview")
			if channelMetaInfo.NextOperation == GetObviouslyRelatedChannelsOverviewOperation {
				log.Println("<-> (4/5): Working on GetObviouslyRelatedChannelsOverview")
				for i, commentatorChannelID := range channelMetaInfo.CommentAuthorChannelIDs {
					go func(index int, inputCommentatorChannelID string) {
						Printfln("<-> (4/5): (1/3) Begin service.Channels.List for ChannelID: %v", inputCommentatorChannelID)
						getChannelCall := service.Channels.List("snippet,contentDetails").Id(inputCommentatorChannelID)
						getChannelResponse, getChannelResponseError := getChannelCall.Do()
						log.Println(fmt.Sprintf("<-> (4/5): (1/3) End service.Channels.List with result: %v error: %v", getChannelResponse, getChannelResponseError))
						if getChannelResponseError != nil {
							formattdErrorMessage := GetFormattedErrorMessage(getChannelResponseError, fmt.Sprintf("GetObviouslyRelatedChannelsOverview#%d Response error!", index))
							if formattdErrorMessage != "" {
								log.Println(formattdErrorMessage)
							}
							return
						}
						Printfln("<-> (4/5): (2/3) Begin service.PlaylistItems.List for PlaylistID: %v", getChannelResponse.Items[0].ContentDetails.RelatedPlaylists.Favorites)
						getPlaylistItemsCall := service.PlaylistItems.List("contentDetails").PlaylistId(getChannelResponse.Items[0].ContentDetails.RelatedPlaylists.Favorites).MaxResults(50)
						getPlaylistItemsResponse, getPlaylistItemsResponseError := getPlaylistItemsCall.Do()
						log.Println("->> (4/5)!!: GetObviouslyRelatedChannelsOverview")
						log.Println(fmt.Sprintf("<-> (4/5): (2/3) End service.PlaylistItems.List with result: %v error: %v", getPlaylistItemsResponse, getPlaylistItemsResponseError))
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

						Printfln("<-> (4/5): (3/3) Begin service.Videos.List for VideoIDs: %v", strings.Join(favoritedVideoIDs, ","))
						getRelatedChannelCall := service.Videos.List("snippet").Id(strings.Join(favoritedVideoIDs, ","))
						getRelatedChannelResponse, getRelatedChannelResponseError := getRelatedChannelCall.Do()
						log.Println(fmt.Sprintf("<-> (4/5): (3/3) End service.Videos.List with result: %v error: %v", getRelatedChannelResponse, getRelatedChannelResponseError))

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
						lastButNotLeastChannel <- channelMetaInfo
						log.Println("--> (4/5): GetObviouslyRelatedChannelsOverview")

					}(i, commentatorChannelID)
				}
			} else {
				log.Println("x-x (4/5): NOT Working on GetObviouslyRelatedChannelsOverview")
				monoChannel <- channelMetaInfo
				log.Println("--> (4/5): GetObviouslyRelatedChannelsOverview")
			}
			time.Sleep(time.Duration(rand.Intn(SleepTime)) * time.Millisecond)
		}
	}()
}

// Exfoliator exfoliates
func Exfoliator(service *youtube.Service, channelMetaInfo ChannelMetaInfo) ChannelMetaInfo {
	monoChannel := make(chan ChannelMetaInfo)
	lastButNotLeastChannel := make(chan ChannelMetaInfo)
	channelMetaInfo.NextOperation = GetChannelOverviewOperation
	accumulatedMetaInfo := ChannelMetaInfo{}
	accumulatedMetaInfo.CustomURL = channelMetaInfo.CustomURL
	go func() {
		monoChannel <- channelMetaInfo
	}()
	go GetChannelOverview(service, monoChannel)
	go GetVideoIDsOverview(service, monoChannel)
	go GetCommentsOverview(service, monoChannel)
	go GetObviouslyRelatedChannelsOverview(service, monoChannel, lastButNotLeastChannel)

	timeout := time.After(5 * 60 * time.Second)
	// time.Sleep(time.Second)

	for {
		log.Println("<<<<<Begin Exfoliator Go Routine")
		select {
		case channelMetaInfo = <-lastButNotLeastChannel:
			log.Println("<-- (5/5): Exfoliator")
			log.Println("<-> (5/5): Working on Exfoliator ++++++")
			// evtl die anderen properties adden
			accumulatedMetaInfo.ObviouslyRelatedChannelIDs = append(accumulatedMetaInfo.ObviouslyRelatedChannelIDs, channelMetaInfo.ObviouslyRelatedChannelIDs...)
			log.Println("--> (5/5): Exfoliator")
			time.Sleep(time.Duration(rand.Intn(SleepTime)) * time.Millisecond)
		case <-timeout:
			log.Println("Request timed out (30 sec)")
			return accumulatedMetaInfo
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
}