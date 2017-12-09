package version

import (
	"fmt"
	"runtime"
)

var (
	Version   = "v0.4.1"
	Commit    = "dev"
	BuildDate string
)

func String() string {
	return fmt.Sprintf("%s-%s", Version, Commit)
}

func Full() string {
	return fmt.Sprintf("%s; built: %s; %s", String(), BuildDate, runtime.Version())
}

func init() {
	if BuildDate == "" {
		BuildDate = "unknown"
	}
}
