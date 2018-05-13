package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/TheDonDope/youtube-tinfoil-expose/api"
)

func main() {
	api.ParseArguments(os.Args)
	api.ConfigureOutput()
	api.ConfigureLogging()

	start := time.Now()
	api.Printfln("Starting youtube-tinfoil-expose @ %v", start.Format(time.RFC3339))

	serviceImpl := &api.YouTubeServiceImpl{}

	youtubeService, serviceError := serviceImpl.GetYouTubeService()
	if serviceError != nil {
		formattdErrorMessage := api.GetFormattedErrorMessage(serviceError, "Error creating YouTube client")
		if formattdErrorMessage != "" {
			log.Println(formattdErrorMessage)
		}
	}
	channelMetaInfo := api.CreateInitialChannelMetaInfo()

	results := api.Exfoliator(youtubeService, channelMetaInfo)
	log.Println("Exfoliator exfoliated successfully.")
	log.Println(fmt.Sprintf("Analysing Exfoliator results (ChannelID: %v, CustomURL: %v)", results.ChannelID, results.CustomURL))
	log.Println(fmt.Sprintf("#videos%v", len(results.ObviouslyRelatedChannelIDs)))
	api.AnalyseChannelMetaInfo(&results)
	api.Printfln("Finishing youtube-tinfoil-expose @ %v", time.Now().Format(time.RFC3339))
	api.Printfln("Overall time spent: %v", time.Since(start))
}
