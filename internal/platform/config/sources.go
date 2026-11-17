package config

import (
	_ "embed"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

//go:embed sources.json
var defaultSources []byte

// SwiftArchive is one official Swift tarball.
type SwiftArchive struct {
	URL    string `json:"url"`
	SHA256 string `json:"sha256"`
}

// SwiftSources is the pinned Swift source.
type SwiftSources struct {
	Version string       `json:"version"`
	AMD64   SwiftArchive `json:"amd64"`
	ARM64   SwiftArchive `json:"arm64"`
}

// ReleaseSources is the idev release manifest source.
type ReleaseSources struct {
	Manifest string `json:"manifest"`
}

// Sources is the user-facing source config.
type Sources struct {
	Release ReleaseSources `json:"release"`
	Swift   SwiftSources   `json:"swift"`
}

// LoadSources reads the user source config and seeds the default copy.
func LoadSources(home string) (Sources, error) {
	path := os.Getenv("IDEV_SOURCES")
	if path == "" && home != "" {
		path = filepath.Join(home, "config", "sources.json")
	}
	data := defaultSources
	if path != "" {
		current, err := os.ReadFile(path)
		switch {
		case err == nil:
			data = current
		case errors.Is(err, os.ErrNotExist):
			// Use the embedded default; setup writes the editable copy.
		default:
			return Sources{}, err
		}
	}
	var sources Sources
	if err := json.Unmarshal(data, &sources); err != nil {
		return Sources{}, err
	}
	if sources.Release.Manifest == "" {
		return Sources{}, errors.New("sources: release manifest is empty")
	}
	if sources.Swift.Version == "" {
		return Sources{}, errors.New("sources: Swift version is empty")
	}
	return sources, nil
}

// EnsureSources writes the editable user copy when it is absent.
func EnsureSources(home string) error {
	if home == "" || os.Getenv("IDEV_SOURCES") != "" {
		return nil
	}
	path := filepath.Join(home, "config", "sources.json")
	if _, err := os.Stat(path); err == nil {
		return nil
	}
	if err := writeDefault(path, defaultSources); err != nil {
		return err
	}
	return writeDefault(AppliedPath(home), defaultSources)
}

func writeDefault(path string, data []byte) error {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}
