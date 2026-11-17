package auth

import (
	"bytes"
	"strings"

	"github.com/tekibo/idev/internal/platform/redact"
)

// AuthStatus is the parsed account state plus redacted raw output.
type AuthStatus struct {
	SignedIn bool   `json:"signed_in"`
	Account  string `json:"account,omitempty"`
	Team     string `json:"team,omitempty"`
	Expiry   string `json:"expiry,omitempty"`
	Raw      string `json:"raw"`
	Parsed   bool   `json:"parsed"`
}

// AuthStatus reads and parses xtool auth status.
func (s *Service) AuthStatus() (AuthStatus, error) {
	var stdout, stderr bytes.Buffer
	if err := s.Auth([]string{"status"}, nil, &stdout, &stderr); err != nil {
		return AuthStatus{}, err
	}
	raw := stdout.String()
	if raw == "" {
		raw = stderr.String()
	}
	raw = redact.RedactText(raw)
	status := AuthStatus{Raw: raw}
	lower := strings.ToLower(raw)
	if strings.Contains(lower, "logged in") {
		status.SignedIn = true
		status.Parsed = true
	}
	if strings.Contains(lower, "not logged in") || strings.Contains(lower, "signed out") {
		status.SignedIn = false
		status.Parsed = true
	}
	status.Account = authField(raw, "account", "apple id", "email")
	status.Team = authField(raw, "team")
	status.Expiry = authField(raw, "expir")
	if status.Account != "" || status.Team != "" || status.Expiry != "" {
		status.Parsed = true
	}
	return status, nil
}

func authField(raw string, labels ...string) string {
	for _, line := range strings.Split(raw, "\n") {
		lower := strings.ToLower(line)
		for _, label := range labels {
			if strings.Contains(lower, label) {
				if _, value, ok := strings.Cut(line, ":"); ok {
					return strings.TrimSpace(value)
				}
			}
		}
	}
	return ""
}
