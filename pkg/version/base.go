package version

import (
	"strings"
)

var (
	major      = "1"
	minor      = "0"
	patch      = "0"
	gitVersion = "v" + strings.Join([]string{major, minor, patch}, ".")
	gitCommit  = "aafda604f437cf1bc26d3db83aa567527613359e"
)
