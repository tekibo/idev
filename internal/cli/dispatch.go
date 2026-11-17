package cli

import (
	"fmt"
	"os"

	helptree "github.com/tekibo/idev/internal/cli/help"
	"github.com/tekibo/idev/internal/manager"
)

// dispatch runs one parsed command.
func dispatch(args []string) int {
	switch args[0] {
	case "version", "-v", "--version":
		fmt.Printf("idev %s\n", manager.New().Version())
		return 0
	case "setup":
		return Setup(args[1:])
	case "next":
		return Next(args[1:])
	case "agents":
		return Agents(args[1:])
	case "doctor":
		return Doctor(args[1:])
	case "system":
		return System(args[1:])
	case "devices":
		return Devices(args[1:])
	case "auth":
		return Auth(args[1:])
	case "xtool":
		return Xtool(args[1:])
	case "swift":
		return Swift(args[1:])
	case "ds":
		return DS(args[1:])
	case "sdk":
		return SDK(args[1:])
	case "xcode":
		return Xcode(args[1:])
	case "new":
		return NewProject(args[1:])
	case "config":
		return ConfigCommand(args[1:])
	case "expo":
		return Expo(args[1:])
	case "rn", "react-native":
		return ReactNative(args[1:])
	case "pod":
		return Pod(args[1:])
	case "ios":
		return Expo(append([]string{"run:ios"}, args[1:]...))
	case "start":
		return Expo(append([]string{"start"}, args[1:]...))
	case "install-app":
		return InstallApp(args[1:])
	case "launch-app":
		return LaunchApp(args[1:])
	case "uninstall-app":
		return UninstallApp(args[1:])
	case "ipad-screen":
		return IpadScreen(args[1:])
	case "device-run":
		return DeviceRun(args[1:])
	case "ship":
		return Ship(args[1:])
	case "lsp":
		return LSP(args[1:])
	case "update":
		return Update(args[1:])
	case "flutter":
		return Flutter(args[1:])
	case "kmp":
		return KMP(args[1:])
	case "host-manifests":
		return HostManifests(args[1:])
	case "simview":
		return SimView(args[1:])
	case "machoload":
		return MachoLoad(args[1:])
	case "compile-app-icon":
		return CompileAppIcon(args[1:])
	case "xcodeproj2xtool":
		return XcodeProject(args[1:])
	case "asc":
		return Asc(args[1:])
	case "provision-dev":
		return ProvisionDev(args[1:])
	case "sign-dev":
		return SignDev(args[1:])
	case "sdk-free":
		return SDKFree(args[1:])
	case "project":
		return Project(args[1:])
	case "uninstall":
		return Uninstall(args[1:])
	default:
		fmt.Fprintf(os.Stderr, "idev: unknown command %q\n\n", args[0])
		helptree.Show(nil, JSON)
		return 2
	}
}
