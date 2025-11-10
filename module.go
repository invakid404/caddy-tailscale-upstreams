package caddy_tailscale_upstreams

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp/reverseproxy"
	"github.com/invakid404/caddy-tailscale-upstreams/tailscale"
	"go.uber.org/zap"
)

func init() {
	caddy.RegisterModule(newModule())
}

type Module struct {
	TargetTag string `json:"tag,omitempty"`
	client    *tailscale.Client

	upstreams   []*reverseproxy.Upstream
	upstreamsMu sync.RWMutex
}

func newModule() *Module {
	return &Module{
		TargetTag:   "",
		client:      nil,
		upstreams:   []*reverseproxy.Upstream{},
		upstreamsMu: sync.RWMutex{},
	}
}

func (m *Module) Provision(ctx caddy.Context) error {
	ctx.Logger().Info("connecting to local tailscale daemon")

	m.client = tailscale.NewClient()
	status, err := m.client.Status(ctx)
	if err != nil {
		return fmt.Errorf("failed to get tailscale status: %w", err)
	}

	ctx.Logger().Info("established connection with local tailscale daemon", zap.String("version", status.Version))

	go m.fetchUpstreamsLoop(ctx)

	return nil
}

func (m *Module) fetchUpstreamsLoop(ctx caddy.Context) {
	timer := time.NewTimer(0)
	defer timer.Stop()

	for {
		select {
		case <-timer.C:
			err := m.fetchUpstreams(ctx)
			if err != nil {
				ctx.Logger().Error("failed to fetch tailscale upstreams", zap.Error(err))
			}
			timer.Reset(10 * time.Second)
		case <-ctx.Done():
			return
		}
	}
}

func (m *Module) fetchUpstreams(ctx caddy.Context) error {
	ctx.Logger().Debug("fetching upstreams", zap.String("tag", m.TargetTag))

	peers, err := m.client.GetPeersByTag(ctx, m.TargetTag)
	if err != nil {
		return fmt.Errorf("failed to get peers by tag: %w", err)
	}

	var upstreams []*reverseproxy.Upstream
	for _, peer := range peers {
		ip := peer.TailscaleIPs[0]

		upstream := reverseproxy.Upstream{
			Dial: ip.String(),
		}
		upstreams = append(upstreams, &upstream)
	}

	ctx.Logger().Debug("fetched upstreams", zap.Int("count", len(upstreams)))

	m.upstreamsMu.Lock()
	defer m.upstreamsMu.Unlock()

	m.upstreams = upstreams

	return nil
}

func (m *Module) GetUpstreams(request *http.Request) ([]*reverseproxy.Upstream, error) {
	m.upstreamsMu.RLock()
	defer m.upstreamsMu.RUnlock()

	return m.upstreams, nil
}

func (*Module) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID: "http.reverse_proxy.upstreams.tailscale",
		New: func() caddy.Module {
			return newModule()
		},
	}
}

var (
	_ caddy.Module                = (*Module)(nil)
	_ caddy.Provisioner           = (*Module)(nil)
	_ reverseproxy.UpstreamSource = (*Module)(nil)
)
