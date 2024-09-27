package manifest

import (
	"time"

	"github.com/coreos/go-semver/semver"
)

// Manifest represents the metadata information related to an application.
type Manifest struct {
	Application string
	Version     semver.Version
	CreatedAt   time.Time
}

// New creates and returns a new Manifest instance initialized with the given application name and version.
func New(name string, versionRaw string) *Manifest {
	version, err := semver.NewVersion(versionRaw)
	if err != nil {
		panic(err)
	}
	return &Manifest{
		Application: name,
		Version:     *version,
		CreatedAt:   time.Now(),
	}
}
