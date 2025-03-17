// pkg/version/version.go
package version

import "fmt"

// Variables set during build time
var (
	Version   = "unknown" // Set via -ldflags during build
	GitCommit = "unknown" // Set via -ldflags during build
	BuildTime = "unknown" // Set via -ldflags during build
)

// GetVersionInfo returns a formatted string with version information.
func GetVersionInfo() string {
	return fmt.Sprintf("Version: %s, GitCommit: %s, BuildTime: %s", Version, GitCommit, BuildTime)
}
