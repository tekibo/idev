package manager

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"

	"github.com/tekibo/idev/internal/tools/xtool"
)

// SDKOptions selects one SDK operation.
type SDKOptions struct {
	XIP  string `json:"xip,omitempty"`
	Out  string `json:"out,omitempty"`
	Arch string `json:"arch,omitempty"`
	Slim bool   `json:"slim"`
}

// SDKStatus is the parsed SDK state plus redacted raw output.
type SDKStatus struct {
	Installed bool   `json:"installed"`
	Version   string `json:"version,omitempty"`
	Flavor    string `json:"flavor,omitempty"`
	Path      string `json:"path,omitempty"`
	Raw       string `json:"raw"`
	Parsed    bool   `json:"parsed"`
}

var sdkVersion = regexp.MustCompile(`\d+\.\d+(\.\d+)?`)

// SDKStatus reads and parses xtool sdk status.
func (a *App) SDKStatus() (SDKStatus, error) {
	if !a.Xtool.Available() {
		return SDKStatus{}, fmt.Errorf("xtool not found; run idev setup")
	}
	var stdout, stderr bytes.Buffer
	if err := a.Xtool.SDKStatus(&stdout, &stderr); err != nil {
		return SDKStatus{}, err
	}
	raw := stdout.String()
	if raw == "" {
		raw = stderr.String()
	}
	raw = RedactText(raw)
	status := SDKStatus{Raw: raw}
	lower := strings.ToLower(raw)
	if strings.Contains(lower, "installed") && !strings.Contains(lower, "not installed") {
		status.Installed = true
		status.Parsed = true
	}
	status.Version = sdkVersion.FindString(raw)
	status.Path = xtool.ParseSDKPath(raw)
	for _, line := range strings.Split(raw, "\n") {
		if strings.HasPrefix(strings.ToLower(line), "flavor:") {
			if _, value, ok := strings.Cut(line, ":"); ok {
				status.Flavor = strings.TrimSpace(value)
			}
		}
	}
	if status.Version != "" || status.Path != "" {
		status.Parsed = true
	}
	return status, nil
}
