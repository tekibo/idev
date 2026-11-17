package update

import (
	"fmt"
	"runtime"
	"sort"
	"strings"

	"github.com/tekibo/idev/internal/manager/core"
	"github.com/tekibo/idev/internal/platform/config"
	"github.com/tekibo/idev/internal/platform/manifest"
	"github.com/tekibo/idev/internal/platform/release"
)

// Update is one component update state.
type Update struct {
	Component string `json:"component"`
	Installed string `json:"installed"`
	Available string `json:"available"`
	Current   bool   `json:"current"`
	Warning   string `json:"warning,omitempty"`
}

// UpdateCheck compares installed components with the release manifest.
func (s *Service) UpdateCheck() ([]Update, error) {
	if s.Config.SourcesErr != nil {
		return nil, s.Config.SourcesErr
	}
	manifestFile, err := release.Load(s.Config.Sources.Release.Manifest)
	if err != nil {
		return nil, err
	}
	installed, err := core.LoadComponents(s.Config.Home)
	if err != nil {
		return nil, err
	}
	updates := []Update{{
		Component: "idev",
		Installed: manifest.Version,
		Available: manifestFile.Version,
		Current:   manifest.Version == manifestFile.Version,
	}}
	add := func(name, available string) {
		current := installed[name]
		updates = append(updates, Update{
			Component: name,
			Installed: current,
			Available: available,
			Current:   current == available,
		})
	}
	addManifest := func(name string) error {
		component, ok := manifestFile.Components[name]
		if !ok {
			return fmt.Errorf("release: component %s missing", name)
		}
		add(name, component.Version)
		return nil
	}
	addArch := func(name string) error {
		component, ok := manifestFile.Components[name+"-"+runtime.GOARCH]
		if !ok {
			return fmt.Errorf("release: component %s missing", name+"-"+runtime.GOARCH)
		}
		add(name, component.Version)
		return nil
	}
	for _, name := range []string{"react-native", "apple-tools", "flutter", "kmp", "xtool", "swift-build"} {
		if err := addManifest(name); err != nil {
			return nil, err
		}
	}
	if err := addArch("rcodesign"); err != nil {
		return nil, err
	}
	if err := addArch("ipsw"); err != nil {
		return nil, err
	}
	add("swift", s.Config.Sources.Swift.Version)
	add("pymobiledevice3", s.Config.Lock.Idev.Pymobiledevice3)
	add("cocoapods", s.Config.Lock.Idev.Cocoapods)
	for _, warning := range config.SourceWarnings(s.Config.Sources) {
		for i := range updates {
			if updates[i].Component == "swift" && strings.HasPrefix(warning.Field, "swift.") {
				updates[i].Warning = warning.Field + " " + warning.Current + " is unsupported; required " + warning.Required
			}
		}
	}
	sort.Slice(updates, func(i, j int) bool { return updates[i].Component < updates[j].Component })
	return updates, nil
}
