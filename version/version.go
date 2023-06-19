package version

import (
	"fmt"

	version "github.com/hashicorp/go-version"
)

const Version = "0.1.0"
const Prerelease = "beta"

var SemVer *version.Version

func init() {
	SemVer = version.Must(version.NewVersion(Version))
}

func String() string {
	if Prerelease != "" {
		return fmt.Sprintf("%s-%s", Version, Prerelease)
	}
	return Version
}

var ServerID = "echoserver/" + String()
