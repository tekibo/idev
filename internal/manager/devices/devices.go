package devices

import (
	"encoding/json"
	"io"
)

// Devices prints connected devices.
func (s *Service) Devices(stdin io.Reader, stdout, stderr io.Writer) error {
	return s.Device.List(stdin, stdout, stderr)
}

// DevicesJSON returns the parsed device list.
func (s *Service) DevicesJSON() ([]byte, error) {
	devices, err := s.ConnectedDevices()
	if err != nil {
		return nil, err
	}
	return json.MarshalIndent(devices, "", "  ")
}

// LaunchApp launches an installed app on one device.
func (s *Service) LaunchApp(bundleID, udid string, stdout, stderr io.Writer) error {
	return s.Device.LaunchApp(bundleID, udid, stdout, stderr)
}

// UninstallApp removes an installed app from one device.
func (s *Service) UninstallApp(bundleID, udid string, stdout, stderr io.Writer) error {
	return s.Device.UninstallApp(bundleID, udid, stdout, stderr)
}
