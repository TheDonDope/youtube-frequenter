package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/user"
	"path/filepath"

	"google.golang.org/api/youtube/v3"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

// GetYouTubeService returns a service to interact with the YouTube API
func GetYouTubeService() (*youtube.Service, error) {
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

	return youtube.New(httpClient)
}

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
	Printfln("Go to the following link in your browser then type the "+
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
	Printfln("Saving credential file to: %s\n", file)
	openedFile, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer openedFile.Close()
	json.NewEncoder(openedFile).Encode(token)
}
