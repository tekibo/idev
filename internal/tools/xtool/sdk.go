package xtool

import (
	"io"
	"strings"
)

// ParseSDKPath extracts the SDK bundle path from xtool sdk status output.
func ParseSDKPath(raw string) string {
	for _, line := range strings.Split(raw, "\n") {
		key, value, ok := strings.Cut(line, ":")
		if ok && strings.EqualFold(strings.TrimSpace(key), "Path") {
			return strings.TrimSpace(value)
		}
	}
	return ""
}

// SDKStatus runs xtool sdk status.
func (c *Client) SDKStatus(stdout, stderr io.Writer) error {
	return c.run([]string{"sdk", "status"}, nil, stdout, stderr)
}
