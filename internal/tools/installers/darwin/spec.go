// Package darwin installs the vendored Darwin tools and OAM into an xtool SDK.
package darwin

import (
	"path/filepath"

	"github.com/tekibo/idev/internal/platform/config"
)

// Spec is the payload and SDK target for one install.
type Spec struct {
	Payload     string
	SDKRoot     string
	SDKPath     string
	SwiftBin    string
	SwiftRoot   string
	ToolsBin    string
	MacToolsBin string
	Plugins     string
	OAMServer   string
}

// NewSpec returns the payload and SDK paths from the manager config.
func NewSpec(cfg config.Config) Spec {
	spec := Spec{
		Payload:   cfg.Path("apple-tools"),
		SDKRoot:   cfg.SDKRoot(),
		SDKPath:   cfg.SDKPath(),
		SwiftBin:  cfg.SwiftBin(),
		SwiftRoot: cfg.SwiftRoot(),
	}
	if spec.SDKPath == "" {
		return spec
	}
	spec.ToolsBin = filepath.Join(spec.SDKPath, "Platforms", "iPhoneOS.platform", "Developer", "usr", "bin")
	spec.MacToolsBin = filepath.Join(spec.SDKPath, "Platforms", "MacOSX.platform", "Developer", "usr", "bin")
	spec.Plugins = filepath.Join(spec.SDKPath, "Platforms", "iPhoneOS.platform", "Developer", "usr", "lib", "swift", "host", "plugins")
	if spec.SDKRoot != "" {
		spec.OAMServer = filepath.Join(spec.SDKRoot, "OpenAppleMacrosServer")
	}
	return spec
}
