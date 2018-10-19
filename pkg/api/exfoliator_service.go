package api

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"gitlab.com/TheDonDope/youtube-frequenter/pkg/types"
	"gitlab.com/TheDonDope/youtube-frequenter/pkg/util/collections"
	"gitlab.com/TheDonDope/youtube-frequenter/pkg/util/configs"
	"gitlab.com/TheDonDope/youtube-frequenter/pkg/util/errors"
	"gitlab.com/TheDonDope/youtube-frequenter/pkg/util/files"
	"gitlab.com/TheDonDope/youtube-frequenter/pkg/util/logs"
	"google.golang.org/api/youtube/v3"
)

// ExfoliatorService implements the Exfoliator interface
type ExfoliatorService struct{}

// GetChannelOverview implements a Search which fills the most basic info about a channel.
// To use the method simply pass a ChannelMetaInfo with either the ChannelID or CustomURL set.
// On search completion the following attributes will be filled, if not already set:
// - ChannelID
// - CustomURL
// - ChannelName
// - Playlists (With their PlaylistID and PlaylistName)
// - SubscriberCount
// - ViewCount
func (impl ExfoliatorService) GetChannelOverview(service *youtube.Service, serviceImpl YouTuber, monoChannel chan types.ChannelMetaInfo) {
	go func() {
		log.Println("<<<<<Begin GetChannelOverview Go Routine")
		defer log.Println("End GetChannelOverview Go Routine>>>>>")
		for {
			channelMetaInfo := <-monoChannel
			log.Println("<-- (1/5): Receiving into GetChannelOverview")
			if channelMetaInfo.NextOperation == types.GetChannelOverviewOperation {
				log.Println("<-> (1/5): Working in GetChannelOverview")
				response, responseError := serviceImpl.ChannelsList(service, channelMetaInfo.ChannelID, channelMetaInfo.CustomURL)

				errors.HandleError(responseError, "GetChannelOverview Response error!")

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
					channelMetaInfo.Playlists = make(map[string]*types.Playlist)
					var videos []*types.Video
					channelMetaInfo.Playlists["uploads"] = &types.Playlist{PlaylistID: firstItem.ContentDetails.RelatedPlaylists.Uploads, PlaylistItems: videos}
					channelMetaInfo.Playlists["favorites"] = &types.Playlist{PlaylistID: firstItem.ContentDetails.RelatedPlaylists.Favorites, PlaylistItems: videos}
				}

				// Fill SubscriberCount
				if channelMetaInfo.SubscriberCount == 0 {
					channelMetaInfo.SubscriberCount = firstItem.Statistics.SubscriberCount
				}

				// Fill ViewCount
				if channelMetaInfo.ViewCount == 0 {
					channelMetaInfo.ViewCount = firstItem.Statistics.ViewCount
				}

				channelMetaInfo.NextOperation = types.GetVideoIDsOverviewOperation
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
func (impl ExfoliatorService) GetVideoIDsOverview(service *youtube.Service, serviceImpl YouTuber, monoChannel chan types.ChannelMetaInfo) {
	go func() {
		log.Println("<<<<<Begin GetVideoIDsOverview Go Routine")
		defer log.Println("End GetVideoIDsOverview Go Routine>>>>>")
		for {
			channelMetaInfo := <-monoChannel
			log.Println("<-- (2/5): Receiving into GetVideoIDsOverview")
			if channelMetaInfo.NextOperation == types.GetVideoIDsOverviewOperation {
				log.Println("<-> (2/5): Working in GetVideoIDsOverview")
				response, responseError := serviceImpl.PlaylistItemsList(service, channelMetaInfo.Playlists["uploads"].PlaylistID, "uploads")

				errors.HandleError(responseError, "GetVideoIDsOverview Response error!")

				var videos []*types.Video
				for _, item := range response.Items {
					video := &types.Video{VideoID: item.Snippet.ResourceId.VideoId}
					videos = append(videos, video)
				}

				uploadPlaylist := channelMetaInfo.Playlists["uploads"]
				uploadPlaylist.PlaylistItems = videos
				channelMetaInfo.NextOperation = types.GetCommentsOverviewOperation
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
func (impl ExfoliatorService) GetCommentsOverview(service *youtube.Service, serviceImpl YouTuber, monoChannel chan types.ChannelMetaInfo) {
	go func() {
		log.Println("<<<<<Begin GetCommentsOverview Go Routine")
		defer log.Println("End GetCommentsOverview Go Routine>>>>>")
		for {
			channelMetaInfo := <-monoChannel
			log.Println("<-- (3/5): Receiving into GetCommentsOverview")
			if channelMetaInfo.NextOperation == types.GetCommentsOverviewOperation {
				log.Println("<-> (3/5): Working in GetCommentsOverview")
				for i, video := range channelMetaInfo.Playlists["uploads"].PlaylistItems {
					go func(index int, inputVideo *types.Video) {
						response, responseError := serviceImpl.CommentThreadsList(service, inputVideo.VideoID)

						errors.HandleError(responseError, fmt.Sprintf("GetCommentsOverview#%d Response error! (videoId: %s)", index, inputVideo.VideoID))

						var comments []*types.Comment
						for _, item := range response.Items {
							comment := new(types.Comment)
							comment.CommentID = item.Snippet.TopLevelComment.Id
							comment.AuthorChannelID = item.Snippet.TopLevelComment.Snippet.AuthorChannelId.(map[string]interface{})["value"].(string)
							channelMetaInfo.CommentAuthorChannelIDs = append(channelMetaInfo.CommentAuthorChannelIDs, comment.AuthorChannelID)

							comments = append(comments, comment)
						}

						inputVideo.Comments = comments
						channelMetaInfo.NextOperation = types.GetObviouslyRelatedChannelsOverviewOperation
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
func (impl ExfoliatorService) GetObviouslyRelatedChannelsOverview(service *youtube.Service, serviceImpl YouTuber, monoChannel chan types.ChannelMetaInfo, lastButNotLeastChannel chan types.ChannelMetaInfo) {
	go func() {
		log.Println("<<<<<Begin GetObviouslyRelatedChannelsOverview Go Routine")
		defer log.Println("End GetObviouslyRelatedChannelsOverview Go Routine>>>>>")
		for {
			channelMetaInfo := <-monoChannel
			log.Println("<-- (4/5): Receiving into GetObviouslyRelatedChannelsOverview")
			if channelMetaInfo.NextOperation == types.GetObviouslyRelatedChannelsOverviewOperation {
				log.Println("<-> (4/5): Working in GetObviouslyRelatedChannelsOverview")
				for i, commentatorChannelID := range channelMetaInfo.CommentAuthorChannelIDs {
					go func(index int, inputCommentatorChannelID string) {
						logs.Printfln("<-> (4/5): (#-----) Begin service.Channels.List for ChannelID: %v", inputCommentatorChannelID)
						getChannelResponse, getChannelResponseError := serviceImpl.ChannelsList(service, inputCommentatorChannelID, "")

						errors.HandleError(getChannelResponseError, fmt.Sprintf("GetObviouslyRelatedChannelsOverview#%d Response error!", index))

						favoritesPlaylistID := getChannelResponse.Items[0].ContentDetails.RelatedPlaylists.Favorites
						if favoritesPlaylistID == "" {
							logs.Printfln("<-> (4/5): (#X----) End service.Channels.List (error: %v)", getChannelResponseError)
							return
						}
						logs.Printfln("<-> (4/5): (##----) End service.Channels.List (error: %v)", getChannelResponseError)
						logs.Printfln("<-> (4/5): (###---) Begin service.PlaylistItems.List for PlaylistID: %v", favoritesPlaylistID)

						getPlaylistItemsResponse, getPlaylistItemsResponseError := serviceImpl.PlaylistItemsList(service, favoritesPlaylistID, "favorites")
						logs.Printfln("<-> (4/5): (####--) End service.PlaylistItems.List (error: %v)", getPlaylistItemsResponseError)

						errors.HandleError(getPlaylistItemsResponseError, "GetObviouslyRelatedChannelsOverview#%d Response error!")

						var favoritedVideoIDs []string
						for _, item := range getPlaylistItemsResponse.Items {
							favoritedVideoIDs = append(favoritedVideoIDs, item.ContentDetails.VideoId)
						}

						logs.Printfln("<-> (4/5): (#####-) Begin service.Videos.List for VideoIDs: %v", strings.Join(favoritedVideoIDs, ","))
						getRelatedChannelResponse, getRelatedChannelResponseError :=
							serviceImpl.VideosList(service, strings.Join(favoritedVideoIDs, ","))

						logs.Printfln("<-> (4/5): (######) End service.Videos.List (error: %v)", getRelatedChannelResponseError)

						errors.HandleError(getRelatedChannelResponseError, fmt.Sprintf("GetObviouslyRelatedChannelsOverview#%d Response error!", index))

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
func (impl ExfoliatorService) CreateInitialChannelMetaInfo() types.ChannelMetaInfo {
	initialChannelMetaInfo := types.ChannelMetaInfo{}
	if configs.Opts.PlaylistID == "" {
		initialChannelMetaInfo.ChannelID = configs.Opts.ChannelID
		initialChannelMetaInfo.CustomURL = configs.Opts.CustomURL
		initialChannelMetaInfo.NextOperation = types.GetChannelOverviewOperation
	} else {
		uploadedPlaylist := &types.Playlist{PlaylistID: configs.Opts.PlaylistID}
		initialChannelMetaInfo.Playlists = make(map[string]*types.Playlist)
		initialChannelMetaInfo.Playlists["uploads"] = uploadedPlaylist
		initialChannelMetaInfo.NextOperation = types.GetVideoIDsOverviewOperation
	}
	return initialChannelMetaInfo
}

// Exfoliate returns the result
func (impl ExfoliatorService) Exfoliate(service *youtube.Service, serviceImpl YouTuber, channelMetaInfo types.ChannelMetaInfo) types.ChannelMetaInfo {
	monoChannel := make(chan types.ChannelMetaInfo)
	lastButNotLeastChannel := make(chan types.ChannelMetaInfo)
	accumulatedMetaInfo := types.ChannelMetaInfo{}
	accumulatedMetaInfo.CustomURL = channelMetaInfo.CustomURL
	accumulatedMetaInfo.ChannelID = channelMetaInfo.ChannelID
	accumulatedMetaInfo.Playlists = channelMetaInfo.Playlists
	go func() {
		monoChannel <- channelMetaInfo
	}()
	go impl.GetChannelOverview(service, serviceImpl, monoChannel)
	go impl.GetVideoIDsOverview(service, serviceImpl, monoChannel)
	go impl.GetCommentsOverview(service, serviceImpl, monoChannel)
	go impl.GetObviouslyRelatedChannelsOverview(service, serviceImpl, monoChannel, lastButNotLeastChannel)

	globalTimeout, globalTimeoutError := time.ParseDuration(configs.Opts.GlobalTimeout)
	if globalTimeoutError != nil {
		log.Println(globalTimeoutError)
	}
	timeout := time.After(globalTimeout)
	for {
		log.Println("<<<<<Begin Exfoliator Main Loop")
		select {
		case channelMetaInfo = <-lastButNotLeastChannel:
			log.Println("<-- (5/5): Exfoliator")
			log.Println("<-> (5/5): Working in Exfoliator")
			// evtl die anderen properties adden
			accumulatedMetaInfo.ObviouslyRelatedChannelIDs = append(accumulatedMetaInfo.ObviouslyRelatedChannelIDs, channelMetaInfo.ObviouslyRelatedChannelIDs...)
			log.Println("--> (5/5): Exfoliator")
		case <-timeout:
			logs.Printfln("Request timed out (%v)", configs.Opts.GlobalTimeout)
			return accumulatedMetaInfo
		}
	}
}

// AnalyseChannelMetaInfo prints additional information for a given channelMetaInfo.
func (impl ExfoliatorService) AnalyseChannelMetaInfo(channelMetaInfo *types.ChannelMetaInfo) {
	relatedChannelIDToNumberOfOccurrences := collections.CountOccurrences(channelMetaInfo.ObviouslyRelatedChannelIDs)

	if len(relatedChannelIDToNumberOfOccurrences) == 0 {
		log.Println("Package to analyse has no ObviouslyRelatedChannelIDs to count.")
	} else {
		sortedRelatedChannelIDsList := types.MapEntryList{}.RankByWordCount(relatedChannelIDToNumberOfOccurrences)

		resultJSONBytes, resultJSONBytesError := json.Marshal(sortedRelatedChannelIDsList)
		errors.HandleError(resultJSONBytesError, "Error marshaling results")
		files.WriteToJSON(configs.GetOutputDirectory()+"/"+configs.GetCustomName()+"-results.json", resultJSONBytes)
		sortedRelatedChannelIDsList.PrintResults()

		dumpJSONBytes, dumpJSONBytesError := json.Marshal(channelMetaInfo)
		errors.HandleError(dumpJSONBytesError, "Error marshaling dump")
		files.WriteToJSON(configs.GetOutputDirectory()+"/"+configs.GetCustomName()+"-dump.json", dumpJSONBytes)
		logs.Printfln("#Results: %d", len(relatedChannelIDToNumberOfOccurrences))
	}
}
