package xray

import (
	"encoding/json"
	"reflector/interfaces"
)

type xrayJSON struct {
	Log xrayJSONLog `json:"log,omitempty"`
	//  API              any                `json:"api,omitempty"`
	//  DNS              any                `json:"dns,omitempty"`
	Routing xrayJSONRouting `json:"routing,omitempty"`
	//	Policy           any                `json:"policy,omitempty"`
	Inbounds        []xrayJSONInbound `json:"inbounds,omitempty"`
	inboundsTagMap  map[string]*xrayJSONInbound
	Outbounds       []xrayJSONOutbound `json:"outbounds,omitempty"`
	outboundsTagMap map[string]*xrayJSONOutbound
	// Transport        any                `json:"transport,omitempty"`
	// Stats            any                `json:"stats,omitempty"`
	// Reverse          any                `json:"reverse,omitempty"`
	// FakeDNS          any                `json:"fakedns,omitempty"`
	// Metrics          any                `json:"metrics,omitempty"`
	// Observatory      any                `json:"observatory,omitempty"`
	// BurstObservatory any                `json:"burstObservatory,omitempty"`
}

func (xj *xrayJSON) UpdateInboundVlessXHTTP(
	tag string,
	listen string,
	port int,
	clients []interfaces.XrayInboundClient,
	xhttpMode string,
	path string,
) {
	inbPtr, exists := xj.inboundsTagMap[tag]
	xrayClients := []xrayJSONInboundSettingsClient{}
	for _, c := range clients {
		xrayClients = append(xrayClients, xrayJSONInboundSettingsClient{
			Email: c.Email,
			ID:    c.ID,
		})
	}
	if !exists {
		newInb := xrayJSONInbound{
			Tag:      tag,
			Listen:   listen,
			Port:     port,
			Protocol: "vless",
			Settings: xrayJSONInboundSettings{
				Clients:    xrayClients,
				Decryption: "none",
			},
			Sniffing: xrayJSONInboundSniffing{
				DestOverride: []string{
					"http",
					"tls",
					"quic",
					"fakedns",
				},
				Enabled: true,
			},
			StreamSettings: xrayJSONInboundStreamSettings{
				Network:  "xhttp",
				Security: "none",
				XHTTPSettings: xrayJSONInboundStreamSettingsXHTTPSettings{
					Mode:                 xhttpMode,
					Path:                 path,
					SCMaxBufferedPosts:   30,
					SCMaxEachPostBytes:   "10000000",
					SCStreamUpServerSecs: "20-80",
					XPaddingBytes:        "100-1000",
				},
			},
		}
		xj.Inbounds = append(xj.Inbounds, newInb)
		xj.inboundsTagMap[tag] = &xj.Inbounds[len(xj.Inbounds)-1]
		return
	}
	inbPtr.Listen = listen
	inbPtr.Port = port
	inbPtr.Protocol = "vless"
	inbPtr.Settings.Clients = xrayClients
	inbPtr.Settings.Decryption = "none"
	inbPtr.Sniffing.DestOverride = []string{
		"http",
		"tls",
		"quic",
		"fakedns",
	}
	inbPtr.Sniffing.Enabled = true
	inbPtr.StreamSettings.Network = "xhttp"
	inbPtr.StreamSettings.Security = "none"
	inbPtr.StreamSettings.XHTTPSettings.Mode = xhttpMode
	inbPtr.StreamSettings.XHTTPSettings.Path = path
	inbPtr.StreamSettings.XHTTPSettings.SCMaxBufferedPosts = 30
	inbPtr.StreamSettings.XHTTPSettings.SCMaxEachPostBytes = "10000000"
	inbPtr.StreamSettings.XHTTPSettings.SCStreamUpServerSecs = "20-80"
	inbPtr.StreamSettings.XHTTPSettings.XPaddingBytes = "100-1000"
}

func (xj *xrayJSON) UpdateOutboundFreedom(
	tag string,
) {
	_, exists := xj.outboundsTagMap[tag]
	if !exists {
		newOb := xrayJSONOutbound{
			Tag:      tag,
			Protocol: "direct",
		}
		xj.Outbounds = append(xj.Outbounds, newOb)
		xj.outboundsTagMap[tag] = &xj.Outbounds[len(xj.Outbounds)-1]
	}
}

type xrayJSONLog struct {
	Access   string `json:"access,omitempty"`
	Loglevel string `json:"loglevel,omitempty"`
}

// BEGIN Routing
type xrayJSONRouting struct {
	DomainStrategy string                `json:"domainStrategy,omitempty"`
	Rules          []xrayJSONRoutingRule `json:"rules,omitempty"`
}

type xrayJSONRoutingRule struct {
	OutboundTag string   `json:"outboundTag,omitempty"`
	Protocol    []string `json:"protocol,omitempty"`
	Type        string   `json:"type,omitempty"`
}

// END Routing

// BEGIN Inbound
type xrayJSONInbound struct {
	Tag            string                        `json:"tag,omitempty"`
	Listen         string                        `json:"listen,omitempty"`
	Port           int                           `json:"port,omitempty"`
	Protocol       string                        `json:"protocol,omitempty"`
	Settings       xrayJSONInboundSettings       `json:"settings,omitempty"`
	Sniffing       xrayJSONInboundSniffing       `json:"sniffing,omitempty"`
	StreamSettings xrayJSONInboundStreamSettings `json:"streamSettings,omitempty"`
}

type xrayJSONInboundSettings struct {
	Clients    []xrayJSONInboundSettingsClient `json:"clients,omitempty"`
	Decryption string                          `json:"decryption,omitempty"`
}

type xrayJSONInboundSettingsClient struct {
	Email string `json:"email,omitempty"`
	ID    string `json:"id,omitempty"`
}

type xrayJSONInboundSniffing struct {
	DestOverride []string `json:"destOverride,omitempty"`
	Enabled      bool     `json:"enabled,omitempty"`
}

type xrayJSONInboundStreamSettings struct {
	Network       string                                     `json:"network,omitempty"`
	Security      string                                     `json:"security,omitempty"`
	XHTTPSettings xrayJSONInboundStreamSettingsXHTTPSettings `json:"xhttpSettings,omitempty"`
}

type xrayJSONInboundStreamSettingsXHTTPSettings struct {
	Mode                 string `json:"mode,omitempty"`
	Path                 string `json:"path,omitempty"`
	SCMaxBufferedPosts   int    `json:"scMaxBufferedPosts,omitempty"`
	SCMaxEachPostBytes   string `json:"scMaxEachPostBytes,omitempty"`
	SCStreamUpServerSecs string `json:"scStreamUpServerSecs,omitempty"`
	XPaddingBytes        string `json:"xPaddingBytes,omitempty"`
}

// END Inbound

// BEGIN Outbound
type xrayJSONOutbound struct {
	Tag      string `json:"tag,omitempty"`
	Protocol string `json:"protocol,omitempty"`
}

// END Outbound

func (cj *xrayJSON) Marshal() []byte {
	bytes, err := json.Marshal(cj)
	if err != nil {
		panic(err)
	}
	return bytes
}

func NewXrayJSON() *xrayJSON {
	newXJ := xrayJSON{}
	newXJ.Inbounds = []xrayJSONInbound{}
	newXJ.Outbounds = []xrayJSONOutbound{}
	newXJ.inboundsTagMap = map[string]*xrayJSONInbound{}
	newXJ.outboundsTagMap = map[string]*xrayJSONOutbound{}
	newXJ.Log.Access = "./access.log"
	newXJ.Log.Loglevel = "info"
	return &newXJ
}
