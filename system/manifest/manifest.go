package manifest

import (
	"time"

	"github.com/coreos/go-semver/semver"
)

type Manifest struct {
	Application string
	Version     semver.Version
	CreatedAt   time.Time
}

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
