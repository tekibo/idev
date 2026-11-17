// Package manifest owns the runtime component lock.
package manifest

// Version is the manager version reported by idev version.
const Version = "0.1.0"

// Component is one installable private component.
type Component struct {
	ID      string `json:"id"`
	Version string `json:"version"`
	Kind    string `json:"kind"`
	Path    string `json:"path"`
}

// Lock is the runtime component lock.
type Lock struct {
	Schema     int         `json:"schema"`
	Components []Component `json:"components"`
}
