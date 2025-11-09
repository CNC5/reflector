package utils_test

import (
	"fmt"
	"reflector/utils"
	"testing"
)

func TestUnpackZipSubpath(t *testing.T) {
	utils.UnpackZipSubpath("/tmp/reflector/Xray-linux-64.zip", "xray", "/tmp/reflector/xray")
}

func TestFindFreePorts(t *testing.T) {
	fmt.Println(utils.FindFreePorts(1))
}
