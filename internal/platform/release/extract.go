package release

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func extract(reader io.Reader, dest string) (int, error) {
	gzipReader, err := gzip.NewReader(reader)
	if err != nil {
		return 0, err
	}
	defer func() { _ = gzipReader.Close() }()
	archive := tar.NewReader(gzipReader)
	count := 0
	for {
		header, err := archive.Next()
		if err == io.EOF {
			return count, nil
		}
		if err != nil {
			return count, err
		}
		name := filepath.ToSlash(filepath.Clean(header.Name))
		name = strings.TrimPrefix(name, "./")
		if name == "." || name == "" {
			continue
		}
		target := filepath.Join(dest, filepath.FromSlash(name))
		if target != dest && !strings.HasPrefix(target, dest+string(os.PathSeparator)) {
			return count, fmt.Errorf("release: archive path escapes destination: %s", header.Name)
		}
		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(target, os.FileMode(header.Mode)&0777|0700); err != nil {
				return count, err
			}
		case tar.TypeReg:
			if err := os.MkdirAll(filepath.Dir(target), 0755); err != nil {
				return count, err
			}
			file, err := os.OpenFile(target, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.FileMode(header.Mode)&0777)
			if err != nil {
				return count, err
			}
			_, copyErr := io.Copy(file, archive)
			closeErr := file.Close()
			if copyErr != nil {
				return count, copyErr
			}
			if closeErr != nil {
				return count, closeErr
			}
			count++
		case tar.TypeSymlink:
			link := header.Linkname
			resolved := link
			if !filepath.IsAbs(resolved) {
				resolved = filepath.Join(filepath.Dir(target), resolved)
			}
			if !inside(dest, resolved) {
				return count, fmt.Errorf("release: symlink escapes destination: %s", header.Linkname)
			}
			if err := os.MkdirAll(filepath.Dir(target), 0755); err != nil {
				return count, err
			}
			_ = os.Remove(target)
			if err := os.Symlink(link, target); err != nil {
				return count, err
			}
			count++
		case tar.TypeLink:
			source := filepath.Join(dest, filepath.FromSlash(strings.TrimPrefix(header.Linkname, "./")))
			if !inside(dest, source) {
				return count, fmt.Errorf("release: hardlink escapes destination: %s", header.Linkname)
			}
			if err := os.Link(source, target); err != nil {
				return count, err
			}
			count++
		}
	}
}

func inside(root, path string) bool {
	root = filepath.Clean(root)
	path = filepath.Clean(path)
	return path == root || strings.HasPrefix(path, root+string(os.PathSeparator))
}
