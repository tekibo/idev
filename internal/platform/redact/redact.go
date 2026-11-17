// Package redact masks secret values in command previews and output.
package redact

import (
	"regexp"
	"strings"
)

var secretFlags = map[string]bool{
	"--password":      true,
	"--token":         true,
	"--api-key":       true,
	"--session":       true,
	"--cookie":        true,
	"--private-key":   true,
	"--secret":        true,
	"--client-secret": true,
	"--access-token":  true,
	"--refresh-token": true,
	"--auth-token":    true,
}

// HasSecretFlag reports a secret value passed in argv.
func HasSecretFlag(args []string) bool {
	for _, arg := range args {
		if secretFlags[arg] {
			return true
		}
		for flag := range secretFlags {
			if strings.HasPrefix(arg, flag+"=") {
				return true
			}
		}
	}
	return false
}

var secretValue = regexp.MustCompile(`(?i)\b(token|password|secret|session|cookie)\b\s*[:=]\s*\S+`)

// RedactText masks secret values in captured output.
func RedactText(text string) string {
	return secretValue.ReplaceAllString(text, "$1=<redacted>")
}
