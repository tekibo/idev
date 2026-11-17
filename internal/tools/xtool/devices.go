package xtool

import (
	"strings"
)

// Device is one connected device.
type Device struct {
	Name       string `json:"name"`
	Connection string `json:"connection"`
	UDID       string `json:"udid"`
}

// ParseDevices parses xtool devices output: Name [connection]: UDID.
func ParseDevices(output []byte) []Device {
	devices := []Device{}
	for _, line := range strings.Split(strings.TrimSpace(string(output)), "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "Waiting for devices") {
			continue
		}
		left, udid, ok := strings.Cut(line, ": ")
		if !ok {
			continue
		}
		name, connection, ok := strings.Cut(left, " [")
		if !ok {
			continue
		}
		devices = append(devices, Device{
			Name:       name,
			Connection: strings.TrimSuffix(connection, "]"),
			UDID:       udid,
		})
	}
	return devices
}
