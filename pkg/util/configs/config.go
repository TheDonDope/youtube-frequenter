package configs

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
