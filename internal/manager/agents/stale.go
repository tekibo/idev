package agents

import (
	"bytes"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

func skillStale(dir string) bool {
	root := "assets/skills/" + skillName
	expected := map[string]bool{}
	stale := false
	_ = fs.WalkDir(skillFS, root, func(path string, entry fs.DirEntry, err error) error {
		if err != nil || entry.IsDir() {
			return nil
		}
		rel := strings.TrimPrefix(path, root+"/")
		expected[rel] = true
		want, err := fs.ReadFile(skillFS, path)
		if err != nil {
			stale = true
			return nil
		}
		got, err := os.ReadFile(filepath.Join(dir, filepath.FromSlash(rel)))
		if err != nil || !bytes.Equal(want, got) {
			stale = true
		}
		return nil
	})
	_ = filepath.WalkDir(dir, func(path string, entry fs.DirEntry, err error) error {
		if err != nil || entry.IsDir() {
			return nil
		}
		rel, relErr := filepath.Rel(dir, path)
		if relErr != nil {
			return nil
		}
		rel = filepath.ToSlash(rel)
		if rel != managedFile && !expected[rel] {
			stale = true
		}
		return nil
	})
	return stale
}
