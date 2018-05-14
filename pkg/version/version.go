package version

// Get returns the overall codebase version. It's for detecting
// what code a binary was built from.
func Get() Info {
	return Info{
		Major:      major,
		Minor:      minor,
		Patch:      patch,
		GitVersion: gitVersion,
		GitCommit:  gitCommit,
	}
}
