package configs

import (
	"io"
	"log"
	"os"

	flags "github.com/jessevdk/go-flags"
	"gitlab.com/TheDonDope/youtube-frequenter/pkg/util/errors"
)

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
	errors.HandleError(logFileError, "LogFileError!")

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
