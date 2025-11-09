package utils

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"os"

	"github.com/blakesmith/ar"
)

func UnpackDebSubpath(archivePath string, subpath string, destinationPath string) {
	dataArchiveName := "data.tar.gz"
	f, err := os.Open(archivePath)
	if err != nil {
		panic(err)
	}
	arReader := ar.NewReader(f)

	// Find data.tar.gz in the .deb archive
	for range 5 {
		currentFileHeaders, err := arReader.Next()
		if err != nil {
			panic(err)
		}
		if currentFileHeaders.Name == dataArchiveName {
			break
		}
	}
	gzipReader, err := gzip.NewReader(arReader)
	tarReader := tar.NewReader(gzipReader)

	for {
		currentFileHeaders, err := tarReader.Next()
		if err != nil {
			panic(err)
		}
		if currentFileHeaders.Name == subpath {
			break
		}
	}
	tf, err := os.Create(destinationPath)
	io.Copy(tf, tarReader)
}
