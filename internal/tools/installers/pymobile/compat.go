package pymobile

import "os"

// ensureCompat exposes the private venv at the upstream apple-dev path.
// A foreign venv is left untouched.
func ensureCompat(spec Spec) error {
	info, err := os.Lstat(spec.Compat)
	if os.IsNotExist(err) {
		return os.Symlink(spec.Root, spec.Compat)
	}
	if err != nil {
		return err
	}
	if info.Mode()&os.ModeSymlink == 0 {
		return nil
	}
	return nil
}
