package uninstall

import (
	"os"
	"path/filepath"
)

// ownedLink reports whether path is a symlink into the private home.
func ownedLink(path, home string) bool {
	info, err := os.Lstat(path)
	if err != nil || info.Mode()&os.ModeSymlink == 0 {
		return false
	}
	target, err := os.Readlink(path)
	if err != nil {
		return false
	}
	if !filepath.IsAbs(target) {
		target = filepath.Join(filepath.Dir(path), target)
	}
	root, err := filepath.Abs(home)
	if err != nil {
		return false
	}
	target, err = filepath.Abs(target)
	if err != nil {
		return false
	}
	return isUnder(root, target)
}

// ownedExecutable reports whether path is the running idev binary.
func ownedExecutable(path string) bool {
	info, err := os.Stat(path)
	if err != nil || !info.Mode().IsRegular() {
		return false
	}
	self, err := os.Executable()
	if err != nil {
		return false
	}
	selfInfo, err := os.Stat(self)
	if err != nil {
		return false
	}
	return os.SameFile(info, selfInfo)
}
