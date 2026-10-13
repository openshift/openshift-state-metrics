package version

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

var (
	// Release returns the release version
	Release = "UNKNOWN"
	// Commit returns the short sha from git
	Commit = "UNKNOWN"
	// BuildDate is the build date
	BuildDate = ""
)

type Version struct {
	GitCommit string
	BuildDate string
	Release   string
	GoVersion string
	Compiler  string
	Platform  string
}

func (v Version) String() string {
	return fmt.Sprintf("%s/%s (%s/%s) openshift-state-metrics/%s",
		filepath.Base(os.Args[0]), v.Release,
		runtime.GOOS, runtime.GOARCH, v.GitCommit)
}

// GetVersion returns openshift-state-metrics version
func GetVersion() Version {
	return Version{
		GitCommit: Commit,
		BuildDate: BuildDate,
		Release:   Release,
		GoVersion: runtime.Version(),
		Compiler:  runtime.Compiler,
		Platform:  fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
	}
}
