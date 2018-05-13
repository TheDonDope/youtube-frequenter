package api

import "testing"

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

func TestShortArgumentsShouldParseWithDefaultValues(t *testing.T) {
	shortDefaultArguments := []string{"-c", expectedChannelID, "-u", expectedCustomURL, "-p", expectedPlaylistID}
	ParseArguments(shortDefaultArguments)
	if actualChannelID := Opts.ChannelID; actualChannelID != expectedChannelID {
		t.Errorf(GetFormattedFailMessage("ParseArguments#ChannelID", expectedChannelID, actualChannelID))
	}
	if actualCustomURL := Opts.CustomURL; actualCustomURL != expectedCustomURL {
		t.Errorf(GetFormattedFailMessage("ParseArguments#CustomURL", expectedCustomURL, actualCustomURL))
	}
	if actualPlaylistID := Opts.PlaylistID; actualPlaylistID != expectedPlaylistID {
		t.Errorf(GetFormattedFailMessage("ParseArguments#PlaylistID", expectedPlaylistID, actualPlaylistID))
	}
}

func TestLongArgumentsShouldParseWithDefaultValues(t *testing.T) {
	shortDefaultArguments := []string{"--channel-id", expectedChannelID, "--custom-url", expectedCustomURL, "--playlist-id", expectedPlaylistID}
	ParseArguments(shortDefaultArguments)
	if actualChannelID := Opts.ChannelID; actualChannelID != expectedChannelID {
		t.Errorf(GetFormattedFailMessage("ParseArguments#ChannelID", expectedChannelID, actualChannelID))
	}
	if actualCustomURL := Opts.CustomURL; actualCustomURL != expectedCustomURL {
		t.Errorf(GetFormattedFailMessage("ParseArguments#CustomURL", expectedCustomURL, actualCustomURL))
	}
	if actualPlaylistID := Opts.PlaylistID; actualPlaylistID != expectedPlaylistID {
		t.Errorf(GetFormattedFailMessage("ParseArguments#PlaylistID", expectedPlaylistID, actualPlaylistID))
	}
}
