package utils

import (
	"archive/zip"
	"io"
	"os"
)

func UnpackZipSubpath(archivePath string, filename string, destinationPath string) {
	zar, err := zip.OpenReader(archivePath)
	if err != nil {
		panic(err)
	}
	defer zar.Close()

	for _, zfile := range zar.File {
		if zfile.Name != filename {
			continue
		}
		zFileDescriptor, err := zfile.Open()
		if err != nil {
			panic(err)
		}
		defer zFileDescriptor.Close()

		if zfile.FileInfo().IsDir() {
			os.MkdirAll(destinationPath, zfile.Mode())
			break
		}
		fileOnDiskDescriptor, err := os.OpenFile(
			destinationPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, zfile.Mode())
		if err != nil {
			panic(err)
		}
		defer fileOnDiskDescriptor.Close()

		_, err = io.Copy(fileOnDiskDescriptor, zFileDescriptor)
		if err != nil {
			panic(err)
		}
	}
}
