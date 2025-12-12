package xray

import (
	"encoding/json"
	"fmt"
	"os"
	"reflector/utils"
	"strings"
)

type XrayConfig struct {
	Log xrayConfigLog `json:"log,omitempty"`
	//  API              any                `json:"api,omitempty"`
	//  DNS              any                `json:"dns,omitempty"`
	Routing xrayConfigRouting `json:"routing,omitempty"`
	//  Policy    any                 `json:"policy,omitempty"`
	Inbounds       []*xrayConfigInbound `json:"inbounds,omitempty"`
	inboundsByTag  map[string]*xrayConfigInbound
	Outbounds      []*xrayConfigOutbound `json:"outbounds,omitempty"`
	outboundsByTag map[string]*xrayConfigOutbound
	// Transport        any                `json:"transport,omitempty"`
	// Stats            any                `json:"stats,omitempty"`
	// Reverse          any                `json:"reverse,omitempty"`
	// FakeDNS          any                `json:"fakedns,omitempty"`
	// Metrics          any                `json:"metrics,omitempty"`
	// Observatory      any                `json:"observatory,omitempty"`
	// BurstObservatory any                `json:"burstObservatory,omitempty"`
}

func (xc *XrayConfig) makeAll() {
	xc.inboundsByTag = make(map[string]*xrayConfigInbound)
	xc.outboundsByTag = make(map[string]*xrayConfigOutbound)
	xc.Routing.rulesByHash = make(map[string]*xrayConfigRoutingRule)
}

func (xc *XrayConfig) UnmarshalJSON(b []byte) error {
	type aliasXrayConfig XrayConfig
	newaxc := &aliasXrayConfig{}
	newaxc.inboundsByTag = make(map[string]*xrayConfigInbound)
	newaxc.outboundsByTag = make(map[string]*xrayConfigOutbound)
	newaxc.Routing.rulesByHash = make(map[string]*xrayConfigRoutingRule)
	err := json.Unmarshal(b, newaxc)
	if err != nil {
		return err
	}
	*xc = XrayConfig(*newaxc)
	for _, inb := range xc.Inbounds {
		xc.inboundsByTag[inb.Tag] = inb
	}
	for _, ob := range xc.Outbounds {
		xc.outboundsByTag[ob.Tag] = ob
	}
	for _, rr := range xc.Routing.Rules {
		xc.Routing.rulesByHash[routingRuleHashString(rr.InboundTag, rr.OutboundTag, rr.Port)] = rr
	}
	return nil
}

func (xc *XrayConfig) EnsureInbound(tag string) *xrayConfigInbound {
	inb, exists := xc.inboundsByTag[tag]
	if !exists {
		inb = &xrayConfigInbound{
			Tag: tag,
		}
		inb.StreamSettings.RealitySettings.shortIdMap = make(map[string]bool)
		xc.inboundsByTag[tag] = inb
		xc.Inbounds = append(xc.Inbounds, inb)
	}
	return inb
}

func (xc *XrayConfig) EnsureInboundVless(
	tag string,
	listen string,
	port int,
) *xrayConfigInbound {
	inb := xc.EnsureInbound(tag)
	inb.Listen = listen
	inb.Port = port
	inb.Protocol = "vless"
	if inb.Settings.Clients == nil {
		inb.Settings.Clients = []*xrayConfigInboundSettingsClient{}
	}
	if inb.Settings.clientsById == nil {
		inb.Settings.clientsById = make(map[string]*xrayConfigInboundSettingsClient)
	}
	inb.Settings.Decryption = "none"

	inb.Sniffing.Enabled = true
	inb.Sniffing.DestOverride = []string{
		"http",
		"tls",
		"quic",
	}
	return inb
}

func (inb *xrayConfigInbound) SecurityRealityAutoShortIDs(sni string) *xrayConfigInbound {
	shortIDs := inb.StreamSettings.RealitySettings.ShortIds
	if len(inb.StreamSettings.RealitySettings.ShortIds) == 0 {
		shortIDs = []string{
			utils.RandomHex(4),
		}
	}
	return inb.SecurityReality(sni, shortIDs)
}

func (inb *xrayConfigInbound) SecurityReality(sni string, shortIDs []string) *xrayConfigInbound {
	inb.StreamSettings.Security = "reality"
	inb.StreamSettings.RealitySettings.Show = false
	inb.StreamSettings.RealitySettings.Dest = fmt.Sprintf("%s:443", sni)
	inb.StreamSettings.RealitySettings.Xver = 0
	inb.StreamSettings.RealitySettings.ServerNames = []string{sni}
	if inb.StreamSettings.RealitySettings.PrivateKey == "" {
		newPrivKey, err := GenerateRealityX25519PrivateKey()
		if err != nil {
			panic(err)
		}
		inb.StreamSettings.RealitySettings.PrivateKey = newPrivKey
	}
	inb.StreamSettings.RealitySettings.ShortIds = shortIDs
	return inb
}

func (inb *xrayConfigInbound) SecurityNone() *xrayConfigInbound {
	inb.StreamSettings.Security = "none"
	return inb
}

func (inb *xrayConfigInbound) EnsureShortID(shortID string) *xrayConfigInbound {
	if _, exists := inb.StreamSettings.RealitySettings.shortIdMap[shortID]; !exists {
		inb.StreamSettings.RealitySettings.shortIdMap[shortID] = true
		inb.StreamSettings.RealitySettings.ShortIds = append(inb.StreamSettings.RealitySettings.ShortIds, shortID)
	}
	return inb
}

func (inb *xrayConfigInbound) TransportXHTTPAutoParams(
	xhttpPath string,
	xhttpMode string,
) *xrayConfigInbound {
	inb.TransportXHTTP(
		xhttpPath,
		xhttpMode,
		30,
		"10000000",
		"20-80",
		"100-1000",
	)
	return inb
}

func (inb *xrayConfigInbound) TransportXHTTP(
	xhttpPath string,
	xhttpMode string,
	maxBufferedPosts int,
	maxEachPostBytes string,
	streamUpServerSecs string,
	xPaddingBytes string,
) *xrayConfigInbound {
	inb.StreamSettings.Network = "xhttp"
	inb.StreamSettings.XHTTPSettings.Mode = xhttpMode
	inb.StreamSettings.XHTTPSettings.Path = xhttpPath
	inb.StreamSettings.XHTTPSettings.SCMaxBufferedPosts = maxBufferedPosts
	inb.StreamSettings.XHTTPSettings.SCMaxEachPostBytes = maxEachPostBytes
	inb.StreamSettings.XHTTPSettings.SCStreamUpServerSecs = streamUpServerSecs
	inb.StreamSettings.XHTTPSettings.XPaddingBytes = xPaddingBytes
	return inb
}

func (inb *xrayConfigInbound) TransportTCP() *xrayConfigInbound {
	inb.StreamSettings.Network = "tcp"
	return inb
}

func (xc *XrayConfig) EnsureOutbound(tag string) *xrayConfigOutbound {
	ob, exists := xc.outboundsByTag[tag]
	if !exists {
		ob = &xrayConfigOutbound{
			Tag: tag,
		}
		xc.outboundsByTag[tag] = ob
		xc.Outbounds = append(xc.Outbounds, ob)
	}
	return ob
}

func routingRuleHashString(
	inboundTag string,
	outboundTag string,
	port string,
) string {
	return fmt.Sprintf("%s%s%s", inboundTag, outboundTag, port)
}

func (xc *XrayConfig) EnsureRoutingRule(
	inboundTag string,
	outboundTag string,
	port string,
) *xrayConfigRoutingRule {
	rrHash := routingRuleHashString(inboundTag, outboundTag, port)
	rr, exists := xc.Routing.rulesByHash[rrHash]
	if exists {
		return rr
	}
	newrr := &xrayConfigRoutingRule{
		Type:        "field",
		InboundTag:  inboundTag,
		OutboundTag: outboundTag,
		Port:        port,
	}
	xc.Routing.Rules = append(xc.Routing.Rules, newrr)
	xc.Routing.rulesByHash[rrHash] = newrr
	return newrr
}

type xrayConfigLog struct {
	Access   string `json:"access,omitempty"`
	Loglevel string `json:"loglevel,omitempty"`
}

// BEGIN Routing
type xrayConfigRouting struct {
	DomainStrategy string                   `json:"domainStrategy,omitempty"`
	Rules          []*xrayConfigRoutingRule `json:"rules,omitempty"`
	rulesByHash    map[string]*xrayConfigRoutingRule
}

type xrayConfigRoutingRule struct {
	Type        string   `json:"type,omitempty"`
	OutboundTag string   `json:"outboundTag,omitempty"`
	InboundTag  string   `json:"inboundTag,omitempty"`
	Protocol    []string `json:"protocol,omitempty"`
	IP          []string `json:"ip,omitempty"`
	Port        string   `json:"port,omitempty"`
}

// END Routing

// BEGIN Inbound

type xrayConfigInbound struct {
	Tag            string                          `json:"tag,omitempty"`
	Listen         string                          `json:"listen,omitempty"`
	Port           int                             `json:"port,omitempty"`
	Protocol       string                          `json:"protocol,omitempty"`
	Settings       xrayConfigInboundSettings       `json:"settings,omitempty"`
	Sniffing       xrayConfigInboundSniffing       `json:"sniffing,omitempty"`
	StreamSettings xrayConfigInboundStreamSettings `json:"streamSettings,omitempty"`
}

func (xi *xrayConfigInbound) EnsureClientReturnClientLink(
	id string,
	flow string,
	email string,
) *XrayLink {
	xl := NewXrayLink(
		xi.Protocol,
		id,
		"", // external hostname is not available in this context,
		// xraylink panic prevents empty hostnames from marshaling
		xi.Port,
		email,
	)
	xl.Parameters.SNI = xi.StreamSettings.RealitySettings.ServerNames[0]
	xl.Parameters.Fingerprint = "chrome"
	pubkey, err := DeriveRealityX25519PublicKey(xi.StreamSettings.RealitySettings.PrivateKey)
	if err != nil {
		panic(err)
	}
	xl.Parameters.PublicKey = pubkey
	xl.Parameters.Security = xi.StreamSettings.Security
	xl.Parameters.ShortID = xi.StreamSettings.RealitySettings.ShortIds[0]
	xl.Parameters.Type = xi.StreamSettings.Network

	newc := &xrayConfigInboundSettingsClient{
		ID:    id,
		Flow:  flow,
		Email: email,
	}
	c, exists := xi.Settings.clientsById[id]
	if !exists {
		xi.Settings.clientsById[id] = newc
		xi.Settings.Clients = append(xi.Settings.Clients, newc)
	} else {
		*c = *newc
	}
	return xl
}

type xrayConfigInboundSettings struct {
	Clients     []*xrayConfigInboundSettingsClient `json:"clients,omitempty"`
	clientsById map[string]*xrayConfigInboundSettingsClient
	Decryption  string `json:"decryption,omitempty"`
}

func (xs *xrayConfigInboundSettings) UnmarshalJSON(b []byte) error {
	type aliasXrayConfigInboundSettings xrayConfigInboundSettings
	newxcs := &aliasXrayConfigInboundSettings{}
	newxcs.clientsById = make(map[string]*xrayConfigInboundSettingsClient)
	err := json.Unmarshal(b, newxcs)
	if err != nil {
		return err
	}
	*xs = xrayConfigInboundSettings(*newxcs)
	for _, cli := range xs.Clients {
		xs.clientsById[cli.ID] = cli
	}
	return nil
}

type xrayConfigInboundSettingsClient struct {
	Email string `json:"email,omitempty"`
	ID    string `json:"id,omitempty"`
	Flow  string `json:"flow,omitempty"`
}

type xrayConfigInboundSniffing struct {
	DestOverride []string `json:"destOverride,omitempty"`
	Enabled      bool     `json:"enabled,omitempty"`
}

type xrayConfigInboundStreamSettings struct {
	Network         string                                         `json:"network,omitempty"`
	Security        string                                         `json:"security,omitempty"`
	XHTTPSettings   xrayConfigInboundStreamSettingsXHTTPSettings   `json:"xhttpSettings,omitempty"`
	RealitySettings xrayConfigInboundStreamSettingsRealitySettings `json:"realitySettings,omitempty"`
}

type xrayConfigInboundStreamSettingsRealitySettings struct {
	Show        bool     `json:"show,omitempty"`
	Dest        string   `json:"dest,omitempty"`
	Xver        int      `json:"xver,omitempty"`
	ServerNames []string `json:"serverNames,omitempty"`
	PrivateKey  string   `json:"privateKey,omitempty"`
	shortIdMap  map[string]bool
	ShortIds    []string `json:"shortIds,omitempty"`
}

func (xs *xrayConfigInboundStreamSettingsRealitySettings) UnmarshalJSON(b []byte) error {
	type aliasXrayConfigInboundStreamSettingsRealitySettings xrayConfigInboundStreamSettingsRealitySettings
	newxrs := &aliasXrayConfigInboundStreamSettingsRealitySettings{}
	newxrs.shortIdMap = make(map[string]bool)
	err := json.Unmarshal(b, newxrs)
	if err != nil {
		return err
	}
	*xs = xrayConfigInboundStreamSettingsRealitySettings(*newxrs)
	for _, shortId := range xs.ShortIds {
		xs.shortIdMap[shortId] = true
	}
	return nil
}

type xrayConfigInboundStreamSettingsXHTTPSettings struct {
	Mode                 string `json:"mode,omitempty"`
	Path                 string `json:"path,omitempty"`
	SCMaxBufferedPosts   int    `json:"scMaxBufferedPosts,omitempty"`
	SCMaxEachPostBytes   string `json:"scMaxEachPostBytes,omitempty"`
	SCStreamUpServerSecs string `json:"scStreamUpServerSecs,omitempty"`
	XPaddingBytes        string `json:"xPaddingBytes,omitempty"`
}

// END Inbound

// BEGIN Outbound

type xrayConfigOutbound struct {
	Tag      string `json:"tag,omitempty"`
	Protocol string `json:"protocol,omitempty"`
}

// END Outbound

func (xc *XrayConfig) Marshal() []byte {
	bytes, err := json.MarshalIndent(xc, "", strings.Repeat(" ", 4))
	if err != nil {
		panic(err)
	}
	return bytes
}

func NewXrayConfig() *XrayConfig {
	newXJ := XrayConfig{}
	newXJ.makeAll()
	return &newXJ
}

func LoadXrayConfig(
	configPath string,
) (
	*XrayConfig,
	error,
) {
	configFile, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer configFile.Close()
	newXJ := NewXrayConfig()
	unmarsh := json.NewDecoder(configFile)
	unmarsh.DisallowUnknownFields()
	err = unmarsh.Decode(&newXJ)
	if err != nil {
		return nil, err
	}

	return newXJ, nil
}

func (xc *XrayConfig) DumpXrayConfig(
	configPath string,
) error {
	file, err := os.Create(configPath)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", strings.Repeat(" ", 4))
	return encoder.Encode(xc)
}
