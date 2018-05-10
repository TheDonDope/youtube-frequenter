package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"os/user"
	"path/filepath"
	"time"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/youtube/v3"
)

var (
	// Web1 Replicae
	PlaylistID1 = MetaSearch("PlaylistID1")
	// Web2 Replicae
	PlaylistID2 = MetaSearch("PlaylistID2")
)

const missingClientSecretsMessage = `
Please configure OAuth 2.0
`

// Result datatype
type Result struct {
	ResultValue string
	ResultError error
}

type Query struct {
	ChannelID    string
	CustomURL    string
	PlaylistName string
}

// Search function declaration
type Search func(serviceservice *youtube.Service, query Query) Result

// getHTTPClient uses a Context and Config to retrieve a Token
// then generate a Client. It returns the generated Client.
func getHTTPClient(context context.Context, oauth2Configuration *oauth2.Config) *http.Client {

	cacheFile, cacheFileError := createTokenCacheFile()
	if cacheFileError != nil {
		log.Fatalf("Unable to get path to cached credential file. %v", cacheFileError)
	}
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
	authURL := oauth2Configuration.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var code string
	if _, err := fmt.Scan(&code); err != nil {
		log.Fatalf("Unable to read authorization code %v", err)
	}

	token, tokenError := oauth2Configuration.Exchange(oauth2.NoContext, code)
	if tokenError != nil {
		log.Fatalf("Unable to retrieve token from web %v", tokenError)
	}
	return token
}

func saveToken(file string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", file)
	openedFile, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer openedFile.Close()
	json.NewEncoder(openedFile).Encode(token)
}

func handleError(errorToHandle error, errorMessage string) {
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
	handleError(responseError, "Response error!")
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

func Exfoliater(service *youtube.Service, query Query) (resultsFromChannel []Result) {
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

func main() {
	start := time.Now()
	fmt.Println("Welcome to youtube-tinfoil-expose")

	backgroundContext := context.Background()

	readBytes, readError := ioutil.ReadFile("client_secret.json")
	if readError != nil {
		log.Fatalf("Unable to read client secret file: %v", readError)
	}

	configFromJSON, configError := google.ConfigFromJSON(readBytes, youtube.YoutubeReadonlyScope)
	if configError != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", configError)
	}
	httpClient := getHTTPClient(backgroundContext, configFromJSON)
	youtubeService, serviceError := youtube.New(httpClient)

	handleError(serviceError, "Error creating YouTube client")
	/*
		getPlaylistIDByChannelIDPlaylistName(youtubeService, "UCu9ljRg6YrwSw64qggVdczQ", "uploads")
		getPlaylistIDByChannelIDPlaylistName(youtubeService, "UCu9ljRg6YrwSw64qggVdczQ", "favorites")
		getPlaylistIDByCustomURLPlaylistName(youtubeService, "wwwKenFMde", "uploads")
		getPlaylistIDByCustomURLPlaylistName(youtubeService, "wwwKenFMde", "favorites")
	*/
	query := Query{}
	query.CustomURL = "wwwKenFMde"
	query.PlaylistName = "uploads"

	results := Exfoliater(youtubeService, query)
	fmt.Println(results)
	elapsed := time.Since(start)
	fmt.Println(elapsed)
}
