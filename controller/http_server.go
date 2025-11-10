package controller

import (
	"reflector/caddy"
	"reflector/log"
	"reflector/utils"
)

type HTTPServer interface {
	Start()
	Reload()
	AddRootStaticLocation(domain string, staticDir string)
	AddProxyLocation(domain string, url string, proxyTarget string)
	Stop()
}

func HTTPServerAutoSelect() HTTPServer {
	existingServer := utils.DetectExistingServer()
	switch existingServer {
	case "nginx":
		panic("can't handle existing nginx")
	case "caddy":
		panic("can't handle existing caddy")
	default:
		log.
			GetDefaultLogger().Info().
			Msg("did not find an existing http server, using portable caddy")
		return caddy.NewPortableCaddy("v2.9.0")
	}

}
