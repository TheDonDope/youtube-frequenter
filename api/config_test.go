package api

import (
	"fmt"
	"log"
	"os"
	"testing"
)

const (
	expectedChannelID                  = "UCHmvQlYVHbkm4oKkgqezgIg"
	expectedCustomURL                  = "wwwKenFMde"
	expectedPlaylistID                 = "PLGAbee9wcbeJU9tpjArgmj2gxlGe_xxfj"
	expectedMaxResultsUploadedVideos   = "25"
	expectedMaxResultsCommentPerVideo  = "25"
	expectedMaxResultsFavouritedVideos = "25"
	expectedAverageAPICallDuration     = "10ms"
	expectedGlobalTimeout              = "60s"
)

func TestShouldParseWithDefaultValues(t *testing.T) {
	originalOpts := Opts
	defaultArguments := []string{}
	ParseArguments(defaultArguments)
	if actualMaxResultsUploadedVideos := fmt.Sprint(Opts.MaxResultsUploadedVideos); actualMaxResultsUploadedVideos != expectedMaxResultsUploadedVideos {
		t.Errorf(GetFormattedFailMessage("ShouldParseWithDefaultValues#MaxResultsUploadedVideos", expectedMaxResultsUploadedVideos, actualMaxResultsUploadedVideos))
	}
	if actualMaxResultsCommentPerVideo := fmt.Sprint(Opts.MaxResultsCommentPerVideo); actualMaxResultsCommentPerVideo != expectedMaxResultsCommentPerVideo {
		t.Errorf(GetFormattedFailMessage("ShouldParseWithDefaultValues#MaxResultsCommentPerVideo", expectedMaxResultsCommentPerVideo, actualMaxResultsCommentPerVideo))
	}
	if actualMaxResultsFavouritedVideos := fmt.Sprint(Opts.MaxResultsFavouritedVideos); actualMaxResultsFavouritedVideos != expectedMaxResultsFavouritedVideos {
		t.Errorf(GetFormattedFailMessage("ShouldParseWithDefaultValues#MaxResultsFavouritedVideos", expectedMaxResultsFavouritedVideos, actualMaxResultsFavouritedVideos))
	}
	if actualAverageAPICallDuration := Opts.AverageAPICallDuration; actualAverageAPICallDuration != expectedAverageAPICallDuration {
		t.Errorf(GetFormattedFailMessage("ShouldParseWithDefaultValues#AverageAPICallDuration", expectedAverageAPICallDuration, actualAverageAPICallDuration))
	}
	if actualGlobalTimeout := Opts.GlobalTimeout; actualGlobalTimeout != expectedGlobalTimeout {
		t.Errorf(GetFormattedFailMessage("ShouldParseWithDefaultValues#GlobalTimeout", expectedGlobalTimeout, actualGlobalTimeout))
	}
	Opts = originalOpts
}

func TestShouldParseShortArguments(t *testing.T) {
	originalOpts := Opts
	shortArguments := []string{"-c", expectedChannelID, "-u", expectedCustomURL, "-p", expectedPlaylistID}
	ParseArguments(shortArguments)
	if actualChannelID := Opts.ChannelID; actualChannelID != expectedChannelID {
		t.Errorf(GetFormattedFailMessage("ShouldParseShortArguments#ChannelID", expectedChannelID, actualChannelID))
	}
	if actualCustomURL := Opts.CustomURL; actualCustomURL != expectedCustomURL {
		t.Errorf(GetFormattedFailMessage("ShouldParseShortArguments#CustomURL", expectedCustomURL, actualCustomURL))
	}
	if actualPlaylistID := Opts.PlaylistID; actualPlaylistID != expectedPlaylistID {
		t.Errorf(GetFormattedFailMessage("ShouldParseShortArguments#PlaylistID", expectedPlaylistID, actualPlaylistID))
	}
	Opts = originalOpts
}

func TestShouldParseLongArguments(t *testing.T) {
	originalOpts := Opts
	longArguments := []string{"--channel-id", expectedChannelID, "--custom-url", expectedCustomURL, "--playlist-id", expectedPlaylistID}
	ParseArguments(longArguments)
	if actualChannelID := Opts.ChannelID; actualChannelID != expectedChannelID {
		t.Errorf(GetFormattedFailMessage("ShouldParseLongArguments#ChannelID", expectedChannelID, actualChannelID))
	}
	if actualCustomURL := Opts.CustomURL; actualCustomURL != expectedCustomURL {
		t.Errorf(GetFormattedFailMessage("ShouldParseLongArguments#CustomURL", expectedCustomURL, actualCustomURL))
	}
	if actualPlaylistID := Opts.PlaylistID; actualPlaylistID != expectedPlaylistID {
		t.Errorf(GetFormattedFailMessage("ShouldParseLongArguments#PlaylistID", expectedPlaylistID, actualPlaylistID))
	}
	Opts = originalOpts
}

func TestShouldConfigureLoggingWithChannelID(t *testing.T) {
	originalOpts := Opts
	// shortArgs
	channelIDShortArguments := []string{"-c", expectedChannelID}
	ParseArguments(channelIDShortArguments)
	ConfigureLogging()
	log.Println("hello")
	if _, err := os.Stat("logs/channel-id-" + expectedChannelID + ".log"); os.IsNotExist(err) {
		t.Errorf(GetFormattedFailMessage("ShouldConfigureLoggingWithChannelID", "Log file created.", err.Error()))
	}
	// longArgs
	channelIDLongArguments := []string{"--channel-id", expectedChannelID}
	ParseArguments(channelIDLongArguments)
	ConfigureLogging()
	if _, err := os.Stat("logs/channel-id-" + expectedChannelID + ".log"); os.IsNotExist(err) {
		t.Errorf(GetFormattedFailMessage("ShouldConfigureLoggingWithChannelID", "Log file created.", err.Error()))
	}
	Opts = originalOpts
}

func TestShouldConfigureLoggingWithCustomURL(t *testing.T) {
	originalOpts := Opts
	// shortArgs
	customURLShortArguments := []string{"-u", expectedCustomURL}
	ParseArguments(customURLShortArguments)
	ConfigureLogging()
	log.Println("hello")
	if _, err := os.Stat("logs/custom-url-" + expectedCustomURL + ".log"); os.IsNotExist(err) {
		t.Errorf(GetFormattedFailMessage("ShouldConfigureLoggingWithCustomURL", "Log file created.", err.Error()))
	}
	// longArgs
	customURLLongArguments := []string{"--custom-url", expectedCustomURL}
	ParseArguments(customURLLongArguments)
	ConfigureLogging()
	if _, err := os.Stat("logs/custom-url-" + expectedCustomURL + ".log"); os.IsNotExist(err) {
		t.Errorf(GetFormattedFailMessage("ShouldConfigureLoggingWithCustomURL", "Log file created.", err.Error()))
	}
	Opts = originalOpts
}

func TestShouldConfigureLoggingWithPlaylistID(t *testing.T) {
	originalOpts := Opts
	// shortArgs
	playlistIDShortArguments := []string{"-p", expectedPlaylistID}
	ParseArguments(playlistIDShortArguments)
	ConfigureLogging()
	log.Println("hello")
	if _, err := os.Stat("logs/playlist-id-" + expectedPlaylistID + ".log"); os.IsNotExist(err) {
		t.Errorf(GetFormattedFailMessage("ShouldConfigureLoggingWithPlaylistID", "Log file created.", err.Error()))
	}
	// longArgs
	playlistIDLongArguments := []string{"--playlist-id", expectedPlaylistID}
	ParseArguments(playlistIDLongArguments)
	ConfigureLogging()
	if _, err := os.Stat("logs/playlist-id-" + expectedPlaylistID + ".log"); os.IsNotExist(err) {
		t.Errorf(GetFormattedFailMessage("ShouldConfigureLoggingWithPlaylistID", "Log file created.", err.Error()))
	}
	Opts = originalOpts
}
