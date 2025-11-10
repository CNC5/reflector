package controller

import (
	"errors"
	"reflector/log"

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
	Spec reflectorConfigV1Spec `yaml:"spec,omitempty"`
}

type reflectorConfigV1Spec struct {
	Inbounds  []reflectorConfigV1SpecInbound  `yaml:"inbounds,omitempty"`
	Outbounds []reflectorConfigV1SpecOutbound `yaml:"outbounds,omitempty"`
	Routes    []reflectorConfigV1SpecRoute    `yaml:"routes,omitempty"`
	Metrics   reflectorConfigV1SpecMetrics    `yaml:"metrics,omitempty"`
}

// BEGIN INBOUND
type reflectorConfigV1SpecInbound struct {
	Name       string                             `yaml:"name,omitempty"`
	Type       string                             `yaml:"type,omitempty"`
	Listen     string                             `yaml:"listen,omitempty"`
	ListenPort int                                `yaml:"listen_port,omitempty"`
	Users      []reflectorConfigV1SpecInboundUser `yaml:"users,omitempty"`
	PrivateKey string                             `yaml:"private_key,omitempty"`
	XHTTPPath  string                             `yaml:"xhttpPath,omitempty"`
}

type reflectorConfigV1SpecInboundUser struct {
	Name    string `yaml:"name,omitempty"`
	UUID    string `yaml:"uuid,omitempty"`
	Flow    string `yaml:"flow,omitempty"`
	ShortID string `yaml:"short_id,omitempty"`
}

type reflectorConfigV1SpecInboundCamo struct {
	Type     string                                 `yaml:"type,omitempty"`
	Template string                                 `yaml:"template,omitempty"`
	FQDN     string                                 `yaml:"fqdn,omitempty"`
	Issuer   reflectorConfigV1SpecInboundCamoIssuer `yaml:"issuer,omitempty"`
}

type reflectorConfigV1SpecInboundCamoIssuer struct {
	Type  string `yaml:"type,omitempty"`
	Email string `yaml:"email,omitempty"`
}

// END INBOUND

type reflectorConfigV1SpecOutbound struct {
	Name string `yaml:"name,omitempty"`
	Type string `yaml:"type,omitempty"`
}

type reflectorConfigV1SpecRoute struct {
	User     string `yaml:"user,omitempty"`
	Outbound string `yaml:"outbound,omitempty"`
}

type reflectorConfigV1SpecMetrics struct {
	Port   int    `yaml:"port,omitempty"`
	Listen string `yaml:"listen,omitempty"`
}
