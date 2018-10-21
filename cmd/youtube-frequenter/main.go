package main

import (
	"fmt"
	"log"
	"os"
	"time"

	cli "gitlab.com/TheDonDope/gocha/v3/pkg/config"
	"gitlab.com/TheDonDope/gocha/v3/pkg/errors"
	"gitlab.com/TheDonDope/gocha/v3/pkg/logging"
	"gitlab.com/TheDonDope/youtube-frequenter/pkg/config"
	"gitlab.com/TheDonDope/youtube-frequenter/pkg/exfoliating"
	"gitlab.com/TheDonDope/youtube-frequenter/pkg/http/youtube"
)

func main() {
	cli.ParseArgs(&config.Opts, os.Args)
	os.MkdirAll(config.GetOutPath(), 0700)
	f := cli.NewLogFile(config.GetOutPath() + "/" + config.GetOutName() + ".log")
	defer f.Close()

	start := time.Now()
	logging.Printfln("Starting youtube-frequenter @ %v", start.Format(time.RFC3339))

	ytV3, err := youtube.NewYouTubeV3()
	errors.Print(err, "Error creating YouTubeV3 Service")

	yt := youtube.NewService(ytV3)
	xf := exfoliating.NewService(yt)

	info := xf.CreateInitialChannelMetaInfo()

	results := xf.Exfoliate(info)
	log.Println("Exfoliator exfoliated successfully.")
	log.Println(fmt.Sprintf("Analysing Exfoliator results (ChannelID: %v, CustomURL: %v)", results.ChannelID, results.CustomURL))
	log.Println(fmt.Sprintf("#videos%v", len(results.ObviouslyRelatedChannelIDs)))
	xf.AnalyseChannelMetaInfo(results)
	logging.Printfln("Program arguments: %+v", config.Opts)
	logging.Printfln("Finishing youtube-frequenter @ %v", time.Now().Format(time.RFC3339))
	logging.Printfln("Overall time spent: %v", time.Since(start))
}
