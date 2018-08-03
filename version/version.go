package version

import (
	"fmt"
)

var (
	// This will be filled in by the linker at build time.
	Version = "0.0.0"

	// This will be filled in by the linker at build time.
	GitCommit = "dev"
)

func GetHumanVersion() string {
	return fmt.Sprintf("%v-%v", Version, GitCommit)
}
