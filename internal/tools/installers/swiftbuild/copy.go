package swiftbuild

import (
	"os"
	"path/filepath"

	"github.com/tekibo/idev/internal/tools/installers/installutil"
)

func copyBundles(sourceRoot, bundles string) error {
	matches, err := filepath.Glob(filepath.Join(sourceRoot, ".build", "release", "*.bundle"))
	if err != nil {
		return err
	}
	for _, bundle := range matches {
		if err := installutil.CopyDir(bundle, filepath.Join(bundles, filepath.Base(bundle))); err != nil {
			return err
		}
	}
	return nil
}

func removeBundles(bundles string) {
	matches, _ := filepath.Glob(filepath.Join(bundles, "*.bundle"))
	for _, bundle := range matches {
		_ = os.RemoveAll(bundle)
	}
}
