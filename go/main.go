package main

import (
	"log"
	"os"
	"time"

	"github.com/jessevdk/go-flags"

	"github.com/TheDonDope/youtube-tinfoil-expose/go/api"
)

var opts struct {
	ChannelID string `short:"i" long:"channel-id" description:"The channel ID of a YouTube Channel"`

	CustomURL string `short:"u" long:"custom-url" description:"The custom URL of a YouTube Channel" required:"true"`

	PlaylistName string `short:"n" long:"playlist-name" description:"The name of the playlist to search. Currently only supports the following playlists: uploads, favorites"`
}

func main() {
	start := time.Now()
	api.Printfln("Starting youtube-tinfoil-expose @ %v", start.Format(time.RFC3339))

	_, argsError := flags.ParseArgs(&opts, os.Args)
	if argsError != nil {
		panic(argsError)
	}

	youtubeService, serviceError := api.GetYouTubeService()
	if serviceError != nil {
		formattdErrorMessage := api.GetFormattedErrorMessage(serviceError, "Error creating YouTube client")
		if formattdErrorMessage != "" {
			log.Println(formattdErrorMessage)
		}
	}
	channelMetaInfo := api.ChannelMetaInfo{}
	channelMetaInfo.CustomURL = opts.CustomURL
	results := api.Exfoliator(youtubeService, channelMetaInfo)
	api.Printfln("Exfoliator exfoliated successfully. Results: %v", results)
	log.Println("Analising Exfoliator results...")
	api.AnaliseChannelMetaInfo(&results)
	api.Printfln("Finishing youtube-tinfoil-expose @ %v", time.Now().Format(time.RFC3339))
	api.Printfln("Overall time spent: %v", time.Since(start))
}
