package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/TheDonDope/youtube-frequenter/api"
)

func main() {
	api.ParseArguments(os.Args)
	api.ConfigureOutput()
	logFile := api.ConfigureLogging()
	defer logFile.Close()

	start := time.Now()
	api.Printfln("Starting youtube-frequenter @ %v", start.Format(time.RFC3339))

	serviceImpl := &api.YouTuberService{}
	exfoliatorImpl := &api.ExfoliatorService{}

	youtubeService, serviceError := serviceImpl.GetService()
	if serviceError != nil {
		formattdErrorMessage := api.GetFormattedErrorMessage(serviceError, "Error creating YouTube client")
		if formattdErrorMessage != "" {
			log.Println(formattdErrorMessage)
		}
	}
	channelMetaInfo := exfoliatorImpl.CreateInitialChannelMetaInfo()

	results := exfoliatorImpl.Exfoliate(youtubeService, serviceImpl, channelMetaInfo)
	log.Println("Exfoliator exfoliated successfully.")
	log.Println(fmt.Sprintf("Analysing Exfoliator results (ChannelID: %v, CustomURL: %v)", results.ChannelID, results.CustomURL))
	log.Println(fmt.Sprintf("#videos%v", len(results.ObviouslyRelatedChannelIDs)))
	exfoliatorImpl.AnalyseChannelMetaInfo(&results)
	api.Printfln("Program arguments: %+v", api.Opts)
	api.Printfln("Finishing youtube-frequenter @ %v", time.Now().Format(time.RFC3339))
	api.Printfln("Overall time spent: %v", time.Since(start))
}
