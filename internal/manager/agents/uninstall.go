package agents

import (
	"os"
	"path/filepath"

	"github.com/tekibo/idev/internal/platform/config"
)

// UninstallResult reports what the uninstall removed.
type UninstallResult struct {
	RemovedLinks        []string `json:"removed_links"`
	RemovedSkill        bool     `json:"removed_skill"`
	RemovedInstructions []string `json:"removed_instructions"`
}

// Uninstall removes the links, the managed skill, and the snippet sections.
func Uninstall(cfg config.Config) (UninstallResult, error) {
	skillDir := cfg.Path("skills", skillName)
	links, err := removeLinks(skillDir)
	if err != nil {
		return UninstallResult{}, err
	}
	result := UninstallResult{RemovedLinks: links}
	if _, err := os.Stat(filepath.Join(skillDir, managedFile)); err == nil {
		if err := os.RemoveAll(skillDir); err != nil {
			return result, err
		}
		result.RemovedSkill = true
	}
	instructions, err := removeSnippet()
	if err != nil {
		return result, err
	}
	result.RemovedInstructions = instructions
	return result, nil
}
