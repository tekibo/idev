package core

import "github.com/tekibo/idev/internal/tools/xtool"

// ConnectedDevices lists the physical devices xtool can see.
func (s *Services) ConnectedDevices() ([]xtool.Device, error) {
	output, err := s.Xtool.Output([]string{"devices", "--no-wait"})
	if err != nil {
		return nil, err
	}
	return xtool.ParseDevices(output), nil
}
