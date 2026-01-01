package camo_test

import (
	"reflector/camo"
	"reflector/log"
	"testing"
)

func dismiss(a any, _ error) any {
	return a
}

func TestCamoController(t *testing.T) {
	log.SetDefaultLogger(log.NewLogger("camo", log.DEBUG))
	cc := camo.NewCamoController()
	templ := "docker://docker.io/z1xs4xg62/camo:example"
	cc.PreLoadCamo(templ)
	log.GetDefaultLogger().
		Info().
		Update("location", dismiss(cc.CamoLocation(templ))).
		Msg("camo loaded")
}
