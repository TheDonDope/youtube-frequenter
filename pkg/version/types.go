package version

// Info contains versioning information.
type Info struct {
	Major      string `json:"major"`
	Minor      string `json:"minor"`
	Patch      string `json:"patch"`
	GitVersion string `json:"gitVersion"`
	GitCommit  string `json:"gitCommit"`
}

// String returns info as a human-friendly version string.
func (info Info) String() string {
	return info.GitVersion
}
