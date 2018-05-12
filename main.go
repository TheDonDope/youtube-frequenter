package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/jessevdk/go-flags"

	"github.com/TheDonDope/youtube-tinfoil-expose/api"
)

func main() {
	_, argsError := flags.ParseArgs(&api.Opts, os.Args)
	if argsError != nil {
		panic(argsError)
	}

	//create your file with desired read/write permissions
	var logFileName string
	if api.Opts.ChannelID != "" {
		logFileName = "channel-id-" + api.Opts.ChannelID
	} else if api.Opts.CustomURL != "" {
		logFileName = "custom-url-" + api.Opts.CustomURL
	} else if api.Opts.PlaylistID != "" {
		logFileName = "playlist-id-" + api.Opts.PlaylistID
	}
	logFileName = logFileName + ".log"
	os.MkdirAll("logs", 0700)
	logFile, logFileError := os.OpenFile("logs/"+logFileName, os.O_WRONLY|os.O_CREATE, 0644)
	if logFileError != nil {
		log.Fatal(logFileError)
	}

	//set output of logs to f
	log.SetOutput(io.MultiWriter(os.Stdout, logFile))
	start := time.Now()
	api.Printfln("Starting youtube-tinfoil-expose @ %v", start.Format(time.RFC3339))

	youtubeService, serviceError := api.GetYouTubeService()
	if serviceError != nil {
		formattdErrorMessage := api.GetFormattedErrorMessage(serviceError, "Error creating YouTube client")
		if formattdErrorMessage != "" {
			log.Println(formattdErrorMessage)
		}
	}
	channelMetaInfo := api.ChannelMetaInfo{}
	if api.Opts.PlaylistID == "" {
		channelMetaInfo.ChannelID = api.Opts.ChannelID
		channelMetaInfo.CustomURL = api.Opts.CustomURL
		channelMetaInfo.NextOperation = api.GetChannelOverviewOperation
	} else {
		uploadedPlaylist := &api.Playlist{PlaylistID: api.Opts.PlaylistID}
		channelMetaInfo.Playlists = make(map[string]*api.Playlist)
		channelMetaInfo.Playlists["uploads"] = uploadedPlaylist
		channelMetaInfo.NextOperation = api.GetVideoIDsOverviewOperation
	}

	results := api.Exfoliator(youtubeService, channelMetaInfo)
	log.Println("Exfoliator exfoliated successfully.")
	log.Println(fmt.Sprintf("Analysing Exfoliator results (ChannelID: %v, CustomURL: %v)", results.ChannelID, results.CustomURL))
	log.Println(fmt.Sprintf("#videos%v", len(results.ObviouslyRelatedChannelIDs)))
	api.AnalyseChannelMetaInfo(&results)
	api.Printfln("Finishing youtube-tinfoil-expose @ %v", time.Now().Format(time.RFC3339))
	api.Printfln("Overall time spent: %v", time.Since(start))
	//defer to close when you're done with it, not because you think it's idiomatic!
	defer logFile.Close()
}
