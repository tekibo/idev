// Package swift installs the official Swift toolchain privately.
package swift

import (
	"fmt"

	"github.com/tekibo/idev/internal/platform/config"
)

// Spec pins one official Swift tarball.
type Spec struct {
	Version string
	Arch    string
	URL     string
	SHA256  string
	Root    string
	Cache   string
}

// ForArch returns the official tarball for one architecture.
func ForArch(cfg config.Config, arch string) (Spec, error) {
	sources := cfg.Sources.Swift
	if sources.Version == "" {
		return Spec{}, fmt.Errorf("source config: Swift version is empty")
	}
	var archive config.SwiftArchive
	switch arch {
	case "amd64":
		archive = sources.AMD64
	case "arm64":
		archive = sources.ARM64
	default:
		return Spec{}, fmt.Errorf("unsupported Swift architecture %q", arch)
	}
	if archive.URL == "" || archive.SHA256 == "" {
		return Spec{}, fmt.Errorf("source config: Swift %s archive is empty", arch)
	}
	spec := Spec{
		Version: sources.Version,
		Arch:    arch,
		URL:     archive.URL,
		SHA256:  archive.SHA256,
		Root:    cfg.SwiftRoot(),
		Cache:   cfg.Path("cache", "downloads", "swift-"+sources.Version+"-"+arch+".tar.gz"),
	}
	return spec, nil
}
