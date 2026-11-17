// Package agents installs the idev agent skill and instruction snippet.
package agents

import "embed"

//go:embed assets/skills/idev
var skillFS embed.FS

//go:embed assets/snippet.md
var snippet string

const (
	skillName   = "idev"
	beginMarker = "<!-- idev:begin -->"
	endMarker   = "<!-- idev:end -->"
	managedFile = ".idev-managed"
)
