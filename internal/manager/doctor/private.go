package doctor

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/tekibo/idev/internal/platform/config"
	"github.com/tekibo/idev/internal/platform/system"
)

func privateChecks(cfg config.Config, projectDir string) []Check {
	sdkInfo := ""
	actool := ""
	if cfg.DarwinSDKReady() {
		sdkInfo = filepath.Join(cfg.SDKRoot(), "info.json")
		if sdk := cfg.SDKPath(); sdk != "" {
			actool = filepath.Join(sdk, "Platforms", "iPhoneOS.platform", "Developer", "usr", "bin", "actool")
		}
	}
	source := Check{
		ID:        "private.source_config",
		Class:     "private",
		Status:    status(cfg.SourcesErr == nil),
		FixableBy: system.FixIdev,
		Hint:      "run idev setup",
		Critical:  true,
	}
	if cfg.SourcesErr != nil {
		source.Detail = cfg.SourcesErr.Error()
	}
	checks := []Check{
		source,
		sourcePinCheck(cfg),
		flowCheck(cfg),
		privateTool("private.swift", cfg, "swift", "install the private Swift toolchain"),
		privateTool("private.xtool", cfg, "xtool", "install the private xtool"),
		privateTool("private.pymobiledevice3", cfg, "pymobiledevice3", "install the private pymobiledevice3"),
		privateTool("private.pod", cfg, "pod", "install the private CocoaPods"),
		privateTool("private.rcodesign", cfg, "rcodesign", "install rcodesign"),
		privateTool("private.ipsw", cfg, "ipsw", "install ipsw"),
		privateFile("private.swift_build", cfg.Path("bin", "expo-build-bin"), "install expo-build-bin"),
		privateFile("private.swbbuildservice", cfg.Path("bin", "SWBBuildService"), "install SWBBuildService"),
		privateFile("private.sdk", sdkInfo, "run xtool setup with Xcode.xip"),
		privateOptionalFile("private.apple_tools", actool, "run idev setup --apple-tools"),
		privateOptionalFile("private.flutter_arm", cfg.Path("flutter", "setup.sh"), "run idev setup"),
		privateOptionalFile("private.kmp", cfg.Path("kmp", "build-framework.sh"), "run idev setup --flow kmp"),
		privateOptionalFile("private.reactnative", cfg.Path("react-native", "tools", "node-darwin.js"), "run idev setup"),
		tailscaleCheck(cfg),
		ipadScreenCheck(cfg),
	}
	return append(checks, projectChecks(projectDir)...)
}

func ipadScreenCheck(cfg config.Config) Check {
	tool := cfg.Path("apple-tools", "tools", "ipad-screen.sh")
	check := Check{
		ID:        "private.ipad_screen",
		Class:     "private",
		FixableBy: system.FixUser,
		Hint:      "run idev ipad-screen install",
	}
	if !config.IsExec(tool) {
		check.Status = system.StatusMissing
		check.FixableBy = system.FixIdev
		check.Hint = "run idev setup --apple-tools"
		return check
	}
	missing := []string{}
	if !config.IsExec(cfg.PymobileBin()) {
		missing = append(missing, "the private pymobiledevice3 venv is missing")
		check.Hint = "run idev setup --pymobile"
	}
	check.Status = system.StatusOK
	if len(missing) > 0 {
		check.Status = system.StatusWarn
		check.Detail = strings.Join(missing, "\n")
	}
	return check
}

func sourcePinCheck(cfg config.Config) Check {
	warnings := config.SourceWarnings(cfg.Sources)
	check := Check{
		ID:        "private.source_pins",
		Class:     "private",
		Status:    status(len(warnings) == 0),
		FixableBy: system.FixIdev,
		Hint:      "run idev config migrate",
	}
	lines := make([]string, 0, len(warnings))
	for _, warning := range warnings {
		lines = append(lines, fmt.Sprintf("%s: %s, required %s (%s)", warning.Field, warning.Current, warning.Required, warning.Reason))
		if warning.Block {
			check.Critical = true
		}
	}
	check.Detail = strings.Join(lines, "\n")
	return check
}
