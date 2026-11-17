package project

import (
	"path/filepath"
	"sort"
)

// DetectWorkspace finds every supported app under one directory.
func DetectWorkspace(root string) Workspace {
	if abs, err := filepath.Abs(root); err == nil {
		root = abs
	}
	workspace := Workspace{Dir: root, Name: filepath.Base(root)}
	if info, ok := detectNonNode(root); ok {
		info.WorkspaceDir = root
		workspace.Apps = []Info{info}
		return workspace
	}
	seen := map[string]bool{}
	for _, dir := range workspaceCandidates(root) {
		info, ok := inspectNode(dir, root)
		if !ok {
			continue
		}
		key := filepath.Clean(info.Dir)
		if seen[key] {
			continue
		}
		seen[key] = true
		workspace.Apps = append(workspace.Apps, info)
	}
	sort.Slice(workspace.Apps, func(i, j int) bool {
		if workspace.Apps[i].Name == workspace.Apps[j].Name {
			return workspace.Apps[i].Dir < workspace.Apps[j].Dir
		}
		return workspace.Apps[i].Name < workspace.Apps[j].Name
	})
	return workspace
}

// App resolves one app by name, directory base, or path.
func (w Workspace) App(name string) (Info, bool) {
	if name == "" {
		return Info{}, false
	}
	abs, _ := filepath.Abs(name)
	for _, app := range w.Apps {
		if app.Name == name || filepath.Base(app.Dir) == name || filepath.Clean(app.Dir) == filepath.Clean(abs) {
			return app, true
		}
	}
	return Info{}, false
}
