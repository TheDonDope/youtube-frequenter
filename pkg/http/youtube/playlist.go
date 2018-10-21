package youtube

// Playlist is the struct for information on playlists
type Playlist struct {
	PlaylistID    string
	PlaylistItems []*Video
}
