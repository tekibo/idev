package xcode

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/tekibo/idev/internal/platform/config"
)

// Import moves one chosen xip into the idev-owned xcode directory.
// Mode is "add" or "replace"; replace removes the previous current file.
func Import(ctx context.Context, cfg config.Config, source, mode string, log func(string),
	progress func(done, total int64)) (List, error) {
	if mode != "add" && mode != "replace" {
		return List{}, fmt.Errorf("choose add or replace")
	}
	if source == "" {
		return List{}, fmt.Errorf("choose an Xcode.xip first")
	}
	info, err := os.Stat(source)
	if err != nil {
		return List{}, fmt.Errorf("Xcode.xip: %v", err)
	}
	if !info.Mode().IsRegular() || !strings.HasSuffix(strings.ToLower(source), ".xip") {
		return List{}, fmt.Errorf("choose an Xcode.xip file")
	}
	directory := cfg.Path("xcode")
	destination := filepath.Join(directory, filepath.Base(source))
	if filepath.Clean(source) == filepath.Clean(destination) {
		if mode == "replace" {
			before, _ := Files(cfg)
			if before.Current != "" && filepath.Clean(before.Current) != filepath.Clean(destination) {
				if err := os.Remove(before.Current); err != nil && !os.IsNotExist(err) {
					return List{}, fmt.Errorf("the old file could not be removed: %v", err)
				}
			}
		}
		return Select(cfg, destination)
	}
	if err := os.MkdirAll(directory, 0755); err != nil {
		return List{}, err
	}
	before, _ := Files(cfg)
	current := before.Current
	if mode == "add" {
		if _, statErr := os.Stat(destination); statErr == nil {
			return List{}, fmt.Errorf("%s already exists; use Replace current", filepath.Base(destination))
		}
	}
	if log != nil {
		log("moving " + source + " to " + destination)
	}
	if err := os.Rename(source, destination); err == nil {
		if log != nil {
			log("moved Xcode.xip")
		}
	} else if !errors.Is(err, syscall.EXDEV) {
		return List{}, err
	} else {
		if err := copyFile(ctx, source, destination, info.Size(), progress); err != nil {
			return List{}, err
		}
		if err := os.Remove(source); err != nil && log != nil {
			log("copied Xcode.xip; the source could not be removed: " + err.Error())
		}
		if log != nil {
			log("moved Xcode.xip")
		}
	}
	if mode == "replace" && current != "" && filepath.Clean(current) != filepath.Clean(destination) {
		if err := os.Remove(current); err != nil && !os.IsNotExist(err) {
			return List{}, fmt.Errorf("the new xip is installed; the old file could not be removed: %v", err)
		}
	}
	if err := config.SaveXcodeSelection(cfg.Home, config.XcodeSelection{XIP: filepath.Base(destination)}); err != nil {
		return List{}, err
	}
	return Files(cfg)
}
