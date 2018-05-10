package main

import (
	"fmt"
	"os"
	"time"

	"github.com/jessevdk/go-flags"

	"github.com/TheDonDope/youtube-tinfoil-expose/go/service"
)

var opts struct {
	ChannelID string `short:"i" long:"channel-id" description:"The channel ID of a YouTube Channel"`

	CustomURL string `short:"u" long:"custom-url" description:"The custom URL of a YouTube Channel"`

	PlaylistName string `short:"n" long:"playlist-name" description:"The name of the playlist to search. Currently only supports the following playlists: uploads, favorites"`
}

func main() {
	start := time.Now()
	fmt.Println("Welcome to youtube-tinfoil-expose")

	_, argsError := flags.ParseArgs(&opts, os.Args)

	if argsError != nil {
		panic(argsError)
	}

	youtubeService, serviceError := service.GetYouTubeService()

	service.HandleError(serviceError, "Error creating YouTube client")

	query := service.Query{}
	query.CustomURL = opts.CustomURL
	query.PlaylistName = opts.PlaylistName

	results := service.Exfoliator(youtubeService, query)
	fmt.Println(results)
	elapsed := time.Since(start)
	fmt.Println(elapsed)
}
