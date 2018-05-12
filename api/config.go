package api

import "time"

const (

	// MaxResultsUploadedVideos is the maximum search result number for the initial PlaylistItems.List call for the uploaded playlist (0 - 50, default if not set: 5)
	MaxResultsUploadedVideos = 50

	// MaxResultsCommentPerVideo is the maximum search result number for the comments per video to be searched (0 - 100, default if not set: 20)
	MaxResultsCommentPerVideo = 100

	// MaxResultsFavouritedVideos is the maximum search result number for the number of playlist items in the favourited playlist (0 - 50, default if not set: 5)
	MaxResultsFavouritedVideos = 50

	// AverageAPICallDuration is the duration we estimate to average for a single API call
	AverageAPICallDuration = 10 * time.Millisecond

	// GlobalTimeout is the timeout for the complete program
	GlobalTimeout = 60 * time.Second
	//GlobalTimeout = MaxResultsUploadedVideos * MaxResultsCommentPerVideo * MaxResultsFavouritedVideos * AverageAPICallDuration

)
