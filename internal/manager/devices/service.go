// Package devices owns the device list and app install operations.
package devices

import (
	"github.com/tekibo/idev/internal/manager/core"
)

// Service owns the device operations.
type Service struct {
	*core.Services
}

// New builds the device service.
func New(services *core.Services) *Service {
	return &Service{Services: services}
}
