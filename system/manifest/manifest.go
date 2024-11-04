package manifest

import (
	"time"

	"github.com/coreos/go-semver/semver"
	"go.uber.org/fx"
)

// Timeouts represents the timeout which are used by the underlying DI system (fx)
type Timeouts struct {
	Start time.Duration
	Stop  time.Duration
}

// Manifest represents the metadata information related to an application.
type Manifest struct {
	Application string
	Version     semver.Version
	CreatedAt   time.Time
	Timeouts    Timeouts
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
		Timeouts: Timeouts{
			Start: fx.DefaultTimeout,
			Stop:  fx.DefaultTimeout,
		},
	}
}
