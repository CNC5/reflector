package logic

import (
	"errors"
	"fmt"
	"io"
	"reflector/log"
	"reflector/utils"
	"strings"

	"gopkg.in/yaml.v3"
)

type configHeader struct {
	ApiVersion string `yaml:"apiVersion"`
	Kind       string `yaml:"kind"`
}

func LoadConfig(c []byte) (*reflectorConfigV1, error) {
	h := configHeader{}
	err := yaml.Unmarshal(c, &h)
	if err != nil {
		errorMsg := "failed unmarshal"
		log.GetDefaultLogger().Error().Update("err", err).Msg(errorMsg)
		return nil, errors.New(errorMsg)
	}
	switch h.Kind + h.ApiVersion {
	case "Reflectorv1":
		{
			newConf := reflectorConfigV1{}
			yaml.Unmarshal(c, &newConf)
			return &newConf, nil
		}
	default:
		{
			errorMsg := "unknown kind + version"
			log.GetDefaultLogger().Error().
				Update("kind", h.Kind).
				Update("version", h.ApiVersion).Msg(errorMsg)
			return nil, errors.New(errorMsg)
		}
	}
}

type reflectorConfigV1 struct {
	Spec reflectorConfigV1Spec `yaml:"spec"`
}

type reflectorConfigV1Spec struct {
	Camos     map[string]reflectorConfigV1SpecInboundCamo `yaml:"camos,omitempty"`
	Inbounds  []reflectorConfigV1SpecInbound              `yaml:"inbounds,omitempty"`
	Outbounds []reflectorConfigV1SpecOutbound             `yaml:"outbounds,omitempty"`
	Routes    []reflectorConfigV1SpecRoute                `yaml:"routes,omitempty"`
	Metrics   reflectorConfigV1SpecMetrics                `yaml:"metrics,omitempty"`
}

// BEGIN INBOUND
type reflectorConfigV1SpecInbound struct {
	Name       string                             `yaml:"name"`
	Type       string                             `yaml:"type"`
	Transport  string                             `yaml:"transport"`
	Listen     string                             `yaml:"listen,omitempty"`
	ListenPort int                                `yaml:"listen_port"`
	Users      []reflectorConfigV1SpecInboundUser `yaml:"users"`
	PrivateKey string                             `yaml:"private_key,omitempty"`
	XHTTPPath  string                             `yaml:"xhttpPath,omitempty"`
	Camo       string                             `yaml:"camo,omitempty"`
}

type reflectorConfigV1SpecInboundUser struct {
	Name    string `yaml:"name"`
	UUID    string `yaml:"uuid"`
	Flow    string `yaml:"flow"`
	ShortID string `yaml:"short_id"`
}

// END INBOUND

type reflectorConfigV1SpecInboundCamo struct {
	Security string `yaml:"security"`
	Template string `yaml:"template,omitempty"`
	FQDN     string `yaml:"fqdn"`
}

type reflectorConfigV1SpecOutbound struct {
	Name string `yaml:"name,omitempty"`
	Type string `yaml:"type,omitempty"`
}

type reflectorConfigV1SpecRoute struct {
	User     string `yaml:"user,omitempty"`
	Outbound string `yaml:"outbound,omitempty"`
}

type reflectorConfigV1SpecMetrics struct {
	Port   int    `yaml:"port"`
	Listen string `yaml:"listen,omitempty"`
}

func (r *reflector) ParseReflectorConfigV1(config io.Reader) error {
	rc := &reflectorConfigV1{}
	configBytes, err := io.ReadAll(config)
	if err != nil {
		log.GetDefaultLogger().Error().Msg("failed to read config bytes")
		return err
	}
	err = yaml.Unmarshal(configBytes, rc)
	if err != nil {
		log.GetDefaultLogger().Error().Msg("failed to unmarshal config")
		return err
	}

	possibleSecurityOptions := map[string]bool{
		"wtls":    true,
		"reality": true,
	}
	localSecurityOption := "wtls"
	possibleSecurityOptionsSlice := []string{}
	for secOpt, _ := range possibleSecurityOptions {
		possibleSecurityOptionsSlice = append(possibleSecurityOptionsSlice, secOpt)
	}

	// check and load all camo
	for camoName, camo := range rc.Spec.Camos {
		if _, exists := possibleSecurityOptions[camo.Security]; !exists {
			log.GetDefaultLogger().
				Error().
				Msgf(
					"camo security should be one of: %s",
					strings.Join(possibleSecurityOptionsSlice, ", "))
			continue
		}

		addresses, err := utils.NSLookup(camo.FQDN)
		if err != nil {
			log.GetDefaultLogger().
				Error().
				Update("err", err.Error()).
				Update("camo_name", camoName).
				Update("fqdn", camo.FQDN).
				Msg("failed to lookup the camo fqdn, skipping camo")
			continue
		}
		if len(addresses) == 0 {
			log.GetDefaultLogger().
				Error().
				Update("camo_name", camoName).
				Update("fqdn", camo.FQDN).
				Msg("fqdn resolves to 0 IPs, skipping camo")
			continue
		}
		if !utils.IsDomainPointingToThisHost(camo.FQDN) {
			log.GetDefaultLogger().Warning().
				Update("camo_name", camoName).
				Update("addresses", addresses).
				Update("fqdn", camo.FQDN).
				Msg("the fqdn does not point to this host")
		}
		if camo.Security == localSecurityOption {
			if camo.Template == "" {
				log.GetDefaultLogger().
					Error().
					Msgf(
						"%s security requires a template",
						camo.Security)
			}
			log.GetDefaultLogger().
				Info().
				Update("camo_name", camoName).
				Update("template", camo.Template).
				Msg("loading camo")
			err := r.CamoController.PreLoadCamo(camo.Template)
			if err != nil {
				log.GetDefaultLogger().
					Error().
					Update("err", err.Error()).
					Msg("failed to load camo")
				continue
			}
			continue
		}
	}

	// check and load all inbounds
	successfullInbounds := 0
	for _, inb := range rc.Spec.Inbounds {
		if !utils.IsPortBindable(inb.ListenPort) {
			log.GetDefaultLogger().
				Error().
				Update("inbound", inb.Name).
				Update("port", inb.ListenPort).
				Msg("specified port is not bindable. missing root/bind capabilities?")
			continue
		}
		// reality
		// 	- caddy:443 -> xray:xrayPort -> caddy:caddyPort/ext:443
		// xhttptls
		// 	- caddy:443 -> xray:xrayPort
		if inb.Type == "vless" {
			// default to direct binding on the specified port
			// in case the configuration is not known
			xinb := r.XrayCore.XrayConfig.EnsureInboundVless(
				inb.Name,
				"0.0.0.0",
				inb.ListenPort,
			)

			if inb.Camo != "" {
				camoSpec, exists := rc.Spec.Camos[inb.Camo]
				if !exists {
					log.GetDefaultLogger().Error().
						Update("camo_name", inb.Camo).
						Msg("undefined camo, skipping inbound")
					continue
				}

				// prepare everything web server
				xrayPorts, err := utils.FindFreePorts(1)
				if err != nil {
					log.GetDefaultLogger().
						Error().
						Update("err", err.Error()).
						Msg("failed to get free ports for xrayPort")
					continue
				}
				xrayPort := xrayPorts[0]
				xrayPath := "/"
				if inb.Transport == "xhttp" {
					xrayPath = inb.XHTTPPath
				}
				r.Caddy.AddProxyLocation(
					camoSpec.FQDN,
					inb.ListenPort,
					xrayPath+"*",
					fmt.Sprintf("127.0.0.1:%d", xrayPort),
				)
				if camoSpec.Security == localSecurityOption {
					camoLocation, err := r.CamoController.CamoLocation(camoSpec.Template)
					if err != nil {
						log.GetDefaultLogger().Error().
							Update("camo_name", inb.Camo).
							Msg("failed to load camo")
					}
					r.Caddy.AddRootStaticLocation(camoSpec.FQDN, inb.ListenPort, camoLocation)
				}

				// prepare everything xray
				xinb.Listen = "127.0.0.1"
				xinb.Port = xrayPort
				if camoSpec.Security == "reality" {
					xinb.SecurityReality(camoSpec.FQDN, []string{})
					for _, user := range inb.Users {
						xinb.EnsureShortID(user.ShortID)
					}
				}
				if camoSpec.Security == localSecurityOption {
					// security is provided by caddy
					xinb.SecurityNone()
				}
			}

			// transports
			if inb.Transport == "xhttp" {
				xinb.TransportXHTTPAutoParams(inb.XHTTPPath, "packet-up")
			}
			if inb.Transport == "tcp" {
				xinb.TransportTCP()
			}

			// users
			for _, user := range inb.Users {
				_ = user
				// clientLink := xinb.EnsureClientReturnClientLink(
				// 	user.UUID,
				// 	user.Flow,
				// 	user.Name)
				// _ = clientLink
			}
		} else {
			log.GetDefaultLogger().
				Error().
				Update("type", inb.Type).
				Msg("unrecognized/unimplemented inbound type")
			continue
		}
		successfullInbounds += 1
	}
	if successfullInbounds == 0 {
		log.GetDefaultLogger().
			Error().
			Msg("0 inbounds were successfully loaded, exiting")
		return errors.New("0 inbounds loaded successfully")
	}

	for _, ob := range rc.Spec.Outbounds {
		switch ob.Type {
		case "direct":
			{
				ob := r.XrayCore.XrayConfig.EnsureOutbound(ob.Name)
				ob.Protocol = "freedom"
			}
		default:
			{
				log.GetDefaultLogger().
					Error().
					Update("type", ob.Type).
					Msg("unrecognized/unimplemented outbound type")
			}
		}
	}
	return nil
}
