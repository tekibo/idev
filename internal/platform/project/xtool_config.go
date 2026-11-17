package project

import (
	"os"
	"strings"
)

type xtoolConfig struct {
	Name             string
	BundleID         string
	Product          string
	DeploymentTarget string
	InfoPath         string
	IconPath         string
	EntitlementsPath string
	HasExtensions    bool
}

func readXtool(path string) xtoolConfig {
	var cfg xtoolConfig
	data, err := os.ReadFile(path)
	if err != nil {
		return cfg
	}
	for _, raw := range strings.Split(string(data), "\n") {
		line := strings.TrimSpace(raw)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		if strings.HasPrefix(line, "extensions:") {
			cfg.HasExtensions = true
			continue
		}
		if strings.HasPrefix(raw, " ") || strings.HasPrefix(raw, "\t") {
			continue
		}
		key, value, ok := strings.Cut(line, ":")
		if !ok {
			continue
		}
		value = scalar(value)
		switch strings.TrimSpace(key) {
		case "name":
			cfg.Name = value
		case "bundleID":
			cfg.BundleID = value
		case "product":
			cfg.Product = value
		case "deploymentTarget":
			cfg.DeploymentTarget = value
		case "infoPath":
			cfg.InfoPath = value
		case "iconPath":
			cfg.IconPath = value
		case "entitlementsPath":
			cfg.EntitlementsPath = value
		}
	}
	return cfg
}

func scalar(value string) string {
	value = strings.TrimSpace(value)
	if len(value) >= 2 {
		if value[0] == '"' && value[len(value)-1] == '"' {
			return value[1 : len(value)-1]
		}
		if value[0] == '\'' && value[len(value)-1] == '\'' {
			return value[1 : len(value)-1]
		}
	}
	if at := strings.Index(value, " #"); at >= 0 {
		value = strings.TrimSpace(value[:at])
	}
	return value
}
