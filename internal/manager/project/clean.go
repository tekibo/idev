package project

import (
	"os"

	proj "github.com/tekibo/idev/internal/platform/project"
)

// CleanResult reports the removed build paths.
type CleanResult struct {
	Removed []string `json:"removed"`
}

// CleanBuild removes the project build output.
func (s *Service) CleanBuild(dir string, log func(string)) (CleanResult, error) {
	result := CleanResult{}
	info := proj.Inspect(dir)
	for _, path := range proj.CleanPaths(info) {
		if _, err := os.Stat(path); err != nil {
			continue
		}
		if err := os.RemoveAll(path); err != nil {
			return result, err
		}
		result.Removed = append(result.Removed, path)
		if log != nil {
			log("removed " + path)
		}
	}
	if len(result.Removed) == 0 && log != nil {
		log("no build output to remove")
	}
	return result, nil
}
