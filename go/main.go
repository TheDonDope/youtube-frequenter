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

	CustomURL string `short:"u" long:"custom-url" description:"The custom URL of a YouTube Channel" required:"true"`

	PlaylistName string `short:"n" long:"playlist-name" description:"The name of the playlist to search. Currently only supports the following playlists: uploads, favorites"`
}

func main() {
	start := time.Now()
	fmt.Println(fmt.Sprintf("Starting youtube-tinfoil-expose @ %v", start.Format(time.RFC3339)))

	_, argsError := flags.ParseArgs(&opts, os.Args)
	if argsError != nil {
		panic(argsError)
	}

	youtubeService, serviceError := service.GetYouTubeService()
	service.HandleError(serviceError, "Error creating YouTube client")
	channelMetaInfo := service.ChannelMetaInfo{}
	channelMetaInfo.CustomURL = opts.CustomURL
	fmt.Println(fmt.Sprintf("Initial input channelMetaInfo: %+v", channelMetaInfo))
	results := service.Exfoliator(youtubeService, channelMetaInfo)
	fmt.Println(fmt.Sprintf("Exfoliator exfoliated successfully. Results: %+v", results))
	elapsed := time.Since(start)
	fmt.Println(fmt.Sprintf("Finishing youtube-tinfoil-expose @ %v. Overall time spent: %v ms.", time.Now().Format(time.RFC3339), elapsed))
}
