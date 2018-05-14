package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/TheDonDope/youtube-frequenter/pkg/services"
	"github.com/TheDonDope/youtube-frequenter/pkg/util/configs"
	"github.com/TheDonDope/youtube-frequenter/pkg/util/errors"
	"github.com/TheDonDope/youtube-frequenter/pkg/util/logs"
)

func main() {
	configs.ParseArguments(os.Args)
	configs.ConfigureOutput()
	logFile := configs.ConfigureLogging()
	defer logFile.Close()

	start := time.Now()
	logs.Printfln("Starting youtube-frequenter @ %v", start.Format(time.RFC3339))

	serviceImpl := &services.YouTuberService{}
	exfoliatorImpl := &services.ExfoliatorService{}

	youtubeService, serviceError := serviceImpl.GetService()
	if serviceError != nil {
		formattdErrorMessage := errors.GetFormattedErrorMessage(serviceError, "Error creating YouTube client")
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
	logs.Printfln("Program arguments: %+v", configs.Opts)
	logs.Printfln("Finishing youtube-frequenter @ %v", time.Now().Format(time.RFC3339))
	logs.Printfln("Overall time spent: %v", time.Since(start))
}
