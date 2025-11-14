package utils

import (
	"archive/tar"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func UnpackTar(tarballReader io.Reader, dest string) error {
	tarr := tar.NewReader(tarballReader)

	for {
		hdr, err := tarr.Next()
		//fmt.Printf("new file %s\n", hdr.Name)
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("reading tar header: %w", err)
		}

		target := filepath.Join(dest, hdr.Name)

		switch hdr.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(target, hdr.FileInfo().Mode()); err != nil {
				return fmt.Errorf("creating dir %s: %w", target, err)
			}

		case tar.TypeReg:
			if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
				return fmt.Errorf("creating parent dirs for %s: %w", target, err)
			}

			f, err := os.OpenFile(target, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, hdr.FileInfo().Mode())
			if err != nil {
				return fmt.Errorf("creating file %s: %w", target, err)
			}

			if _, err := io.Copy(f, tarr); err != nil {
				f.Close()
				return fmt.Errorf("writing file %s: %w", target, err)
			}

			f.Close()

		case tar.TypeSymlink:
			if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
				return fmt.Errorf("creating parent dirs for symlink %s: %w", target, err)
			}
			if err := os.Symlink(hdr.Linkname, target); err != nil && !os.IsExist(err) {
				return fmt.Errorf("creating symlink %s: %w", target, err)
			}

		case tar.TypeLink:
			if err := os.Link(filepath.Join(dest, hdr.Linkname), target); err != nil {
				return fmt.Errorf("creating hard link %s: %w", target, err)
			}

		default:
			// ignore other types
		}
	}

	return nil
}
