package uninstall

import (
	"github.com/tekibo/idev/internal/platform/config"
)

// Component is one managed install unit.
type Component struct {
	ID        string   `json:"id"`
	Label     string   `json:"label"`
	Keepable  bool     `json:"keepable"`
	Installed bool     `json:"installed"`
	Keep      bool     `json:"keep"`
	Paths     []string `json:"paths"`
}

// spec describes one component's owned paths.
type spec struct {
	ID       string
	Label    string
	Keepable bool
	Paths    func(config.Config) []string
}

// specs returns the component map in report order.
func specs() []spec {
	return []spec{
		{ID: "swift", Label: "Swift toolchain and SourceKit-LSP", Keepable: true, Paths: swiftPaths},
		{ID: "xtool", Label: "xtool signer", Keepable: true, Paths: xtoolPaths},
		{ID: "pymobiledevice3", Label: "pymobiledevice3 device tool", Keepable: true, Paths: pymobilePaths},
		{ID: "cocoapods", Label: "CocoaPods gem tree", Keepable: true, Paths: cocoapodsPaths},
		{ID: "rcodesign", Label: "rcodesign binary", Keepable: true, Paths: binPaths("rcodesign")},
		{ID: "ipsw", Label: "ipsw binary", Keepable: true, Paths: binPaths("ipsw")},
		{ID: "swift-build", Label: "Swift Build and React Native engine", Paths: swiftBuildPaths},
		{ID: "react-native", Label: "React Native & Expo build engine", Paths: homePaths("react-native")},
		{ID: "flutter", Label: "Flutter scripts", Paths: homePaths("flutter")},
		{ID: "kmp", Label: "Kotlin Multiplatform scripts", Paths: homePaths("kmp")},
		{ID: "apple-tools", Label: "Darwin tools, OAM payload, and iPad viewer", Paths: appleToolsPaths},
		{ID: "agents", Label: "Agent skill", Paths: agentPaths},
		{ID: "sdk", Label: "Darwin SDK", Paths: sdkPaths},
		{ID: "state", Label: "idev state and logs", Paths: homePaths("state")},
		{ID: "xcode", Label: "Xcode.xip", Paths: homePaths("xcode")},
		{ID: "config", Label: "idev config and flow plan", Paths: homePaths("config")},
		{ID: "cache", Label: "idev caches", Paths: homePaths("cache")},
	}
}

// buildComponents returns the installed components with existing paths.
func buildComponents(cfg config.Config) []Component {
	var components []Component
	for _, item := range specs() {
		paths := existing(item.Paths(cfg))
		if len(paths) == 0 {
			continue
		}
		components = append(components, Component{
			ID:        item.ID,
			Label:     item.Label,
			Keepable:  item.Keepable,
			Installed: true,
			Paths:     paths,
		})
	}
	return components
}
