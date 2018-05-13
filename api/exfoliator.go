package api

import (
	"encoding/json"
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
	serviceImpl := &YouTubeServiceImpl{}
	go func() {
		monoChannel <- channelMetaInfo
	}()
	go GetChannelOverview(service, serviceImpl, monoChannel)
	go GetVideoIDsOverview(service, serviceImpl, monoChannel)
	go GetCommentsOverview(service, serviceImpl, monoChannel)
	go GetObviouslyRelatedChannelsOverview(service, serviceImpl, monoChannel, lastButNotLeastChannel)

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

// AnalyseChannelMetaInfo prints additional information for a given channelMetaInfo.
func AnalyseChannelMetaInfo(channelMetaInfo *ChannelMetaInfo) {
	relatedChannelIDToNumberOfOccurrences := CountOccurrences(channelMetaInfo.ObviouslyRelatedChannelIDs)

	if len(relatedChannelIDToNumberOfOccurrences) == 0 {
		log.Println("Package to analyse has no ObviouslyRelatedChannelIDs to count.")
	} else {
		sortedRelatedChannelIDsList := RankByWordCount(relatedChannelIDToNumberOfOccurrences)

		resultJSONBytes, resultJSONBytesError := json.Marshal(sortedRelatedChannelIDsList)
		HandleError(resultJSONBytesError, "Error marshaling results")
		WriteToJSON(GetOutputDirectory()+"/"+GetCustomName()+"-results.json", resultJSONBytes)
		printResults(sortedRelatedChannelIDsList)

		dumpJSONBytes, dumpJSONBytesError := json.Marshal(channelMetaInfo)
		HandleError(dumpJSONBytesError, "Error marshaling dump")
		WriteToJSON(GetOutputDirectory()+"/"+GetCustomName()+"-dump.json", dumpJSONBytes)
	}
}

// printResults prints the results
func printResults(results MapEntryList) {
	for _, item := range results {
		log.Println(item)
	}
}
