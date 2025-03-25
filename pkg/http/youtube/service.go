package youtube

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/user"
	"path/filepath"

	"github.com/TheDonDope/youtube-frequenter/pkg/config"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	youtubeV3 "google.golang.org/api/youtube/v3"
)

// Service describes methods to interact with the YouTube API
type Service interface {
	// ChannelsList returns the result of the Channels.list API call
	ChannelsList(channelID string, customURL string) (*youtubeV3.ChannelListResponse, error)

	// PlaylistItemsList returns the result of the PlaylistItems.list API call
	PlaylistItemsList(playlistID string, playlistName string) (*youtubeV3.PlaylistItemListResponse, error)

	// CommentThreadsList returns the result of the CommentThreads.list API call
	CommentThreadsList(videoID string) (*youtubeV3.CommentThreadListResponse, error)

	// VideosList returns the result of the Videos.list API call
	VideosList(videoIDs string) (*youtubeV3.VideoListResponse, error)
}

// service is the struct implementing the Service interface
type service struct {
	ytV3 *youtubeV3.Service
}

// NewService creates an youtube service with the necessary dependencies
func NewService(ytV3 *youtubeV3.Service) Service {
	return &service{ytV3}
}

// NewYouTubeV3 returns a service to interact with the YouTube API
func NewYouTubeV3() (*youtubeV3.Service, error) {
	backgroundContext := context.Background()

	readBytes, _ := os.ReadFile("client_secret.json")
	// TODO: rewrite
	// errors.Print(readError, fmt.Sprintf("Unable to read client secret file: %v", readError))

	configFromJSON, _ := google.ConfigFromJSON(readBytes, youtubeV3.YoutubeForceSslScope)
	// TOOD: rewrite
	// errors.Print(configError, fmt.Sprintf("Unable to parse client secret file to config: %v", configError))
	httpClient := getHTTPClient(backgroundContext, configFromJSON)

	return youtubeV3.NewService(context.Background(), option.WithHTTPClient(httpClient))
}

// ChannelsList returns the result of the Channels.list API call
func (s *service) ChannelsList(channelID string, customURL string) (*youtubeV3.ChannelListResponse, error) {
	call := s.ytV3.Channels.List([]string{"contentDetails", "snippet", "statistics"})
	if channelID != "" {
		call = call.Id(channelID)
	} else if customURL != "" {
		call = call.ForUsername(customURL)
	}
	return call.Do()
}

// PlaylistItemsList returns the result of the PlaylistItems.list API call
func (s *service) PlaylistItemsList(playlistID string, playlistName string) (*youtubeV3.PlaylistItemListResponse, error) {
	call := s.ytV3.PlaylistItems.List([]string{"contentDetails", "snippet"}).PlaylistId(playlistID)
	if playlistName == "uploads" {
		call = call.MaxResults(config.Opts.MaxResultsUploadedVideos)
	} else if playlistName == "favorites" {
		call = call.MaxResults(config.Opts.MaxResultsFavouritedVideos)
	}
	return call.Do()
}

// CommentThreadsList returns the result of the CommentThreads.list API call
func (s *service) CommentThreadsList(videoID string) (*youtubeV3.CommentThreadListResponse, error) {
	return s.ytV3.CommentThreads.List([]string{"snippet"}).VideoId(videoID).MaxResults(config.Opts.MaxResultsCommentPerVideo).Do()
}

// VideosList returns the result of the Videos.list API call
func (s *service) VideosList(videoIDs string) (*youtubeV3.VideoListResponse, error) {
	return s.ytV3.Videos.List([]string{"snippet"}).Id(videoIDs).Do()
}

// getHTTPClient uses a Context and Config to retrieve a Token
// then generate a Client. It returns the generated Client.
func getHTTPClient(context context.Context, oauth2Configuration *oauth2.Config) *http.Client {

	cacheFile, _ := createTokenCacheFile()
	// TODO: rewrite
	// errors.Print(cacheFileError, fmt.Sprintf("Unable to get path to cached credential file. %v", cacheFileError))
	token, tokenError := getTokenFromFile(cacheFile)
	if tokenError != nil {
		token = getTokenFromWeb(oauth2Configuration)
		saveToken(cacheFile, token)
	}
	return oauth2Configuration.Client(context, token)
}

// createTokenCacheFile generates a credential file path/filename.
// It returns the generated credential path/filename.
func createTokenCacheFile() (string, error) {
	currentUser, err := user.Current()
	if err != nil {
		return "", err
	}
	tokenCacheDirectory := filepath.Join(currentUser.HomeDir, ".credentials")
	os.MkdirAll(tokenCacheDirectory, 0700)
	return filepath.Join(tokenCacheDirectory, url.QueryEscape("youtube-tinfoil-expose.json")), err
}

// getTokenFromFile retrieves a Token from a given file path.
// It returns the retrieved Token and any read error encountered.
func getTokenFromFile(filePath string) (*oauth2.Token, error) {
	openedFile, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	resultToken := &oauth2.Token{}
	err = json.NewDecoder(openedFile).Decode(resultToken)
	defer openedFile.Close()
	return resultToken, err
}

// getTokenFromWeb uses Config to request a Token.
// It returns the retrieved Token.
func getTokenFromWeb(oauth2Configuration *oauth2.Config) *oauth2.Token {
	// TODO: rewrite
	// authURL := oauth2Configuration.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	// logging.Printfln("Go to the following link in your browser then type the "+
	//	"authorization code: \n%v\n", authURL)

	var code string
	if _, err := fmt.Scan(&code); err != nil {
		log.Fatalf("Unable to read authorization code %v", err)
	}

	token, _ := oauth2Configuration.Exchange(context.Background(), code)
	// TODO: rewrite
	// errors.Print(tokenError, fmt.Sprintf("Unable to retrieve token from web %v", tokenError))

	return token
}

func saveToken(file string, token *oauth2.Token) {
	// TODO: rewrite
	// logging.Printfln("Saving credential file to: %s\n", file)
	openedFile, _ := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	// TODO: rewrite
	// errors.Print(err, fmt.Sprintf("Unable to cache oauth token: %v", err))
	defer openedFile.Close()
	json.NewEncoder(openedFile).Encode(token)
}
