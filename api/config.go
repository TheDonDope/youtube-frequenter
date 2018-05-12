package api

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
}
