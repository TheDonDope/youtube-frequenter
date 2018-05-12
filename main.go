package main

import (
	"io"
	"log"
	"os"
	"time"

	"github.com/jessevdk/go-flags"

	"github.com/TheDonDope/youtube-tinfoil-expose/api"
)

var opts struct {
	ChannelID string `short:"i" long:"channel-id" description:"The channel ID of a YouTube Channel"`

	CustomURL string `short:"u" long:"custom-url" description:"The custom URL of a YouTube Channel"`

	PlaylistName string `short:"n" long:"playlist-name" description:"The name of the playlist to search. Currently only supports the following playlists: uploads, favorites"`
}

func main() {
	//create your file with desired read/write permissions
	logFile, logFileError := os.OpenFile("exfoliate.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if logFileError != nil {
		log.Fatal(logFileError)
	}

	//set output of logs to f
	log.SetOutput(io.MultiWriter(os.Stdout, logFile))
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
	channelMetaInfo.ChannelID = opts.ChannelID
	channelMetaInfo.CustomURL = opts.CustomURL
	results := api.Exfoliator(youtubeService, channelMetaInfo)
	log.Println("Exfoliator exfoliated successfully.")
	log.Println("Analysing Exfoliator results:")
	api.AnalyseChannelMetaInfo(&results)
	api.Printfln("Finishing youtube-tinfoil-expose @ %v", time.Now().Format(time.RFC3339))
	api.Printfln("Overall time spent: %v", time.Since(start))
	//defer to close when you're done with it, not because you think it's idiomatic!
	defer logFile.Close()
}
