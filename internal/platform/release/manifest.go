// Package release reads idev release manifests and installs component archives.
package release

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/tekibo/idev/internal/platform/download"
)

// Component is one release component.
type Component struct {
	Version   string `json:"version"`
	URL       string `json:"url"`
	SHA256    string `json:"sha256"`
	PinSource string `json:"pin_source,omitempty"`
}

// Manifest is the idev release manifest.
type Manifest struct {
	Version    string               `json:"version"`
	Components map[string]Component `json:"components"`
}

// Load reads a manifest from a URL or a local path.
func Load(source string) (Manifest, error) {
	data, err := read(source)
	if err != nil {
		return Manifest{}, err
	}
	var manifest Manifest
	if err := json.Unmarshal(data, &manifest); err != nil {
		return Manifest{}, err
	}
	if len(manifest.Components) == 0 {
		return Manifest{}, errors.New("release: manifest has no components")
	}
	return manifest, nil
}

// Resolve resolves a component URL against the manifest source.
func Resolve(source, reference string) string {
	if strings.HasPrefix(reference, "http://") ||
		strings.HasPrefix(reference, "https://") ||
		strings.HasPrefix(reference, "file://") ||
		filepath.IsAbs(reference) {
		return reference
	}
	if strings.HasPrefix(source, "http://") || strings.HasPrefix(source, "https://") {
		base, err := url.Parse(source)
		if err == nil {
			ref, err := url.Parse(reference)
			if err == nil {
				return base.ResolveReference(ref).String()
			}
		}
	}
	path := strings.TrimPrefix(source, "file://")
	return filepath.Join(filepath.Dir(path), reference)
}

func read(source string) ([]byte, error) {
	if strings.HasPrefix(source, "http://") || strings.HasPrefix(source, "https://") {
		data, err := download.Bytes(context.Background(), source)
		if err != nil {
			return nil, fmt.Errorf("release: %s: %w", source, err)
		}
		return data, nil
	}
	path := strings.TrimPrefix(source, "file://")
	return os.ReadFile(path)
}
