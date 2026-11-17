package project

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type packageJSON struct {
	Name            string            `json:"name"`
	Dependencies    map[string]string `json:"dependencies"`
	DevDependencies map[string]string `json:"devDependencies"`
	Workspaces      json.RawMessage   `json:"workspaces"`
}

func inspectNode(dir, workspaceDir string) (Info, bool) {
	pkg, ok := readPackageJSON(dir)
	if !ok {
		return Info{}, false
	}
	kind, expoVersion, rnVersion := nodeKind(dir, pkg)
	if kind == Unknown {
		return Info{}, false
	}
	info := Info{
		Kind:         kind,
		Dir:          dir,
		WorkspaceDir: workspaceDir,
		Name:         nodeName(dir, pkg),
		Label:        LabelFor(kind),
		ExpoVersion:  expoVersion,
		RNVersion:    rnVersion,
		DevClient:    dependency(pkg, "expo-dev-client") != "",
	}
	info.NativeDir, info.ProjectPath, info.Scheme = nodeNative(dir)
	info.Workflow = nodeWorkflow(kind, info.NativeDir)
	info.BundleID, info.BundleIDSource = nodeBundleID(dir, info.ProjectPath)
	info.PackageManager = packageManager(workspaceDir, dir)
	if info.BundleID == "" {
		addWarning(&info, "No iOS bundle ID was found in the app config or native project.")
	}
	return info, true
}

func readPackageJSON(dir string) (packageJSON, bool) {
	data, err := os.ReadFile(filepath.Join(dir, "package.json"))
	if err != nil {
		return packageJSON{}, false
	}
	var pkg packageJSON
	if json.Unmarshal(data, &pkg) != nil {
		return packageJSON{}, false
	}
	return pkg, true
}

func nodeKind(dir string, pkg packageJSON) (Kind, string, string) {
	expoVersion := pkg.Dependencies["expo"]
	rnVersion := pkg.Dependencies["react-native"]
	if expoVersion != "" {
		return Expo, expoVersion, rnVersion
	}
	if rnVersion != "" {
		return ReactNative, "", rnVersion
	}
	// Libraries often keep react-native in devDependencies. Require an app
	// project before treating that as a bare RN app.
	if dev := pkg.DevDependencies["react-native"]; dev != "" && firstGlob(filepath.Join(dir, "ios", "*.xcodeproj")) != "" {
		return ReactNative, "", dev
	}
	return Unknown, "", ""
}

func dependency(pkg packageJSON, name string) string {
	if value := pkg.Dependencies[name]; value != "" {
		return value
	}
	return pkg.DevDependencies[name]
}

func nodeName(dir string, pkg packageJSON) string {
	if config, ok := staticAppConfig(dir); ok && config.Name != "" {
		return config.Name
	}
	if pkg.Name != "" {
		return pkg.Name
	}
	return filepath.Base(dir)
}

func nodeWorkflow(kind Kind, nativeDir string) string {
	if kind == Expo {
		if nativeDir != "" {
			return WorkflowExpoPrebuild
		}
		return WorkflowExpoManaged
	}
	if nativeDir != "" {
		return WorkflowBareRN
	}
	return WorkflowBareRNMissing
}
