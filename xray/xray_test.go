package xray_test

import (
	"fmt"
	"reflector/xray"
	"testing"
	"time"
)

func TestPortableXray(t *testing.T) {
	xr := xray.NewPortableXray("v25.9.11")
	fmt.Println("xray init complete")
	time.Sleep(1 * time.Second)
	xr.Start()
}
