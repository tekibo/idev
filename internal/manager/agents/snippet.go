package agents

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/tekibo/idev/internal/platform/config"
)

var snippetRe = regexp.MustCompile(`(?s)<!-- idev:begin -->.*?<!-- idev:end -->\n?`)

func instructionFiles() []string {
	home := config.UserHome()
	return []string{
		filepath.Join(home, ".claude", "CLAUDE.md"),
		filepath.Join(home, ".codex", "AGENTS.md"),
	}
}

func installSnippet() ([]string, error) {
	installed := []string{}
	for _, path := range instructionFiles() {
		if info, err := os.Stat(filepath.Dir(path)); err != nil || !info.IsDir() {
			continue
		}
		if err := backup(path); err != nil {
			return installed, err
		}
		old, err := os.ReadFile(path)
		if err != nil && !os.IsNotExist(err) {
			return installed, err
		}
		if err := os.WriteFile(path, []byte(withSnippet(string(old), snippet)), 0644); err != nil {
			return installed, err
		}
		installed = append(installed, path)
	}
	return installed, nil
}

func removeSnippet() ([]string, error) {
	removed := []string{}
	for _, path := range instructionFiles() {
		data, err := os.ReadFile(path)
		if err != nil {
			continue
		}
		if !snippetRe.Match(data) {
			continue
		}
		if err := backup(path); err != nil {
			return removed, err
		}
		if err := os.WriteFile(path, []byte(withoutSnippet(string(data))), 0644); err != nil {
			return removed, err
		}
		removed = append(removed, path)
	}
	return removed, nil
}

func backup(path string) error {
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return err
	}
	return os.WriteFile(path+".bak-idev", data, 0644)
}

func withSnippet(old, body string) string {
	block := beginMarker + "\n" + strings.TrimSpace(body) + "\n" + endMarker
	if snippetRe.MatchString(old) {
		return snippetRe.ReplaceAllString(old, block+"\n")
	}
	if strings.TrimSpace(old) == "" {
		return block + "\n"
	}
	return strings.TrimRight(old, "\n") + "\n\n" + block + "\n"
}

func withoutSnippet(old string) string {
	text := snippetRe.ReplaceAllString(old, "")
	return strings.TrimRight(text, "\n") + "\n"
}
