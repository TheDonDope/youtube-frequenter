package youtube

// Video is the struct for information on videos
type Video struct {
	VideoID           string
	UploaderChannelID string
	Comments          []*Comment
}
