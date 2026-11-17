package manager

import (
	"fmt"
	"strings"

	"github.com/tekibo/idev/internal/platform/config"
)

// UpstreamCommand is one script or helper invocation.
type UpstreamCommand struct {
	Name string
	Args []string
	Env  []string
}

// KMPCommand returns the Kotlin Multiplatform framework command.
func (a *App) KMPCommand(args []string) (UpstreamCommand, error) {
	return a.scriptCommand(a.Config.Path("kmp", "build-framework.sh"), args,
		"run idev setup --flow kmp")
}

// HostManifestsCommand returns the SwiftPM host-manifest command.
func (a *App) HostManifestsCommand(args []string) (UpstreamCommand, error) {
	return a.pythonTool("host-manifests.py", args)
}

// XcodeProjectCommand returns the Xcode project converter.
func (a *App) XcodeProjectCommand(args []string) (UpstreamCommand, error) {
	return a.pythonTool("xcodeproj2xtool.py", args)
}

// AscCommand returns the App Store Connect helper.
func (a *App) AscCommand(args []string) (UpstreamCommand, error) {
	return a.signedPythonTool("asc.py", args)
}

// ProvisionDevCommand returns the development profile helper.
func (a *App) ProvisionDevCommand(args []string) (UpstreamCommand, error) {
	return a.signedPythonTool("provision-dev.py", args)
}

// SignDevCommand returns the development signing helper.
func (a *App) SignDevCommand(args []string) (UpstreamCommand, error) {
	return a.scriptCommand(a.Config.Path("apple-tools", "tools", "sign-dev.sh"), args,
		"run idev setup --apple-tools")
}

// SDKFreeCommand returns one sdk-free script command.
func (a *App) SDKFreeCommand(script string, args []string) (UpstreamCommand, error) {
	relative, ok := sdkFreeScripts[script]
	if !ok {
		return UpstreamCommand{}, fmt.Errorf("unknown sdk-free script %q", script)
	}
	command, err := a.scriptCommand(a.Config.Path("apple-tools", "sdk-free", relative), args,
		"run idev setup --apple-tools")
	if err != nil {
		return UpstreamCommand{}, err
	}
	command.Env = a.Config.SwiftEnvWith(map[string]string{
		"SDKFREE_HOME": a.Config.Path("sdk-free"),
	})
	return command, nil
}

var sdkFreeScripts = map[string]string{
	"setup":              "setup.sh",
	"flutter-build":      "flutter-build.sh",
	"ship":               "ship.sh",
	"cc":                 "cc.sh",
	"swiftc":             "swiftc.sh",
	"swift-build-stdlib": "swift/build-stdlib.sh",
	"swift-overlays":     "swift/overlays.sh",
	"swiftui-build":      "swiftui/swiftui-build.sh",
}

// SimViewCommand returns one simview script command.
func (a *App) SimViewCommand(script string, args []string) (UpstreamCommand, error) {
	if !safeName(script) {
		return UpstreamCommand{}, fmt.Errorf("invalid simview script %q", script)
	}
	return a.scriptCommand(a.Config.Path("apple-tools", "simview", script+".sh"), args,
		"run idev setup --apple-tools")
}

// MachoLoadCommand returns the machoload build command.
func (a *App) MachoLoadCommand(args []string) (UpstreamCommand, error) {
	return a.scriptCommand(a.Config.Path("apple-tools", "machoload", "build.sh"), args,
		"run idev setup --apple-tools")
}

// CompileAppIconCommand returns the app icon compiler command.
func (a *App) CompileAppIconCommand(args []string) (UpstreamCommand, error) {
	return a.scriptCommand(a.Config.Path("apple-tools", "tools", "compile-app-icon.sh"), args,
		"run idev setup --apple-tools")
}

func (a *App) scriptCommand(path string, args []string, hint string) (UpstreamCommand, error) {
	if !config.IsExec(path) {
		return UpstreamCommand{}, fmt.Errorf("%s not found; %s", path, hint)
	}
	return UpstreamCommand{Name: path, Args: args, Env: a.Config.SwiftEnv()}, nil
}

func safeName(name string) bool {
	return name != "" && name != "." && name != ".." && !strings.ContainsAny(name, `/\`)
}
