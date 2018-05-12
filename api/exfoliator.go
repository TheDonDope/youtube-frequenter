package api

import (
	"log"
	"time"

	"google.golang.org/api/youtube/v3"
)

// Exfoliator exfoliates
func Exfoliator(service *youtube.Service, channelMetaInfo ChannelMetaInfo) ChannelMetaInfo {
	monoChannel := make(chan ChannelMetaInfo)
	lastButNotLeastChannel := make(chan ChannelMetaInfo)
	accumulatedMetaInfo := ChannelMetaInfo{}
	accumulatedMetaInfo.CustomURL = channelMetaInfo.CustomURL
	accumulatedMetaInfo.ChannelID = channelMetaInfo.ChannelID
	accumulatedMetaInfo.Playlists = channelMetaInfo.Playlists
	go func() {
		monoChannel <- channelMetaInfo
	}()
	go GetChannelOverview(service, monoChannel)
	go GetVideoIDsOverview(service, monoChannel)
	go GetCommentsOverview(service, monoChannel)
	go GetObviouslyRelatedChannelsOverview(service, monoChannel, lastButNotLeastChannel)

	globalTimeout, globalTimeoutError := time.ParseDuration(Opts.GlobalTimeout)
	if globalTimeoutError != nil {
		log.Println(globalTimeoutError)
	}
	timeout := time.After(globalTimeout)
	for {
		log.Println("<<<<<Begin Exfoliator Main Loop")
		select {
		case channelMetaInfo = <-lastButNotLeastChannel:
			log.Println("<-- (5/5): Exfoliator")
			log.Println("<-> (5/5): Working in Exfoliator")
			// evtl die anderen properties adden
			accumulatedMetaInfo.ObviouslyRelatedChannelIDs = append(accumulatedMetaInfo.ObviouslyRelatedChannelIDs, channelMetaInfo.ObviouslyRelatedChannelIDs...)
			log.Println("--> (5/5): Exfoliator")
		case <-timeout:
			Printfln("Request timed out (%v)", Opts.GlobalTimeout)
			return accumulatedMetaInfo
		}
	}
}
