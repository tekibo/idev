package cli

import (
	"fmt"
	"strings"

	"github.com/tekibo/idev/internal/manager"
)

func printDoctor(report manager.DoctorReport) {
	fmt.Printf("idev %s\n", report.Version)
	fmt.Printf("home:    %s\n", report.Home)
	fmt.Printf("manager: %s\n", report.Manager)
	fmt.Printf("distro:  %s\n", report.Distro)
	lastClass := ""
	for _, check := range report.Checks {
		if check.Class != lastClass {
			fmt.Printf("\n%s:\n", classLabel(check.Class))
			lastClass = check.Class
		}
		status := doctorStatus(check)
		fmt.Printf("  %-4s %s\n", status, doctorLabel(check.ID))
		if status == "ok" {
			continue
		}
		detail := check.Detail
		if detail == "" {
			detail = check.Hint
		}
		if detail != "" {
			fmt.Printf("       %s\n", detail)
		}
	}
	if report.SystemPlan.Command != "" {
		fmt.Printf("\nsystem command:\n  %s\n", report.SystemPlan.Command)
	}
	if report.SystemPlan.Blocked != "" {
		fmt.Printf("blocked: %s\n", report.SystemPlan.Blocked)
	}
	if next := doctorNext(report); next != "" {
		fmt.Printf("\nnext: %s\n", next)
	}
}

func classLabel(class string) string {
	if class == "" {
		return class
	}
	return strings.ToUpper(class[:1]) + class[1:]
}

func doctorStatus(check manager.DoctorCheck) string {
	if check.OK() {
		return "ok"
	}
	if check.Status == "warn" {
		return "WARN"
	}
	if check.Critical {
		return "FAIL"
	}
	return "WARN"
}

func doctorLabel(id string) string {
	labels := map[string]string{
		"private.source_config":   "Source config",
		"private.source_pins":     "Source pins",
		"private.flows":           "Flows",
		"private.swift":           "Swift",
		"private.xtool":           "xtool",
		"private.pymobiledevice3": "pymobiledevice3",
		"private.pod":             "CocoaPods",
		"private.rcodesign":       "rcodesign",
		"private.ipsw":            "ipsw",
		"private.apple_tools":     "Darwin tools and OAM",
		"private.swift_build":     "Swift Build",
		"private.swbbuildservice": "SWBBuildService",
		"private.sdk":             "Darwin SDK",
		"private.flutter_arm":     "Flutter",
		"private.kmp":             "Kotlin Multiplatform",
		"private.reactnative":     "React Native & Expo",
		"private.ipad_screen":     "iPad screen viewer",
	}
	if label, ok := labels[id]; ok {
		return label
	}
	for _, prefix := range []string{"system.", "project.", "device.", "auth."} {
		if strings.HasPrefix(id, prefix) {
			return strings.ReplaceAll(strings.TrimPrefix(id, prefix), "_", "-")
		}
	}
	return id
}

func doctorNext(report manager.DoctorReport) string {
	if len(report.SystemPlan.Packages) > 0 {
		return "idev system apply"
	}
	for _, check := range report.Checks {
		if check.OK() {
			continue
		}
		if check.Hint != "" {
			return check.Hint
		}
	}
	return "idev next"
}
