package utils

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/google/go-containerregistry/pkg/crane"
)

type ContainerImageController struct {
	DefaultRegistry string
}

func NewContainerImageController() *ContainerImageController {
	return &ContainerImageController{DefaultRegistry: "docker.io/library"}
}

func (c *ContainerImageController) UnpackImage(ref, dest string) error {
	switch {
	case strings.HasPrefix(ref, "docker://"):
		return c.unpackContainerImage(strings.TrimPrefix(ref, "docker://"), dest)
	case strings.HasPrefix(ref, "oci://"):
		return c.unpackContainerImage(strings.TrimPrefix(ref, "oci://"), dest)
	case strings.HasPrefix(ref, "./") || strings.HasPrefix(ref, "/"):
		return CopyDir(ref, dest)
	default:
		// Assume it's a plain image name, prepend default registry
		imageRef := fmt.Sprintf("%s/%s", strings.TrimSuffix(c.DefaultRegistry, "/"), ref)
		return c.unpackContainerImage(imageRef, dest)
	}
}

func (c *ContainerImageController) unpackContainerImage(ref, dest string) error {
	img, err := crane.Pull(ref)
	if err != nil {
		return err
	}
	buff := bytes.NewBuffer(make([]byte, 0))

	err = crane.Export(img, buff)
	if err != nil {
		return err
	}

	err = UnpackTar(buff, dest)
	if err != nil {
		return err
	}
	return nil
}
