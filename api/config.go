package api

import (
	"io"
	"log"
	"os"

	"github.com/jessevdk/go-flags"
)

// Opts are the program options, configurable by command line argument
var Opts struct {
	ChannelID string `short:"c" long:"channel-id" description:"The channel ID of a YouTube Channel."`

	CustomURL string `short:"u" long:"custom-url" description:"The custom URL of a YouTube Channel."`

	PlaylistID string `short:"p" long:"playlist-id" description:"The id of the playlist to search."`

	MaxResultsUploadedVideos int64 `short:"U" long:"max-results-uploaded-videos" description:"The maximum search result number for the initial PlaylistItems.List call for the uploaded playlist (default: 25, range: 0-50)" default:"25"`

	MaxResultsCommentPerVideo int64 `short:"C" long:"max-results-comments-per-video" description:"The maximum search result number for the comments per video to be searched (default: 25, range: 0-100)" default:"25"`

	MaxResultsFavouritedVideos int64 `short:"F" long:"max-results-favourited-videos" description:"The maximum search result number for the number of playlist items in the favourited playlist (default: 25, range: 0-50)" default:"25"`

	AverageAPICallDuration string `short:"d" long:"average-api-call-duration" description:"The duration we estimate to average for a single API call (default: 10ms, format: 1h10m10s)" default:"10ms"`

	GlobalTimeout string `short:"t" long:"global-timeout" description:"The timeout for the complete program (default: 60sec, format: 1h10m10s)" default:"60s"`

	OutputDirectory string `short:"o" long:"output-directory" description:"The output directory for the log file, results.json and dump.json (default: output)" default:"output"`
}

// ParseArguments parses the program arguments
func ParseArguments(args []string) {
	_, argsError := flags.ParseArgs(&Opts, args)
	if argsError != nil {
		panic(argsError)
	}
}

// ConfigureLogging configures the logging
func ConfigureLogging() *os.File {
	logFileName := GetCustomName() + ".log"
	logFile, logFileError := os.OpenFile(GetOutputDirectory()+"/"+logFileName, os.O_WRONLY|os.O_CREATE, 0644)
	HandleError(logFileError, "LogFileError!")

	//set output of logs to f
	log.SetOutput(io.MultiWriter(os.Stdout, logFile))
	//defer to close when you're done with it, not because you think it's idiomatic!
	return logFile
}

// ConfigureOutput creates the necessary output folders
func ConfigureOutput() {
	os.MkdirAll(Opts.OutputDirectory+"/"+GetCustomName(), 0700)
}

// GetCustomName returns the custom file/directory name
func GetCustomName() string {
	result := ""
	if Opts.ChannelID != "" {
		result = "channel-id-" + Opts.ChannelID
	} else if Opts.CustomURL != "" {
		result = "custom-url-" + Opts.CustomURL
	} else if Opts.PlaylistID != "" {
		result = "playlist-id-" + Opts.PlaylistID
	}
	return result
}

// GetOutputDirectory returns the complete path to the output directory
func GetOutputDirectory() string {
	return Opts.OutputDirectory + "/" + GetCustomName()
}
