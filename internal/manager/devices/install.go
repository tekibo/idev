package devices

import (
	"github.com/tekibo/idev/internal/tools/device"
)

// InstallPlan is the ordered install sequence.
type InstallPlan = device.InstallPlan

// InstallPlan builds the install sequence.
func (s *Service) InstallPlan(bundle, bundleID, udid, devURL string) (InstallPlan, error) {
	return s.Device.InstallPlan(bundle, bundleID, udid, devURL)
}
