package swift

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// Extract unpacks one official tarball and strips its top-level directory.
func Extract(archive, destination string) error {
	file, err := os.Open(archive)
	if err != nil {
		return err
	}
	defer func() { _ = file.Close() }()
	gzipReader, err := gzip.NewReader(file)
	if err != nil {
		return err
	}
	defer func() { _ = gzipReader.Close() }()
	reader := tar.NewReader(gzipReader)
	for {
		header, err := reader.Next()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		name := filepath.ToSlash(filepath.Clean(header.Name))
		parts := strings.Split(name, "/")
		if len(parts) > 1 && strings.HasPrefix(parts[0], "swift-") {
			name = strings.Join(parts[1:], "/")
		}
		if name == "." || name == "" {
			continue
		}
		target := filepath.Join(destination, filepath.FromSlash(name))
		if target != destination && !strings.HasPrefix(target, destination+string(os.PathSeparator)) {
			return fmt.Errorf("archive path escapes destination: %s", header.Name)
		}
		if err := extractEntry(reader, header, target, destination); err != nil {
			return err
		}
	}
}

func extractEntry(reader io.Reader, header *tar.Header, target, destination string) error {
	mode := os.FileMode(header.Mode) & 0777
	if mode == 0 {
		mode = 0600
	}
	switch header.Typeflag {
	case tar.TypeDir:
		return os.MkdirAll(target, mode|0700)
	case tar.TypeReg:
		if err := os.MkdirAll(filepath.Dir(target), 0755); err != nil {
			return err
		}
		file, err := os.OpenFile(target, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, mode)
		if err != nil {
			return err
		}
		_, copyErr := io.Copy(file, reader)
		closeErr := file.Close()
		if copyErr != nil {
			return copyErr
		}
		return closeErr
	case tar.TypeSymlink:
		link := header.Linkname
		resolved := link
		if !filepath.IsAbs(resolved) {
			resolved = filepath.Join(filepath.Dir(target), resolved)
		}
		if !inside(destination, resolved) {
			return fmt.Errorf("archive: symlink escapes destination: %s", header.Linkname)
		}
		if err := os.MkdirAll(filepath.Dir(target), 0755); err != nil {
			return err
		}
		_ = os.Remove(target)
		return os.Symlink(link, target)
	case tar.TypeLink:
		source := filepath.Join(destination, filepath.FromSlash(header.Linkname))
		if !inside(destination, source) {
			return fmt.Errorf("archive: hardlink escapes destination: %s", header.Linkname)
		}
		return os.Link(source, target)
	default:
		return nil
	}
}

func inside(root, path string) bool {
	root = filepath.Clean(root)
	path = filepath.Clean(path)
	return path == root || strings.HasPrefix(path, root+string(os.PathSeparator))
}
