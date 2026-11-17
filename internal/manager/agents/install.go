package agents

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/tekibo/idev/internal/platform/config"
)

// Result reports one agent install.
type Result struct {
	SkillDir     string   `json:"skill_dir"`
	Links        []string `json:"links"`
	Instructions []string `json:"instructions"`
	Skipped      []string `json:"skipped,omitempty"`
}

// Install writes the embedded skill and links it into the agent directories.
func Install(cfg config.Config) (Result, error) {
	skillDir := cfg.Path("skills", skillName)
	if err := writeSkill(skillDir); err != nil {
		return Result{}, err
	}
	links, skipped, err := linkSkill(skillDir)
	if err != nil {
		return Result{}, err
	}
	instructions, err := installSnippet()
	if err != nil {
		return Result{}, err
	}
	return Result{
		SkillDir:     skillDir,
		Links:        links,
		Instructions: instructions,
		Skipped:      skipped,
	}, nil
}

func writeSkill(dst string) error {
	if info, err := os.Stat(dst); err == nil {
		if !info.IsDir() {
			return fmt.Errorf("%s exists and is not a directory", dst)
		}
		if _, err := os.Stat(filepath.Join(dst, managedFile)); err != nil {
			return fmt.Errorf("%s exists and is not idev-managed; move it first", dst)
		}
		if err := os.RemoveAll(dst); err != nil {
			return err
		}
	}
	root := "assets/skills/" + skillName
	err := fs.WalkDir(skillFS, root, func(path string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		rel := strings.TrimPrefix(path, root)
		rel = strings.TrimPrefix(rel, "/")
		target := dst
		if rel != "" {
			target = filepath.Join(dst, filepath.FromSlash(rel))
		}
		if entry.IsDir() {
			return os.MkdirAll(target, 0755)
		}
		data, err := fs.ReadFile(skillFS, path)
		if err != nil {
			return err
		}
		return os.WriteFile(target, data, 0644)
	})
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(dst, managedFile), []byte("managed by idev agents install\n"), 0644)
}
