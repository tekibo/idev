package darwin

import (
	"os"
	"path/filepath"

	"github.com/tekibo/idev/internal/tools/installers/installutil"
)

func installTools(spec Spec, actool string) error {
	for _, bin := range []string{spec.ToolsBin, spec.MacToolsBin} {
		if err := os.MkdirAll(bin, 0755); err != nil {
			return err
		}
		if err := installutil.CopyFile(actool, filepath.Join(bin, "actool"), 0755); err != nil {
			return err
		}
		for _, name := range []string{"xcstringstool", "ibtool", "momc"} {
			if err := installutil.CopyFile(filepath.Join(spec.Payload, "tools", name), filepath.Join(bin, name), 0755); err != nil {
				return err
			}
		}
		for _, name := range []string{"nibarchive.py", "keyorder.py", "xcstrings_symbols.py"} {
			if err := installutil.CopyFile(filepath.Join(spec.Payload, "tools", name), filepath.Join(bin, name), 0644); err != nil {
				return err
			}
		}
	}
	return installutil.CopyFile(filepath.Join(spec.Payload, "tools", "ibmac.py"), filepath.Join(spec.MacToolsBin, "ibmac.py"), 0644)
}

func installOAM(spec Spec, server string) error {
	if err := os.MkdirAll(filepath.Dir(spec.OAMServer), 0755); err != nil {
		return err
	}
	if err := installutil.CopyFile(server, spec.OAMServer, 0755); err != nil {
		return err
	}
	if err := os.MkdirAll(spec.Plugins, 0755); err != nil {
		return err
	}
	for _, name := range []string{"libSwiftDataMacros.so", "libFoundationMacros.so"} {
		if err := os.WriteFile(filepath.Join(spec.Plugins, name), nil, 0644); err != nil {
			return err
		}
	}
	return nil
}
