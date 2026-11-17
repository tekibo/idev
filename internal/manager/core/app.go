// Package core owns the shared manager services and the typed error contract.
package core

import (
	"bytes"
	"errors"
	"io"
	"path/filepath"

	"github.com/tekibo/idev/internal/platform/config"
	"github.com/tekibo/idev/internal/platform/execx"
	"github.com/tekibo/idev/internal/platform/manifest"
	"github.com/tekibo/idev/internal/tools/device"
	"github.com/tekibo/idev/internal/tools/xtool"
)

// Services holds the shared services used by the CLI.
type Services struct {
	Config config.Config
	Runner execx.ContextRunner
	Xtool  *xtool.Client
	Device *device.Client
}

// New builds the shared services.
func New() *Services {
	cfg := config.Load()
	runner := execx.Exec{}
	services := &Services{
		Config: cfg,
		Runner: runner,
		Xtool:  xtool.NewWithRunner(cfg, runner),
	}
	services.ResolveSDK()
	services.Device = device.NewWithRunner(services.Config, runner)
	return services
}

// ResolveSDK asks xtool for the registered Darwin SDK path.
func (s *Services) ResolveSDK() {
	if root := s.Config.SDKRoot(); root != "" {
		s.Config.SDKRootPath = root
		return
	}
	if !s.Xtool.Available() {
		return
	}
	var out bytes.Buffer
	if err := s.Xtool.SDKStatus(&out, io.Discard); err != nil {
		return
	}
	root := xtool.ParseSDKPath(out.String())
	if root != "" && config.IsFile(filepath.Join(root, "info.json")) {
		s.Config.SDKRootPath = root
	}
}

// Version returns the manager version.
func (s *Services) Version() string {
	return manifest.Version
}

// ExitCode maps lower-layer errors to the CLI exit code.
func (s *Services) ExitCode(err error) int {
	if err == nil {
		return 0
	}
	if errors.Is(err, device.ErrNoDeviceTool) || errors.Is(err, xtool.ErrMissing) {
		return 127
	}
	return execx.ExitCode(err)
}
