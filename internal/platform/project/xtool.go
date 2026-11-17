package project

import (
	"path/filepath"
)

func inspectXtool(dir string) Info {
	path := filepath.Join(dir, "xtool.yml")
	cfg := readXtool(path)
	name := cfg.Name
	if name == "" {
		name = filepath.Base(dir)
	}
	product := cfg.Product
	if product == "" {
		product = name
	}
	info := Info{
		Kind:             SwiftUI,
		Dir:              dir,
		Name:             name,
		Label:            LabelFor(SwiftUI),
		BundleID:         cfg.BundleID,
		BundleIDEditable: true,
		Product:          product,
		DeploymentTarget: cfg.DeploymentTarget,
		HasExtensions:    cfg.HasExtensions,
		XtoolFile:        path,
	}
	info.AppEntry = hasAppEntry(dir)
	info.PackageFile = filepath.Join(dir, "Package.swift")
	if !fileExists(info.PackageFile) {
		addWarning(&info, "Package.swift is missing.")
	} else {
		if info.DeploymentTarget == "" {
			info.DeploymentTarget = packageIOS(info.PackageFile)
		}
		if cfg.Product == "" {
			if product := packageProduct(info.PackageFile); product != "" {
				info.Product = product
			}
		}
	}
	if info.BundleID == "" {
		addWarning(&info, "Bundle ID is not declared in xtool.yml.")
	}
	if !info.AppEntry {
		addWarning(&info, "No @main SwiftUI App was found in Sources.")
	}
	if cfg.IconPath != "" {
		info.IconPath = absolute(dir, cfg.IconPath)
		if !pathExists(info.IconPath) {
			info.IconPath = ""
			addWarning(&info, "Configured app icon was not found.")
		}
	}
	if info.IconPath == "" {
		info.IconPath = findIcon(dir)
	}
	if info.IconPath == "" && cfg.IconPath == "" {
		addWarning(&info, "No app icon was found.")
	}
	if cfg.InfoPath != "" {
		info.InfoPath = absolute(dir, cfg.InfoPath)
		if !fileExists(info.InfoPath) {
			info.InfoPath = ""
			addWarning(&info, "Configured Info.plist was not found.")
		}
	} else {
		info.InfoPath = findInfoPlist(dir, name)
	}
	if info.InfoPath == "" && cfg.InfoPath == "" {
		addWarning(&info, "No source Info.plist; version and display name are generated.")
	}
	if cfg.EntitlementsPath != "" {
		info.EntitlementsFile = absolute(dir, cfg.EntitlementsPath)
		if !fileExists(info.EntitlementsFile) {
			info.EntitlementsFile = ""
			addWarning(&info, "Configured entitlements file was not found.")
		}
	} else {
		def := filepath.Join(dir, "xtool", "ship-signing", "entitlements.plist")
		if fileExists(def) {
			info.EntitlementsFile = def
		}
	}
	return info
}
