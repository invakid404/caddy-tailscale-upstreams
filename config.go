package caddy_tailscale_upstreams

import (
	"fmt"
	"strconv"

	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
)

func (m *Module) UnmarshalCaddyfile(d *caddyfile.Dispenser) error {
	for d.Next() {
		// No positional arguments expected after "tailscale"
		if d.NextArg() {
			return d.ArgErr()
		}

		// Handle optional block for sub-options
		for nesting := d.Nesting(); d.NextBlock(nesting); {
			key := d.Val()
			switch key {
			case "port":
				if !d.NextArg() {
					return d.ArgErr()
				}

				rawValue := d.Val()
				parsedValue, err := strconv.ParseUint(rawValue, 10, 16)
				if err != nil || parsedValue < 0 || parsedValue > 65535 {
					return d.Errf("invalid port '%s'", rawValue)
				}

				m.Port = uint16(parsedValue)

				// Ensure no extra args after the port value
				if d.NextArg() {
					return d.ArgErr()
				}
			case "tag":
				if !d.NextArg() {
					return d.ArgErr()
				}

				m.TargetTag = d.Val()

				// Ensure no extra args after the tag value
				if d.NextArg() {
					return d.ArgErr()
				}
			default:
				return d.Errf("unrecognized tailscale option '%s'", key)
			}
		}
	}

	// Enforce that the tag is required
	if m.TargetTag == "" {
		return fmt.Errorf("tag is required for tailscale upstreams")
	}

	return nil
}

var (
	_ caddyfile.Unmarshaler
)
