package controller

import (
	"reflector/interfaces"
	"reflector/log"
	"reflector/xray"
)

type XrayCore interface {
	Start()
	Stop()
	Reload()
}

func XrayCoreAutoSelect() XrayCore {
	return xray.NewPortableXray("v25.9.11")
}

type XrayConfig interface {
	UpdateInboundVlessXHTTP(
		tag string,
		listen string,
		port int,
		clients []interfaces.XrayInboundClient,
		xhttpMode string,
		path string,
	)
	UpdateOutboundFreedom(
		tag string,
	)
}

func ParseReflectorConfigV1(rc *reflectorConfigV1, xc XrayConfig) {
	for _, inb := range rc.Spec.Inbounds {
		switch inb.Type {
		case "vless-xhttp":
			{
				inbClients := []interfaces.XrayInboundClient{}
				for _, iu := range inb.Users {
					inbClients = append(inbClients, interfaces.XrayInboundClient{
						Email: iu.Name,
						ID:    iu.UUID,
					})
				}
				xc.UpdateInboundVlessXHTTP(
					inb.Name,
					inb.Listen,
					inb.ListenPort,
					inbClients,
					"stream-up",
					inb.XHTTPPath,
				)
			}
		default:
			{
				log.GetDefaultLogger().Error().
					Update("type", inb.Type).Msg("Unrecognized inbound type")
			}
		}
	}
	for _, ob := range rc.Spec.Outbounds {
		switch ob.Type {
		case "direct":
			{
				xc.UpdateOutboundFreedom(ob.Name)
			}
		default:
			{
				log.GetDefaultLogger().Error().
					Update("type", ob.Type).Msg("Unrecognized inbound type")
			}
		}
	}
}
