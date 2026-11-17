package agents

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/tekibo/idev/internal/platform/config"
)

func agentDirs() []string {
	home := config.UserHome()
	return []string{
		filepath.Join(home, ".claude", "skills"),
		filepath.Join(home, ".agents", "skills"),
	}
}

func linkSkill(skillDir string) ([]string, []string, error) {
	links := []string{}
	skipped := []string{}
	for _, dir := range agentDirs() {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return nil, nil, err
		}
		link := filepath.Join(dir, skillName)
		if info, err := os.Lstat(link); err == nil {
			if info.Mode()&os.ModeSymlink != 0 {
				target, readErr := os.Readlink(link)
				if readErr == nil && target == skillDir {
					if err := os.Remove(link); err != nil {
						return nil, nil, err
					}
				} else {
					skipped = append(skipped, link+" (foreign link)")
					continue
				}
			} else {
				skipped = append(skipped, link+" (not a link)")
				continue
			}
		}
		if err := os.Symlink(skillDir, link); err != nil {
			return nil, nil, fmt.Errorf("%s: %w", link, err)
		}
		links = append(links, link)
	}
	return links, skipped, nil
}

func removeLinks(skillDir string) ([]string, error) {
	removed := []string{}
	for _, dir := range agentDirs() {
		link := filepath.Join(dir, skillName)
		info, err := os.Lstat(link)
		if err != nil {
			continue
		}
		if info.Mode()&os.ModeSymlink == 0 {
			continue
		}
		target, err := os.Readlink(link)
		if err != nil || target != skillDir {
			continue
		}
		if err := os.Remove(link); err != nil {
			return removed, err
		}
		removed = append(removed, link)
	}
	return removed, nil
}
