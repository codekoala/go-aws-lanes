package version

import (
	"fmt"
	"runtime"
)

var (
	Version   = "v0.1.0"
	Tag       = "dev"
	BuildDate string
)

func String() string {
	return fmt.Sprintf("%s-%s", Version, Tag)
}

func Full() string {
	return fmt.Sprintf("%s; built: %s; %s", String(), BuildDate, runtime.Version())
}

func init() {
	if BuildDate == "" {
		BuildDate = "unknown"
	}
}
