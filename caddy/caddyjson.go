package caddy

import (
	"encoding/json"
	"slices"
)

type caddyJSON struct {
	Apps        caddyJSONApp `json:"apps"`
	httpsListen []string
}

type caddyJSONApp struct {
	Http caddyJSONAppHTTP `json:"http"`
}

type caddyJSONAppHTTP struct {
	Servers  map[string]*caddyJSONAppHTTPServer `json:"servers,omitempty"`
	HTTPPort int                                `json:"http_port,omitempty"`
}

type caddyJSONAppHTTPServer struct {
	Listen []string                      `json:"listen,omitempty"`
	Routes []caddyJSONAppHTTPServerRoute `json:"routes,omitempty"`
}

func (cs *caddyJSONAppHTTPServer) findRouteHost(
	hostname string) (int, bool) {
	for id, route := range cs.Routes {
		if _, exists := route.findMatchHost(hostname); exists {
			return id, true
		}
	}
	return -1, false
}

func (cs *caddyJSONAppHTTPServer) ensureRouteToHost(hostname string) int {
	if cs.Routes == nil {
		cs.Routes = []caddyJSONAppHTTPServerRoute{}
	}
	routeId, exists := cs.findRouteHost(hostname)
	if !exists {
		cs.Routes = append(cs.Routes, caddyJSONAppHTTPServerRoute{
			Match: []caddyJSONAppHTTPServerRouteMatch{{
				Host: []string{hostname},
			}},
			Terminal: true,
		})
		return len(cs.Routes) - 1
	}
	return routeId
}

func (cs *caddyJSONAppHTTPServer) addRouteReverseProxy(
	hostname string, dial string, path string) *caddyJSONAppHTTPServer {
	route := &cs.Routes[cs.ensureRouteToHost(hostname)]
	subrouteHandler := &route.Handle[route.ensureHandleSubroute()]
	subrouteHandler.Routes = append(subrouteHandler.Routes,
		caddyJSONAppHTTPServerRouteHandleRoute{
			Handle: []caddyJSONAppHTTPServerRouteHandleRouteHandle{{
				Handler: "reverse_proxy",
				Upstreams: []caddyJSONAppHTTPServerRouteHandleRouteHandleUpstream{{
					Dial: dial,
				}},
			}},
			Match: []caddyJSONAppHTTPServerRouteHandleRouteMatch{{
				Path: []string{path},
			}},
		})
	return cs
}

func (cs *caddyJSONAppHTTPServer) addRouteRootStaticLocation(
	hostname string, directory string) *caddyJSONAppHTTPServer {
	route := &cs.Routes[cs.ensureRouteToHost(hostname)]
	subrouteHandler := &route.Handle[route.ensureHandleSubroute()]
	subrouteHandler.Routes = append(subrouteHandler.Routes,
		caddyJSONAppHTTPServerRouteHandleRoute{
			Handle: []caddyJSONAppHTTPServerRouteHandleRouteHandle{{
				Handler: "vars",
				Root:    directory,
			}, {
				Handler: "file_server",
			}},
		})
	return cs
}

type caddyJSONAppHTTPServerRoute struct {
	Match    []caddyJSONAppHTTPServerRouteMatch  `json:"match,omitempty"`
	Handle   []caddyJSONAppHTTPServerRouteHandle `json:"handle,omitempty"`
	Terminal bool                                `json:"terminal,omitempty"`
}

func (cr *caddyJSONAppHTTPServerRoute) ensureHandleSubroute() int {
	id, exists := cr.findHandlerSubroute()
	if !exists {
		cr.Handle = append(cr.Handle, caddyJSONAppHTTPServerRouteHandle{
			Handler: "subroute",
		})
		return len(cr.Handle) - 1
	}
	return id
}

func (cr *caddyJSONAppHTTPServerRoute) findMatchHost(hostname string) (int, bool) {
	for id, match := range cr.Match {
		if slices.Contains(match.Host, hostname) {
			return id, true
		}
	}
	return -1, false
}

func (cr *caddyJSONAppHTTPServerRoute) findHandlerSubroute() (int, bool) {
	for id, handle := range cr.Handle {
		if handle.Handler == "subroute" {
			return id, true
		}
	}
	return -1, false
}

type caddyJSONAppHTTPServerRouteMatch struct {
	Host []string `json:"host,omitempty"`
}

type caddyJSONAppHTTPServerRouteHandle struct {
	Handler string                                   `json:"handler,omitempty"`
	Routes  []caddyJSONAppHTTPServerRouteHandleRoute `json:"routes,omitempty"`
}

type caddyJSONAppHTTPServerRouteHandleRoute struct {
	Handle []caddyJSONAppHTTPServerRouteHandleRouteHandle `json:"handle,omitempty"`
	Match  []caddyJSONAppHTTPServerRouteHandleRouteMatch  `json:"match,omitempty"`
}

type caddyJSONAppHTTPServerRouteHandleRouteHandle struct {
	Handler   string                                                 `json:"handler,omitempty"`
	Root      string                                                 `json:"root,omitempty"`
	Upstreams []caddyJSONAppHTTPServerRouteHandleRouteHandleUpstream `json:"upstreams,omitempty"`
	Hide      []string                                               `json:"hide,omitempty"`
}

type caddyJSONAppHTTPServerRouteHandleRouteHandleUpstream struct {
	Dial string `json:"dial,omitempty"`
}

type caddyJSONAppHTTPServerRouteHandleRouteMatch struct {
	Path []string `json:"path,omitempty"`
}

func NewCaddyJSON(httpsListen []string) *caddyJSON {
	return &caddyJSON{httpsListen: httpsListen}
}

func (cj *caddyJSON) Marshal() []byte {
	bytes, err := json.Marshal(cj)
	if err != nil {
		panic(err)
	}
	return bytes
}

func (cj *caddyJSON) ensureServer(name string, listen []string) *caddyJSONAppHTTPServer {
	if cj.Apps.Http.Servers == nil {
		cj.Apps.Http.Servers = map[string]*caddyJSONAppHTTPServer{}
	}
	if _, exists := cj.Apps.Http.Servers[name]; !exists {
		cj.Apps.Http.Servers[name] =
			&caddyJSONAppHTTPServer{Listen: listen}
	}
	server := cj.Apps.Http.Servers[name]
	return server
}

func (cj *caddyJSON) AddProxyLocation(domain string, url string, proxyTarget string) {
	cj.ensureServer(domain, cj.httpsListen).
		addRouteReverseProxy(domain, proxyTarget, url)
}

func (cj *caddyJSON) AddRootStaticLocation(domain string, staticDir string) {
	cj.ensureServer(domain, cj.httpsListen).
		addRouteRootStaticLocation(domain, staticDir)
}
