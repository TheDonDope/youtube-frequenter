package test

import (
	"fmt"
	"os"
	"testing"

	"gitlab.com/TheDonDope/youtube-frequenter/pkg/util/configs"
	"gitlab.com/TheDonDope/youtube-frequenter/pkg/util/errors"
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
	originalOpts := configs.Opts
	defaultArguments := []string{}
	configs.ParseArguments(defaultArguments)
	if actualMaxResultsUploadedVideos := fmt.Sprint(configs.Opts.MaxResultsUploadedVideos); actualMaxResultsUploadedVideos != expectedMaxResultsUploadedVideos {
		t.Errorf(errors.GetFormattedFailMessage("ShouldParseWithDefaultValues#MaxResultsUploadedVideos", expectedMaxResultsUploadedVideos, actualMaxResultsUploadedVideos))
	}
	if actualMaxResultsCommentPerVideo := fmt.Sprint(configs.Opts.MaxResultsCommentPerVideo); actualMaxResultsCommentPerVideo != expectedMaxResultsCommentPerVideo {
		t.Errorf(errors.GetFormattedFailMessage("ShouldParseWithDefaultValues#MaxResultsCommentPerVideo", expectedMaxResultsCommentPerVideo, actualMaxResultsCommentPerVideo))
	}
	if actualMaxResultsFavouritedVideos := fmt.Sprint(configs.Opts.MaxResultsFavouritedVideos); actualMaxResultsFavouritedVideos != expectedMaxResultsFavouritedVideos {
		t.Errorf(errors.GetFormattedFailMessage("ShouldParseWithDefaultValues#MaxResultsFavouritedVideos", expectedMaxResultsFavouritedVideos, actualMaxResultsFavouritedVideos))
	}
	if actualAverageAPICallDuration := configs.Opts.AverageAPICallDuration; actualAverageAPICallDuration != expectedAverageAPICallDuration {
		t.Errorf(errors.GetFormattedFailMessage("ShouldParseWithDefaultValues#AverageAPICallDuration", expectedAverageAPICallDuration, actualAverageAPICallDuration))
	}
	if actualGlobalTimeout := configs.Opts.GlobalTimeout; actualGlobalTimeout != expectedGlobalTimeout {
		t.Errorf(errors.GetFormattedFailMessage("ShouldParseWithDefaultValues#GlobalTimeout", expectedGlobalTimeout, actualGlobalTimeout))
	}
	configs.Opts = originalOpts
	os.RemoveAll(configs.Opts.OutputDirectory)
}

func TestShouldParseShortArguments(t *testing.T) {
	originalOpts := configs.Opts
	shortArguments := []string{"-c", expectedChannelID, "-u", expectedCustomURL, "-p", expectedPlaylistID}
	configs.ParseArguments(shortArguments)
	if actualChannelID := configs.Opts.ChannelID; actualChannelID != expectedChannelID {
		t.Errorf(errors.GetFormattedFailMessage("ShouldParseShortArguments#ChannelID", expectedChannelID, actualChannelID))
	}
	if actualCustomURL := configs.Opts.CustomURL; actualCustomURL != expectedCustomURL {
		t.Errorf(errors.GetFormattedFailMessage("ShouldParseShortArguments#CustomURL", expectedCustomURL, actualCustomURL))
	}
	if actualPlaylistID := configs.Opts.PlaylistID; actualPlaylistID != expectedPlaylistID {
		t.Errorf(errors.GetFormattedFailMessage("ShouldParseShortArguments#PlaylistID", expectedPlaylistID, actualPlaylistID))
	}
	configs.Opts = originalOpts
	os.RemoveAll(configs.Opts.OutputDirectory)
}

func TestShouldParseLongArguments(t *testing.T) {
	originalOpts := configs.Opts
	longArguments := []string{"--channel-id", expectedChannelID, "--custom-url", expectedCustomURL, "--playlist-id", expectedPlaylistID}
	configs.ParseArguments(longArguments)
	if actualChannelID := configs.Opts.ChannelID; actualChannelID != expectedChannelID {
		t.Errorf(errors.GetFormattedFailMessage("ShouldParseLongArguments#ChannelID", expectedChannelID, actualChannelID))
	}
	if actualCustomURL := configs.Opts.CustomURL; actualCustomURL != expectedCustomURL {
		t.Errorf(errors.GetFormattedFailMessage("ShouldParseLongArguments#CustomURL", expectedCustomURL, actualCustomURL))
	}
	if actualPlaylistID := configs.Opts.PlaylistID; actualPlaylistID != expectedPlaylistID {
		t.Errorf(errors.GetFormattedFailMessage("ShouldParseLongArguments#PlaylistID", expectedPlaylistID, actualPlaylistID))
	}
	configs.Opts = originalOpts
	os.RemoveAll(configs.Opts.OutputDirectory)
}

func TestShouldConfigureLoggingWithChannelID(t *testing.T) {
	originalOpts := configs.Opts
	// shortArgs
	channelIDShortArguments := []string{"-c", expectedChannelID}
	configs.ParseArguments(channelIDShortArguments)
	configs.ConfigureOutput()
	configs.ConfigureLogging()
	if _, err := os.Stat(configs.GetOutputDirectory() + "/channel-id-" + expectedChannelID + ".log"); os.IsNotExist(err) {
		t.Errorf(errors.GetFormattedFailMessage("ShouldConfigureLoggingWithChannelID", "Log file created.", err.Error()))
	}
	// longArgs
	channelIDLongArguments := []string{"--channel-id", expectedChannelID}
	configs.ParseArguments(channelIDLongArguments)
	configs.ConfigureOutput()
	configs.ConfigureLogging()
	if _, err := os.Stat(configs.GetOutputDirectory() + "/channel-id-" + expectedChannelID + ".log"); os.IsNotExist(err) {
		t.Errorf(errors.GetFormattedFailMessage("ShouldConfigureLoggingWithChannelID", "Log file created.", err.Error()))
	}
	configs.Opts = originalOpts
	os.RemoveAll(configs.Opts.OutputDirectory)
}

func TestShouldConfigureLoggingWithCustomURL(t *testing.T) {
	originalOpts := configs.Opts
	// shortArgs
	customURLShortArguments := []string{"-u", expectedCustomURL}
	configs.ParseArguments(customURLShortArguments)
	configs.ConfigureOutput()
	configs.ConfigureLogging()
	if _, err := os.Stat(configs.GetOutputDirectory() + "/custom-url-" + expectedCustomURL + ".log"); os.IsNotExist(err) {
		t.Errorf(errors.GetFormattedFailMessage("ShouldConfigureLoggingWithCustomURL", "Log file created.", err.Error()))
	}
	// longArgs
	customURLLongArguments := []string{"--custom-url", expectedCustomURL}
	configs.ParseArguments(customURLLongArguments)
	configs.ConfigureOutput()
	configs.ConfigureLogging()
	if _, err := os.Stat(configs.GetOutputDirectory() + "/custom-url-" + expectedCustomURL + ".log"); os.IsNotExist(err) {
		t.Errorf(errors.GetFormattedFailMessage("ShouldConfigureLoggingWithCustomURL", "Log file created.", err.Error()))
	}
	configs.Opts = originalOpts
	os.RemoveAll(configs.Opts.OutputDirectory)
}

func TestShouldConfigureLoggingWithPlaylistID(t *testing.T) {
	originalOpts := configs.Opts
	// shortArgs
	playlistIDShortArguments := []string{"-p", expectedPlaylistID}
	configs.ParseArguments(playlistIDShortArguments)
	configs.ConfigureOutput()
	configs.ConfigureLogging()
	if _, err := os.Stat(configs.GetOutputDirectory() + "/playlist-id-" + expectedPlaylistID + ".log"); os.IsNotExist(err) {
		t.Errorf(errors.GetFormattedFailMessage("ShouldConfigureLoggingWithPlaylistID", "Log file created.", err.Error()))
	}
	// longArgs
	playlistIDLongArguments := []string{"--playlist-id", expectedPlaylistID}
	configs.ParseArguments(playlistIDLongArguments)
	configs.ConfigureOutput()
	configs.ConfigureLogging()
	if _, err := os.Stat(configs.GetOutputDirectory() + "/playlist-id-" + expectedPlaylistID + ".log"); os.IsNotExist(err) {
		t.Errorf(errors.GetFormattedFailMessage("ShouldConfigureLoggingWithPlaylistID", "Log file created.", err.Error()))
	}
	configs.Opts = originalOpts
	os.RemoveAll(configs.Opts.OutputDirectory)
}
