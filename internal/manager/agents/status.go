package agents

import (
	"os"
	"path/filepath"

	"github.com/tekibo/idev/internal/platform/config"
)

// LinkStatus reports one agent skill link.
type LinkStatus struct {
	Path   string `json:"path"`
	OK     bool   `json:"ok"`
	Detail string `json:"detail,omitempty"`
}

// InstructionStatus reports one global instruction file.
type InstructionStatus struct {
	Path    string `json:"path"`
	Present bool   `json:"present"`
}

// Status reports the installed skill and links.
type Status struct {
	SkillDir     string              `json:"skill_dir"`
	Installed    bool                `json:"installed"`
	Stale        bool                `json:"stale"`
	Links        []LinkStatus        `json:"links"`
	Instructions []InstructionStatus `json:"instructions"`
}

// Inspect reads the agent install state.
func Inspect(cfg config.Config) Status {
	skillDir := cfg.Path("skills", skillName)
	status := Status{SkillDir: skillDir}
	if _, err := os.Stat(filepath.Join(skillDir, managedFile)); err == nil {
		status.Installed = true
		status.Stale = skillStale(skillDir)
	}
	for _, dir := range agentDirs() {
		link := filepath.Join(dir, skillName)
		item := LinkStatus{Path: link}
		info, lstatErr := os.Lstat(link)
		target, err := os.Readlink(link)
		switch {
		case lstatErr == nil && info.Mode()&os.ModeSymlink == 0:
			item.Detail = "not a link"
		case err == nil && target == skillDir:
			item.OK = true
		case err == nil:
			item.Detail = "points to " + target
		case os.IsNotExist(err):
			item.Detail = "missing"
		default:
			item.Detail = err.Error()
		}
		status.Links = append(status.Links, item)
	}
	for _, path := range instructionFiles() {
		data, err := os.ReadFile(path)
		status.Instructions = append(status.Instructions, InstructionStatus{
			Path:    path,
			Present: err == nil && snippetRe.Match(data),
		})
	}
	return status
}
