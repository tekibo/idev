package darwin

import (
	"context"
	"io"
	"os"
	"path/filepath"

	"github.com/tekibo/idev/internal/platform/config"
	"github.com/tekibo/idev/internal/platform/execx"
)

// Result reports one Darwin tools and OAM install.
type Result struct {
	Payload     string
	SDKRoot     string
	ToolsBin    string
	MacToolsBin string
	OAMServer   string
	Plugins     string
	Already     bool
}

// Install builds the vendored tools and copies them into the xtool SDK.
func Install(ctx context.Context, cfg config.Config, runner execx.Runner, force bool, stdout, stderr io.Writer) (Result, error) {
	if err := ctx.Err(); err != nil {
		return Result{}, err
	}
	spec := NewSpec(cfg)
	if err := validate(cfg, spec); err != nil {
		return Result{}, err
	}
	if !force && config.IsExec(filepath.Join(spec.ToolsBin, "actool")) && config.IsExec(spec.OAMServer) {
		if err := patchSDK(cfg, runner, spec, stdout, stderr); err != nil {
			return Result{}, err
		}
		return Result{
			Payload: spec.Payload, SDKRoot: spec.SDKRoot, ToolsBin: spec.ToolsBin,
			MacToolsBin: spec.MacToolsBin, OAMServer: spec.OAMServer, Plugins: spec.Plugins,
			Already: true,
		}, nil
	}
	actool, err := buildActool(cfg, runner, spec.Payload, force, stdout, stderr)
	if err != nil {
		return Result{}, err
	}
	oam, err := buildOAM(cfg, runner, spec.Payload, force, stdout, stderr)
	if err != nil {
		return Result{}, err
	}
	if err := installTools(spec, actool); err != nil {
		return Result{}, err
	}
	if err := installOAM(spec, oam); err != nil {
		return Result{}, err
	}
	if err := ctx.Err(); err != nil {
		return Result{}, err
	}
	if err := patchSDK(cfg, runner, spec, stdout, stderr); err != nil {
		return Result{}, err
	}
	_ = os.RemoveAll(filepath.Join(spec.Payload, "darwin-tools", ".build"))
	_ = os.RemoveAll(filepath.Join(spec.Payload, "oam", ".build"))
	return Result{
		Payload:     spec.Payload,
		SDKRoot:     spec.SDKRoot,
		ToolsBin:    spec.ToolsBin,
		MacToolsBin: spec.MacToolsBin,
		OAMServer:   spec.OAMServer,
		Plugins:     spec.Plugins,
	}, nil
}

func validate(cfg config.Config, spec Spec) error {
	required := []string{
		"darwin-tools/Package.swift",
		"oam/Package.swift",
		"tools/ibtool",
		"tools/xcstringstool",
		"tools/momc",
		"tools/nibarchive.py",
		"tools/keyorder.py",
		"tools/xcstrings_symbols.py",
		"tools/ibmac.py",
	}
	for _, relative := range required {
		if !config.IsFile(filepath.Join(spec.Payload, relative)) {
			return ErrPayloadMissing
		}
	}
	if spec.SDKRoot == "" || spec.SDKPath == "" {
		return ErrSDKMissing
	}
	if !cfg.DarwinSDKReady() {
		return ErrSDKMissing
	}
	if !config.IsExec(filepath.Join(cfg.SwiftBin(), "swift")) {
		return ErrSwiftMissing
	}
	if _, ok := cfg.FindBinary("python3"); !ok {
		return ErrPythonMissing
	}
	return nil
}
