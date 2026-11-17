// Package project detects supported project types.
package project

// Kind is one supported project type.
type Kind string

// Project kinds.
const (
	SwiftUI      Kind = "xtool-swiftui"
	Flutter      Kind = "flutter"
	Expo         Kind = "expo"
	ReactNative  Kind = "react-native"
	Xcode        Kind = "xcode"
	SwiftPackage Kind = "swift-package"
	Unknown      Kind = "unknown"
)

// Node workflows.
const (
	WorkflowExpoManaged   = "expo-managed"
	WorkflowExpoPrebuild  = "expo-prebuild"
	WorkflowBareRN        = "bare-rn"
	WorkflowBareRNMissing = "bare-rn-missing"
)

// Info is one detected project.
type Info struct {
	Kind             Kind     `json:"kind"`
	Dir              string   `json:"dir"`
	WorkspaceDir     string   `json:"workspace_dir,omitempty"`
	Name             string   `json:"name"`
	Label            string   `json:"label"`
	BundleID         string   `json:"bundle_id"`
	BundleIDEditable bool     `json:"bundle_id_editable"`
	BundleIDSource   string   `json:"bundle_id_source,omitempty"`
	Product          string   `json:"product"`
	DeploymentTarget string   `json:"deployment_target"`
	ExpoVersion      string   `json:"expo_version"`
	RNVersion        string   `json:"rn_version,omitempty"`
	Workflow         string   `json:"workflow,omitempty"`
	DevClient        bool     `json:"dev_client,omitempty"`
	NativeDir        string   `json:"native_dir,omitempty"`
	ProjectPath      string   `json:"project_path,omitempty"`
	Scheme           string   `json:"scheme,omitempty"`
	PackageManager   string   `json:"package_manager,omitempty"`
	AppEntry         bool     `json:"app_entry"`
	HasExtensions    bool     `json:"has_extensions"`
	XtoolFile        string   `json:"xtool_file,omitempty"`
	PackageFile      string   `json:"package_file,omitempty"`
	InfoPath         string   `json:"info_path,omitempty"`
	IconPath         string   `json:"icon_path,omitempty"`
	EntitlementsFile string   `json:"entitlements_file,omitempty"`
	Warnings         []string `json:"warnings,omitempty"`
}

// Workspace is one opened directory and its detected apps.
type Workspace struct {
	Dir      string   `json:"dir"`
	Name     string   `json:"name"`
	Apps     []Info   `json:"apps"`
	Warnings []string `json:"warnings,omitempty"`
}

// Known reports whether the project kind is supported.
func (i Info) Known() bool {
	return i.Kind != "" && i.Kind != Unknown
}

// LabelFor names one project kind.
func LabelFor(kind Kind) string {
	switch kind {
	case SwiftUI:
		return "SwiftUI"
	case Flutter:
		return "Flutter"
	case Expo:
		return "Expo"
	case ReactNative:
		return "React Native"
	case Xcode:
		return "Xcode"
	case SwiftPackage:
		return "Swift package"
	default:
		return "Unknown"
	}
}
