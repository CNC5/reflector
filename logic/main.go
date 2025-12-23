package logic

import (
	"os"
	"os/signal"
	"reflector/caddy"
	"reflector/camo"
	"reflector/log"
	"reflector/xray"
	"syscall"
)

type reflector struct {
	Caddy          *caddy.PortableCaddy
	XrayCore       *xray.PortableXray
	CamoController *camo.CamoController
}

func NewReflector(caddyVersion, xrayVersion string) *reflector {
	return &reflector{
		Caddy:          caddy.NewPortableCaddy(caddyVersion),
		XrayCore:       xray.NewPortableXray(xrayVersion),
		CamoController: camo.NewCamoController(),
	}
}

func (r *reflector) Start() {
	r.XrayCore.Start()
	r.Caddy.Start()
	r.Caddy.Reload()
	log.GetDefaultLogger().Info().Msg("reflector started")
}

func (r *reflector) Stop() {
	r.Caddy.Stop()
	r.XrayCore.Stop()
}

func (r *reflector) RunWithSignalHandling() {
	r.Start()
	defer r.Stop()
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs
}
