package controller

import (
	"os"
	"os/signal"
	"reflector/log"
	"syscall"
)

type reflector struct {
	HTTPServer HTTPServer
	XrayCore   XrayCore
}

func NewReflector() *reflector {
	return &reflector{
		HTTPServer: HTTPServerAutoSelect(),
		XrayCore:   XrayCoreAutoSelect(),
	}
}

func (r *reflector) Start() {
	r.XrayCore.Start()
	// r.HTTPServer.Start()
	// r.HTTPServer.
	// 	AddProxyLocation("localhost", "/prox*", "127.0.0.1:8080")
	// r.HTTPServer.
	// 	AddRootStaticLocation("localhost", "/tmp/reflector/camo/")
	// r.HTTPServer.Reload()
	log.GetDefaultLogger().Info().Msg("reflector started")
}

func (r *reflector) Stop() {
	r.HTTPServer.Stop()
	r.XrayCore.Stop()
}

func (r *reflector) RunWithSignalHandling() {
	r.Start()
	defer r.Stop()
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs
}
