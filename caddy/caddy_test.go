package caddy_test

import (
	"fmt"
	"reflector/caddy"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestCaddyVersion(t *testing.T) {
	v := "v1.2.3"
	cvV := caddy.LoadCaddyVersion(&v)
	fmt.Printf("%v\n", cvV)
	v = "1.2.3"
	cv := caddy.LoadCaddyVersion(&v)
	fmt.Printf("%v\n", cv)
	v = "1.2.3.2"
	cvMV := caddy.LoadCaddyVersion(&v)
	fmt.Printf("%v\n", cvMV)
	if !cmp.Equal(cvV, cv) {
		t.Fatalf("version mismatch 'v*' != '*'")
		return
	}
	if !cmp.Equal(cv, cvMV) {
		t.Fatalf("version mismatch '*' != '*.prerelease'")
	}
	cvVRepr := cvV.ReprV()
	fmt.Printf("%s\n", cvVRepr)
	if cvVRepr != "v1.2.3" {
		t.Fatalf("invalid version representation")
	}
}

func TestCaddyJSON(t *testing.T) {
	cj := caddy.NewCaddyJSON([]string{":8443"})
	cj.AddProxyLocation("example.com", "/", "localhost:8080")
	jsn := cj.Marshal()
	fmt.Print(string(jsn))
}

func TestCaddy(t *testing.T) {
	c := caddy.NewPortableCaddy("v2.10.0")
	c.AddProxyLocation("localhost", "/", "http://localhost:8080")
	c.Start()
	time.Sleep(10 * time.Second)
	c.Stop()
}
