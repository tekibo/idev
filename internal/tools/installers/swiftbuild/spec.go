// Package swiftbuild builds the pinned Swift Build source and the Expo engine.
package swiftbuild

import (
	"github.com/tekibo/idev/internal/platform/config"
	"github.com/tekibo/idev/internal/platform/release"
)

// Spec is the Swift Build and Expo engine install plan.
type Spec struct {
	Version      string
	Source       string
	SourceRoot   string
	SwiftBin     string
	SwiftRoot    string
	EngineSource string
	Bin          string
	EngineBin    string
	Bundles      string
	Component    release.Component
}

// NewSpec returns the Swift Build and Expo engine paths.
func NewSpec(cfg config.Config, component release.Component) Spec {
	return Spec{
		Version:      component.Version,
		Source:       cfg.Sources.Release.Manifest,
		SourceRoot:   cfg.Path("cache", "src", "swift-build-"+component.Version),
		SwiftBin:     cfg.SwiftBin(),
		SwiftRoot:    cfg.SwiftRoot(),
		EngineSource: cfg.Path("react-native", "engine"),
		Bin:          cfg.Path("bin", "SWBBuildService"),
		EngineBin:    cfg.Path("bin", "expo-build-bin"),
		Bundles:      cfg.Path("bin"),
		Component:    component,
	}
}
