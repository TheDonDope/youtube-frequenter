package api

import "time"

const (

	// MaxResultsUploadedVideos is the maximum search result number for the initial PlaylistItems.List call for the uploaded playlist
	MaxResultsUploadedVideos = 50

	// MaxResultsCommentPerVideo is the maximum search result number for the comments per video to be searched
	MaxResultsCommentPerVideo = 10

	// MaxResultsFavouritedVideos is the maximum search result number for the number of playlist items in the favourited playlist
	MaxResultsFavouritedVideos = 10

	// AverageAPICallDuration is the duration we estimate to average for a single API call
	AverageAPICallDuration = 8 * time.Millisecond

	// GlobalTimeout is the timeout for the complete program
	GlobalTimeout = MaxResultsUploadedVideos * MaxResultsCommentPerVideo * MaxResultsFavouritedVideos * AverageAPICallDuration
	//
)
