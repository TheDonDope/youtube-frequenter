package exfoliating

import (
	"fmt"
	"log"
	"strings"
	"time"

	"gitlab.com/TheDonDope/gocha/v3/pkg/errors"
	"gitlab.com/TheDonDope/gocha/v3/pkg/logging"
	"gitlab.com/TheDonDope/youtube-frequenter/pkg/config"
	"gitlab.com/TheDonDope/youtube-frequenter/pkg/http/youtube"
)

// Service describes methods to interact with the Exfoliator API
type Service interface {
	// GetChannelOverview implements a Search which fills the most basic info about a channel.
	// To use the method simply pass a ChannelMetaInfo with either the ChannelID or CustomURL set.
	// On search completion the following attributes will be filled, if not already set:
	// - ChannelID
	// - CustomURL
	// - ChannelName
	// - Playlists (With their PlaylistID and PlaylistName)
	// - SubscriberCount
	// - ViewCount
	GetChannelOverview(infos chan youtube.ChannelMetaInfo)

	// GetVideoIDsOverview gets all videos.
	GetVideoIDsOverview(infos chan youtube.ChannelMetaInfo)

	// GetCommentsOverview returns the overview over the comments
	GetCommentsOverview(infos chan youtube.ChannelMetaInfo)

	// GetObviouslyRelatedChannelsOverview gets the related channels for a YouTube channel
	GetObviouslyRelatedChannelsOverview(infos chan youtube.ChannelMetaInfo, lasts chan youtube.ChannelMetaInfo)

	// CreateInitialChannelMetaInfo creates the initial request context
	CreateInitialChannelMetaInfo() youtube.ChannelMetaInfo

	// Exfoliate returns the result
	Exfoliate(info youtube.ChannelMetaInfo) youtube.ChannelMetaInfo

	// AnalyseChannelMetaInfo prints additional information for a given channelMetaInfo.
	AnalyseChannelMetaInfo(info youtube.ChannelMetaInfo)
}

// service is the struct implementing the Service interface
type service struct {
	yt youtube.Service
}

// NewService creates an Exfoliator service with the necessary dependencies
func NewService(yt youtube.Service) Service {
	return &service{yt}
}

// GetChannelOverview implements a Search which fills the most basic info about a channel.
// To use the method simply pass a ChannelMetaInfo with either the ChannelID or CustomURL set.
// On search completion the following attributes will be filled, if not already set:
// - ChannelID
// - CustomURL
// - ChannelName
// - Playlists (With their PlaylistID and PlaylistName)
// - SubscriberCount
// - ViewCount
func (s *service) GetChannelOverview(infos chan youtube.ChannelMetaInfo) {
	go func() {
		log.Println("<<<<<Begin GetChannelOverview Go Routine")
		defer log.Println("End GetChannelOverview Go Routine>>>>>")
		for {
			info := <-infos
			log.Println("<-- (1/5): Receiving into GetChannelOverview")
			if info.NextOperation == youtube.GetChannelOverviewOperation {
				log.Println("<-> (1/5): Working in GetChannelOverview")
				response, responseError := s.yt.ChannelsList(info.ChannelID, info.CustomURL)
				errors.Print(responseError, "GetChannelOverview Response error!")

				firstItem := response.Items[0]

				// Fill ChannelID
				if info.ChannelID == "" {
					info.ChannelID = firstItem.Id
				}

				// Fill CustomURL
				if info.CustomURL == "" {
					info.CustomURL = firstItem.Snippet.CustomUrl
				}

				// Fill ChannelName
				if info.ChannelName == "" {
					info.ChannelName = firstItem.Snippet.Title
				}

				// Fill Playlists
				if info.Playlists == nil {
					info.Playlists = make(map[string]*youtube.Playlist)
					var videos []*youtube.Video
					info.Playlists["uploads"] = &youtube.Playlist{PlaylistID: firstItem.ContentDetails.RelatedPlaylists.Uploads, PlaylistItems: videos}
					info.Playlists["favorites"] = &youtube.Playlist{PlaylistID: firstItem.ContentDetails.RelatedPlaylists.Favorites, PlaylistItems: videos}
				}

				// Fill SubscriberCount
				if info.SubscriberCount == 0 {
					info.SubscriberCount = firstItem.Statistics.SubscriberCount
				}

				// Fill ViewCount
				if info.ViewCount == 0 {
					info.ViewCount = firstItem.Statistics.ViewCount
				}

				info.NextOperation = youtube.GetVideoIDsOverviewOperation
				infos <- info
				log.Println("--> (1/5): Sending out from GetChannelOverview")
			} else {
				log.Println("x-x (1/5): NOT Working in GetChannelOverview")
				infos <- info
				log.Println("--> (1/5): Sendubg out from GetChannelOverview")
			}
		}
	}()

}

// GetVideoIDsOverview gets all videos.
func (s *service) GetVideoIDsOverview(infos chan youtube.ChannelMetaInfo) {
	go func() {
		log.Println("<<<<<Begin GetVideoIDsOverview Go Routine")
		defer log.Println("End GetVideoIDsOverview Go Routine>>>>>")
		for {
			info := <-infos
			log.Println("<-- (2/5): Receiving into GetVideoIDsOverview")
			if info.NextOperation == youtube.GetVideoIDsOverviewOperation {
				log.Println("<-> (2/5): Working in GetVideoIDsOverview")
				response, responseError := s.yt.PlaylistItemsList(info.Playlists["uploads"].PlaylistID, "uploads")

				errors.Print(responseError, "GetVideoIDsOverview Response error!")

				var videos []*youtube.Video
				for _, item := range response.Items {
					video := &youtube.Video{VideoID: item.Snippet.ResourceId.VideoId}
					videos = append(videos, video)
				}

				uploadPlaylist := info.Playlists["uploads"]
				uploadPlaylist.PlaylistItems = videos
				info.NextOperation = youtube.GetCommentsOverviewOperation
				infos <- info
				log.Println("--> (2/5): Sending out from GetVideoIDsOverview")
			} else {
				log.Println("x-x (2/5): NOT Working in GetVideoIDsOverview")
				infos <- info
				log.Println("--> (2/5): Sending out from GetVideoIDsOverview")
			}
		}
	}()
}

// GetCommentsOverview foo
func (s *service) GetCommentsOverview(infos chan youtube.ChannelMetaInfo) {
	go func() {
		log.Println("<<<<<Begin GetCommentsOverview Go Routine")
		defer log.Println("End GetCommentsOverview Go Routine>>>>>")
		for {
			info := <-infos
			log.Println("<-- (3/5): Receiving into GetCommentsOverview")
			if info.NextOperation == youtube.GetCommentsOverviewOperation {
				log.Println("<-> (3/5): Working in GetCommentsOverview")
				for i, video := range info.Playlists["uploads"].PlaylistItems {
					go func(index int, inputVideo *youtube.Video) {
						response, responseError := s.yt.CommentThreadsList(inputVideo.VideoID)

						errors.Print(responseError, fmt.Sprintf("GetCommentsOverview#%d Response error! (videoId: %s)", index, inputVideo.VideoID))

						var comments []*youtube.Comment
						for _, item := range response.Items {
							comment := new(youtube.Comment)
							comment.CommentID = item.Snippet.TopLevelComment.Id
							comment.AuthorChannelID = item.Snippet.TopLevelComment.Snippet.AuthorChannelId.(map[string]interface{})["value"].(string)
							info.CommentAuthorChannelIDs = append(info.CommentAuthorChannelIDs, comment.AuthorChannelID)

							comments = append(comments, comment)
						}

						inputVideo.Comments = comments
						info.NextOperation = youtube.GetObviouslyRelatedChannelsOverviewOperation
						infos <- info
						log.Println("--> (3/5): Sending out from GetCommentsOverview")
					}(i, video)
				}
			} else {
				log.Println("x-x (3/5): NOT Working in GetCommentsOverview")
				infos <- info
				log.Println("--> (3/5): Sending out from GetCommentsOverview")
			}
		}
	}()
}

// GetObviouslyRelatedChannelsOverview gets the related channels for a YouTube channel
func (s *service) GetObviouslyRelatedChannelsOverview(infos chan youtube.ChannelMetaInfo, lasts chan youtube.ChannelMetaInfo) {
	go func() {
		log.Println("<<<<<Begin GetObviouslyRelatedChannelsOverview Go Routine")
		defer log.Println("End GetObviouslyRelatedChannelsOverview Go Routine>>>>>")
		for {
			info := <-infos
			log.Println("<-- (4/5): Receiving into GetObviouslyRelatedChannelsOverview")
			if info.NextOperation == youtube.GetObviouslyRelatedChannelsOverviewOperation {
				log.Println("<-> (4/5): Working in GetObviouslyRelatedChannelsOverview")
				for i, commentatorChannelID := range info.CommentAuthorChannelIDs {
					go func(index int, inputCommentatorChannelID string) {
						logging.Printfln("<-> (4/5): (#-----) Begin service.Channels.List for ChannelID: %v", inputCommentatorChannelID)
						getChannelResponse, getChannelResponseError := s.yt.ChannelsList(inputCommentatorChannelID, "")

						errors.Print(getChannelResponseError, fmt.Sprintf("GetObviouslyRelatedChannelsOverview#%d Response error!", index))

						favoritesPlaylistID := getChannelResponse.Items[0].ContentDetails.RelatedPlaylists.Favorites
						if favoritesPlaylistID == "" {
							logging.Printfln("<-> (4/5): (#X----) End service.Channels.List (error: %v)", getChannelResponseError)
							return
						}
						logging.Printfln("<-> (4/5): (##----) End service.Channels.List (error: %v)", getChannelResponseError)
						logging.Printfln("<-> (4/5): (###---) Begin service.PlaylistItems.List for PlaylistID: %v", favoritesPlaylistID)

						getPlaylistItemsResponse, getPlaylistItemsResponseError := s.yt.PlaylistItemsList(favoritesPlaylistID, "favorites")
						logging.Printfln("<-> (4/5): (####--) End service.PlaylistItems.List (error: %v)", getPlaylistItemsResponseError)

						errors.Print(getPlaylistItemsResponseError, "GetObviouslyRelatedChannelsOverview#%d Response error!")

						var favoritedVideoIDs []string
						for _, item := range getPlaylistItemsResponse.Items {
							favoritedVideoIDs = append(favoritedVideoIDs, item.ContentDetails.VideoId)
						}

						logging.Printfln("<-> (4/5): (#####-) Begin service.Videos.List for VideoIDs: %v", strings.Join(favoritedVideoIDs, ","))
						getRelatedChannelResponse, getRelatedChannelResponseError :=
							s.yt.VideosList(strings.Join(favoritedVideoIDs, ","))

						logging.Printfln("<-> (4/5): (######) End service.Videos.List (error: %v)", getRelatedChannelResponseError)

						errors.Print(getRelatedChannelResponseError, fmt.Sprintf("GetObviouslyRelatedChannelsOverview#%d Response error!", index))

						var obviouslyRelatedChannelNames []string
						// ObviouslyRelatedChannelIDs
						for _, item := range getRelatedChannelResponse.Items {
							obviouslyRelatedChannelNames = append(obviouslyRelatedChannelNames, item.Snippet.ChannelTitle)
						}

						info.ObviouslyRelatedChannelIDs = obviouslyRelatedChannelNames
						lasts <- info
						log.Println("--> (4/5): Sending out from GetObviouslyRelatedChannelsOverview to lastButNotLeastChannel")

					}(i, commentatorChannelID)
				}
			} else {
				log.Println("x-x (4/5): NOT Working in GetObviouslyRelatedChannelsOverview")
				infos <- info
				log.Println("--> (4/5): Sending out from GetObviouslyRelatedChannelsOverview to monoChannel")
			}
		}
	}()
}

// CreateInitialChannelMetaInfo creates the initial request context
func (s *service) CreateInitialChannelMetaInfo() youtube.ChannelMetaInfo {
	initialChannelMetaInfo := youtube.ChannelMetaInfo{}
	if config.Opts.PlaylistID == "" {
		initialChannelMetaInfo.ChannelID = config.Opts.ChannelID
		initialChannelMetaInfo.CustomURL = config.Opts.CustomURL
		initialChannelMetaInfo.NextOperation = youtube.GetChannelOverviewOperation
	} else {
		uploadedPlaylist := &youtube.Playlist{PlaylistID: config.Opts.PlaylistID}
		initialChannelMetaInfo.Playlists = make(map[string]*youtube.Playlist)
		initialChannelMetaInfo.Playlists["uploads"] = uploadedPlaylist
		initialChannelMetaInfo.NextOperation = youtube.GetVideoIDsOverviewOperation
	}
	return initialChannelMetaInfo
}

// Exfoliate returns the result
func (s *service) Exfoliate(info youtube.ChannelMetaInfo) youtube.ChannelMetaInfo {
	monoChannel := make(chan youtube.ChannelMetaInfo)
	lastButNotLeastChannel := make(chan youtube.ChannelMetaInfo)
	accumulatedMetaInfo := youtube.ChannelMetaInfo{}
	accumulatedMetaInfo.CustomURL = info.CustomURL
	accumulatedMetaInfo.ChannelID = info.ChannelID
	accumulatedMetaInfo.Playlists = info.Playlists
	go func() {
		monoChannel <- info
	}()
	go s.GetChannelOverview(monoChannel)
	go s.GetVideoIDsOverview(monoChannel)
	go s.GetCommentsOverview(monoChannel)
	go s.GetObviouslyRelatedChannelsOverview(monoChannel, lastButNotLeastChannel)

	globalTimeout, globalTimeoutError := time.ParseDuration(config.Opts.GlobalTimeout)
	if globalTimeoutError != nil {
		log.Println(globalTimeoutError)
	}
	timeout := time.After(globalTimeout)
	for {
		log.Println("<<<<<Begin Exfoliator Main Loop")
		select {
		case channelMetaInfo := <-lastButNotLeastChannel:
			log.Println("<-- (5/5): Exfoliator")
			log.Println("<-> (5/5): Working in Exfoliator")
			// evtl die anderen properties adden
			accumulatedMetaInfo.ObviouslyRelatedChannelIDs = append(accumulatedMetaInfo.ObviouslyRelatedChannelIDs, channelMetaInfo.ObviouslyRelatedChannelIDs...)
			log.Println("--> (5/5): Exfoliator")
		case <-timeout:
			logging.Printfln("Request timed out (%v)", config.Opts.GlobalTimeout)
			return accumulatedMetaInfo
		}
	}
}

// AnalyseChannelMetaInfo prints additional information for a given channelMetaInfo.
func (s *service) AnalyseChannelMetaInfo(info youtube.ChannelMetaInfo) {
	// relatedChannelIDToNumberOfOccurrences := collections.CountOccurrences(info.ObviouslyRelatedChannelIDs)

	// if len(relatedChannelIDToNumberOfOccurrences) == 0 {
	// 	log.Println("Package to analyse has no ObviouslyRelatedChannelIDs to count.")
	// } else {
	// 	sortedRelatedChannelIDsList := commonTypes.MapEntryList{}.RankByWordCount(relatedChannelIDToNumberOfOccurrences)

	// 	resultJSONBytes, resultJSONBytesError := json.Marshal(sortedRelatedChannelIDsList)
	// 	errors.Print(resultJSONBytesError, "Error marshaling results")
	// 	storage.ToJSON(config.GetOutPath()+"/"+config.GetOutName()+"-results.json", resultJSONBytes)
	// 	sortedRelatedChannelIDsList.PrintResults()

	// 	dumpJSONBytes, dumpJSONBytesError := json.Marshal(info)
	// 	errors.Print(dumpJSONBytesError, "Error marshaling dump")
	// 	storage.ToJSON(config.GetOutPath()+"/"+config.GetOutName()+"-dump.json", dumpJSONBytes)
	// 	logging.Printfln("#Results: %d", len(relatedChannelIDToNumberOfOccurrences))
	// }
}
