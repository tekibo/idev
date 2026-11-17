package xcode

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/tekibo/idev/internal/platform/config"
)

// Validate reports whether one path is a regular xip inside the idev folder.
func Validate(cfg config.Config, path string) error {
	directory := filepath.Clean(cfg.Path("xcode"))
	cleaned := filepath.Clean(path)
	if filepath.Dir(cleaned) != directory {
		return fmt.Errorf("the xip is outside the idev xcode folder")
	}
	info, err := os.Lstat(cleaned)
	if err != nil {
		return fmt.Errorf("Xcode.xip: %v", err)
	}
	if !info.Mode().IsRegular() || !strings.HasSuffix(strings.ToLower(cleaned), ".xip") {
		return fmt.Errorf("choose an Xcode.xip file")
	}
	return nil
}

// Select records one existing xip as current.
func Select(cfg config.Config, path string) (List, error) {
	if err := Validate(cfg, path); err != nil {
		return List{}, err
	}
	if err := config.SaveXcodeSelection(cfg.Home, config.XcodeSelection{XIP: filepath.Base(path)}); err != nil {
		return List{}, err
	}
	return Files(cfg)
}

// Remove deletes one xip and selects the next file.
func Remove(cfg config.Config, path string) (List, error) {
	if err := Validate(cfg, path); err != nil {
		return List{}, err
	}
	if err := os.Remove(path); err != nil {
		return List{}, err
	}
	list, err := Files(cfg)
	if err != nil {
		return List{}, err
	}
	current := ""
	if len(list.Files) > 0 {
		current = list.Files[0].Path
	}
	name := ""
	if current != "" {
		name = filepath.Base(current)
	}
	if err := config.SaveXcodeSelection(cfg.Home, config.XcodeSelection{XIP: name}); err != nil {
		return List{}, err
	}
	list.Current = current
	markCurrent(&list)
	return list, nil
}
