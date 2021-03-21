package version

import "runtime"

var (
	// The current built version
	version = ""
	// The git branch associated with the current built version
	gitBranch = ""
	// The git SHA1 of the commit
	gitCommit = ""
)

// BuildInfo contains build time information about the application
type BuildInfo struct {
	Version   string `json:"version,omitempty"`
	GitBranch string `json:"gitBranch,omitempty"`
	GitCommit string `json:"gitCommit,omitempty"`
	GoVersion string `json:"goVersion,omitempty"`
}

// Short returns the semantic version of the application
func Short() string {
	return version
}

// Long returns the build time version information of the application
func Long() BuildInfo {
	return BuildInfo{
		Version:   version,
		GitBranch: gitBranch,
		GitCommit: gitCommit,
		GoVersion: runtime.Version(),
	}
}
