package api

import (
	"fmt"
	"log"
	"strings"

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
			log.Println("<-- (1/5): Receiving into GetChannelOverview")
			if channelMetaInfo.NextOperation == GetChannelOverviewOperation {
				log.Println("<-> (1/5): Working in GetChannelOverview")
				response, responseError := ChannelsList(service, channelMetaInfo.ChannelID, channelMetaInfo.CustomURL)

				HandleError(responseError, "GetChannelOverview Response error!")

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
				log.Println("--> (1/5): Sending out from GetChannelOverview")
			} else {
				log.Println("x-x (1/5): NOT Working in GetChannelOverview")
				monoChannel <- channelMetaInfo
				log.Println("--> (1/5): Sendubg out from GetChannelOverview")
			}
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
			log.Println("<-- (2/5): Receiving into GetVideoIDsOverview")
			if channelMetaInfo.NextOperation == GetVideoIDsOverviewOperation {
				log.Println("<-> (2/5): Working in GetVideoIDsOverview")
				call := service.PlaylistItems.List("contentDetails,snippet").PlaylistId(channelMetaInfo.Playlists["uploads"].PlaylistID).MaxResults(Opts.MaxResultsUploadedVideos)
				response, responseError := call.Do()

				HandleError(responseError, "GetVideoIDsOverview Response error!")

				var videos []*Video
				for _, item := range response.Items {
					video := &Video{VideoID: item.Snippet.ResourceId.VideoId}
					videos = append(videos, video)
				}

				uploadPlaylist := channelMetaInfo.Playlists["uploads"]
				uploadPlaylist.PlaylistItems = videos
				channelMetaInfo.NextOperation = GetCommentsOverviewOperation
				monoChannel <- channelMetaInfo
				log.Println("--> (2/5): Sending out from GetVideoIDsOverview")
			} else {
				log.Println("x-x (2/5): NOT Working in GetVideoIDsOverview")
				monoChannel <- channelMetaInfo
				log.Println("--> (2/5): Sending out from GetVideoIDsOverview")
			}
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
			log.Println("<-- (3/5): Receiving into GetCommentsOverview")
			if channelMetaInfo.NextOperation == GetCommentsOverviewOperation {
				log.Println("<-> (3/5): Working in GetCommentsOverview")
				for i, video := range channelMetaInfo.Playlists["uploads"].PlaylistItems {
					go func(index int, inputVideo *Video) {
						call := service.CommentThreads.List("snippet").VideoId(inputVideo.VideoID).MaxResults(Opts.MaxResultsCommentPerVideo)
						response, responseError := call.Do()

						HandleError(responseError, fmt.Sprintf("GetCommentsOverview#%d Response error! (videoId: %s)", index, inputVideo.VideoID))

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
						log.Println("--> (3/5): Sending out from GetCommentsOverview")
					}(i, video)
				}
			} else {
				log.Println("x-x (3/5): NOT Working in GetCommentsOverview")
				monoChannel <- channelMetaInfo
				log.Println("--> (3/5): Sending out from GetCommentsOverview")
			}
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
			log.Println("<-- (4/5): Receiving into GetObviouslyRelatedChannelsOverview")
			if channelMetaInfo.NextOperation == GetObviouslyRelatedChannelsOverviewOperation {
				log.Println("<-> (4/5): Working in GetObviouslyRelatedChannelsOverview")
				for i, commentatorChannelID := range channelMetaInfo.CommentAuthorChannelIDs {
					go func(index int, inputCommentatorChannelID string) {
						Printfln("<-> (4/5): (#-----) Begin service.Channels.List for ChannelID: %v", inputCommentatorChannelID)
						getChannelCall := service.Channels.List("snippet,contentDetails").Id(inputCommentatorChannelID)
						getChannelResponse, getChannelResponseError := getChannelCall.Do()

						HandleError(getChannelResponseError, fmt.Sprintf("GetObviouslyRelatedChannelsOverview#%d Response error!", index))

						favoritesPlaylistID := getChannelResponse.Items[0].ContentDetails.RelatedPlaylists.Favorites
						if favoritesPlaylistID == "" {
							Printfln("<-> (4/5): (#X----) End service.Channels.List (error: %v)", getChannelResponseError)
							return
						}
						Printfln("<-> (4/5): (##----) End service.Channels.List (error: %v)", getChannelResponseError)
						Printfln("<-> (4/5): (###---) Begin service.PlaylistItems.List for PlaylistID: %v", favoritesPlaylistID)
						getPlaylistItemsCall := service.PlaylistItems.List("contentDetails").PlaylistId(favoritesPlaylistID).MaxResults(Opts.MaxResultsFavouritedVideos)
						getPlaylistItemsResponse, getPlaylistItemsResponseError := getPlaylistItemsCall.Do()
						Printfln("<-> (4/5): (####--) End service.PlaylistItems.List (error: %v)", getPlaylistItemsResponseError)

						HandleError(getPlaylistItemsResponseError, "GetObviouslyRelatedChannelsOverview#%d Response error!")

						var favoritedVideoIDs []string
						for _, item := range getPlaylistItemsResponse.Items {
							favoritedVideoIDs = append(favoritedVideoIDs, item.ContentDetails.VideoId)
						}

						Printfln("<-> (4/5): (#####-) Begin service.Videos.List for VideoIDs: %v", strings.Join(favoritedVideoIDs, ","))
						getRelatedChannelCall := service.Videos.List("snippet").Id(strings.Join(favoritedVideoIDs, ","))
						getRelatedChannelResponse, getRelatedChannelResponseError := getRelatedChannelCall.Do()
						Printfln("<-> (4/5): (######) End service.Videos.List (error: %v)", getRelatedChannelResponseError)

						HandleError(getRelatedChannelResponseError, fmt.Sprintf("GetObviouslyRelatedChannelsOverview#%d Response error!", index))

						var obviouslyRelatedChannelNames []string
						// ObviouslyRelatedChannelIDs
						for _, item := range getRelatedChannelResponse.Items {
							obviouslyRelatedChannelNames = append(obviouslyRelatedChannelNames, item.Snippet.ChannelTitle)
						}

						channelMetaInfo.ObviouslyRelatedChannelIDs = obviouslyRelatedChannelNames
						lastButNotLeastChannel <- channelMetaInfo
						log.Println("--> (4/5): Sending out from GetObviouslyRelatedChannelsOverview to lastButNotLeastChannel")

					}(i, commentatorChannelID)
				}
			} else {
				log.Println("x-x (4/5): NOT Working in GetObviouslyRelatedChannelsOverview")
				monoChannel <- channelMetaInfo
				log.Println("--> (4/5): Sending out from GetObviouslyRelatedChannelsOverview to monoChannel")
			}
		}
	}()
}

// CreateInitialChannelMetaInfo creates the initial request context
func CreateInitialChannelMetaInfo() ChannelMetaInfo {
	initialChannelMetaInfo := ChannelMetaInfo{}
	if Opts.PlaylistID == "" {
		initialChannelMetaInfo.ChannelID = Opts.ChannelID
		initialChannelMetaInfo.CustomURL = Opts.CustomURL
		initialChannelMetaInfo.NextOperation = GetChannelOverviewOperation
	} else {
		uploadedPlaylist := &Playlist{PlaylistID: Opts.PlaylistID}
		initialChannelMetaInfo.Playlists = make(map[string]*Playlist)
		initialChannelMetaInfo.Playlists["uploads"] = uploadedPlaylist
		initialChannelMetaInfo.NextOperation = GetVideoIDsOverviewOperation
	}
	return initialChannelMetaInfo
}
