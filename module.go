package caddy_tailscale_upstreams

import "github.com/caddyserver/caddy/v2"

func init() {
	caddy.RegisterModule(Module{})
}

type Module struct{}

func (Module) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "http.reverse_proxy.upstreams.tailscale",
		New: func() caddy.Module {
			return new(Module)
		},
	}
}

var _ caddy.Module = (*Module)(nil)
