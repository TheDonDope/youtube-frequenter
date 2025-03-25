package main

import (
	"fmt"
	"log"
	"os"

	"github.com/TheDonDope/youtube-frequenter/pkg/config"
	"github.com/TheDonDope/youtube-frequenter/pkg/exfoliating"
	"github.com/TheDonDope/youtube-frequenter/pkg/http/youtube"
)

func main() {
	// TODO: rewrite
	// cli.ParseArgs(&config.Opts, os.Args)
	os.MkdirAll(config.GetOutPath(), 0700)

	// TODO: rewrite
	// f := cli.NewLogFile(config.GetOutPath() + "/" + config.GetOutName() + ".log")
	// defer f.Close()

	// TODO: rewrite
	// start := time.Now()
	// logging.Printfln("Starting youtube-frequenter @ %v", start.Format(time.RFC3339))

	ytV3, _ := youtube.NewYouTubeV3()
	// TODO: rewrite
	// errors.Print(err, "Error creating YouTubeV3 Service")

	yt := youtube.NewService(ytV3)
	xf := exfoliating.NewService(yt)

	info := xf.CreateInitialChannelMetaInfo()

	results := xf.Exfoliate(info)
	log.Println("Exfoliator exfoliated successfully.")
	log.Println(fmt.Sprintf("Analysing Exfoliator results (ChannelID: %v, CustomURL: %v)", results.ChannelID, results.CustomURL))
	log.Println(fmt.Sprintf("#videos%v", len(results.ObviouslyRelatedChannelIDs)))
	xf.AnalyseChannelMetaInfo(results)
	// TODO: rewrite
	// logging.Printfln("Program arguments: %+v", config.Opts)
	// logging.Printfln("Finishing youtube-frequenter @ %v", time.Now().Format(time.RFC3339))
	// logging.Printfln("Overall time spent: %v", time.Since(start))
}
