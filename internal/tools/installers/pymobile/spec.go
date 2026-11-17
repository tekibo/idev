// Package pymobile installs the private pymobiledevice3 venv.
package pymobile

import (
	"path/filepath"

	"github.com/tekibo/idev/internal/platform/config"
)

// Spec is the private venv install plan.
type Spec struct {
	Version string
	Root    string
	Bin     string
	Link    string
	Cache   string
	Compat  string
}

// NewSpec returns the private venv paths.
func NewSpec(cfg config.Config) Spec {
	return Spec{
		Version: cfg.Lock.Idev.Pymobiledevice3,
		Root:    cfg.PymobileRoot(),
		Bin:     cfg.PymobileBin(),
		Link:    cfg.PymobileLink(),
		Cache:   cfg.PymobileCache(),
		Compat:  filepath.Join(config.UserHome(), "pymobile3-venv"),
	}
}
