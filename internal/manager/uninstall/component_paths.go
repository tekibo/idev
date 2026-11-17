package uninstall

import (
	"os"
	"path/filepath"

	"github.com/tekibo/idev/internal/platform/config"
)

// swiftPaths returns the Swift toolchain, its links, and its cache root.
func swiftPaths(cfg config.Config) []string {
	paths := scanDirs(cfg.Path("toolchains"), "swift-")
	for _, name := range []string{"swift", "swiftc", "sourcekit-lsp"} {
		paths = append(paths, cfg.Path("bin", name))
	}
	return paths
}

// xtoolPaths returns every xtool build root and the manager link.
func xtoolPaths(cfg config.Config) []string {
	paths := scanDirs(cfg.Path("lib", "xtool"), "")
	return append(paths,
		cfg.Path("bin", "xtool"),
		filepath.Join(config.UserHome(), ".local", "bin", "xtool"),
	)
}

// pymobilePaths returns the private venv, its link, and the compat link.
func pymobilePaths(cfg config.Config) []string {
	paths := []string{cfg.PymobileRoot(), cfg.Path("bin", "pymobiledevice3")}
	if compat := pymobileCompat(cfg); compat != "" {
		paths = append(paths, compat)
	}
	return paths
}

// pymobileCompat returns the apple-dev venv link when idev owns it.
func pymobileCompat(cfg config.Config) string {
	link := filepath.Join(config.UserHome(), "pymobile3-venv")
	target, err := os.Readlink(link)
	if err != nil || target != cfg.PymobileRoot() {
		return ""
	}
	return link
}

// cocoapodsPaths returns the private gem tree and its wrapper.
func cocoapodsPaths(cfg config.Config) []string {
	return []string{cfg.GemRoot(), cfg.Path("bin", "pod")}
}

// swiftBuildPaths returns the build service, the Expo engine, and its bundles.
func swiftBuildPaths(cfg config.Config) []string {
	paths := []string{cfg.Path("bin", "SWBBuildService"), cfg.Path("bin", "expo-build-bin")}
	return append(paths, scanSuffix(cfg.Path("bin"), ".bundle")...)
}

// sdkPaths returns the one managed Darwin SDK bundle.
func sdkPaths(cfg config.Config) []string {
	return []string{cfg.Path(".config", "swiftpm", "swift-sdks", "darwin.artifactbundle")}
}

// appleToolsPaths returns the payload, the viewer link, and its user entries.
func appleToolsPaths(cfg config.Config) []string {
	return []string{
		cfg.Path("apple-tools"),
		cfg.Path("bin", "ipad-screen"),
		filepath.Join(config.UserHome(), ".local", "share", "applications", "idev-ipad-screen.desktop"),
		filepath.Join(config.UserHome(), ".config", "systemd", "user", "ipad-screen@.service"),
	}
}

// agentPaths returns the embedded skill and its agent links.
func agentPaths(cfg config.Config) []string {
	home := config.UserHome()
	return []string{
		cfg.Path("skills"),
		filepath.Join(home, ".claude", "skills", "idev"),
		filepath.Join(home, ".agents", "skills", "idev"),
	}
}

// binPaths returns one exact binary path.
func binPaths(name string) func(config.Config) []string {
	return func(cfg config.Config) []string {
		return []string{cfg.Path("bin", name)}
	}
}

// homePaths returns one exact private-home child.
func homePaths(name string) func(config.Config) []string {
	return func(cfg config.Config) []string {
		return []string{cfg.Path(name)}
	}
}
