// Package xcode owns the Xcode.xip folder.
package xcode

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/tekibo/idev/internal/platform/config"
)

// XIP is one xip file in the idev-owned folder.
type XIP struct {
	Name     string    `json:"name"`
	Path     string    `json:"path"`
	Size     int64     `json:"size"`
	Modified time.Time `json:"modified"`
	Current  bool      `json:"current"`
}

// List is the xip folder state.
type List struct {
	Files   []XIP  `json:"files"`
	Current string `json:"current"`
}

// Files lists the xips and the selected one.
func Files(cfg config.Config) (List, error) {
	directory := cfg.Path("xcode")
	entries, err := os.ReadDir(directory)
	if err != nil {
		if os.IsNotExist(err) {
			return List{}, nil
		}
		return List{}, err
	}
	selection, _ := config.LoadXcodeSelection(cfg.Home)
	list := List{}
	for _, entry := range entries {
		if entry.IsDir() || entry.Type()&os.ModeSymlink != 0 || !strings.HasSuffix(strings.ToLower(entry.Name()), ".xip") {
			continue
		}
		info, err := entry.Info()
		if err != nil || !info.Mode().IsRegular() {
			continue
		}
		list.Files = append(list.Files, XIP{
			Name: entry.Name(), Path: filepath.Join(directory, entry.Name()),
			Size: info.Size(), Modified: info.ModTime(),
		})
	}
	sort.Slice(list.Files, func(i, j int) bool {
		if list.Files[i].Modified.Equal(list.Files[j].Modified) {
			return list.Files[i].Name < list.Files[j].Name
		}
		return list.Files[i].Modified.After(list.Files[j].Modified)
	})
	for index := range list.Files {
		if list.Files[index].Name == selection.XIP {
			list.Current = list.Files[index].Path
		}
	}
	if list.Current == "" && len(list.Files) > 0 {
		list.Current = list.Files[0].Path
		_ = config.SaveXcodeSelection(cfg.Home, config.XcodeSelection{XIP: list.Files[0].Name})
	}
	markCurrent(&list)
	return list, nil
}

func markCurrent(list *List) {
	for index := range list.Files {
		list.Files[index].Current = list.Files[index].Path == list.Current
	}
}
