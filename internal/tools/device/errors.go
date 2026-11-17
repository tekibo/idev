package device

import "errors"

// ErrNoDeviceTool reports that no device transport is installed.
var ErrNoDeviceTool = errors.New("neither xtool nor pymobiledevice3 is available")

// ErrNoLaunchTool reports that pymobiledevice3 is missing.
var ErrNoLaunchTool = errors.New("pymobiledevice3 not found")

// ErrNoDeviceRun reports that the device-run script is missing.
var ErrNoDeviceRun = errors.New("device-run.sh not found; run idev setup")
