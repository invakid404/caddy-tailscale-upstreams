## Caddy Tailscale Upstreams

This repository provides a Caddy `reverse_proxy` upstream source that discovers
backends through the local Tailscale daemon. By pointing the module at a
Tailscale tag, it resolves all online peers with that tag and keeps the upstream
list fresh automatically.

- Caddy module ID: `http.reverse_proxy.upstreams.tailscale`
- Polls the local `tailscaled` every ten seconds for updates
- Works with any Caddy site block that uses `reverse_proxy`

### Requirements

- Go `1.25` or newer (for building from source)
- A working Tailscale node with the local daemon (`tailscaled`) running
- Tagged Tailscale devices that expose the desired service on a consistent port

### Installation

Add this module when building Caddy with `xcaddy`:

```shell
xcaddy build \
  --with github.com/invakid404/caddy-tailscale-upstreams
```

Alternatively, add the module to your existing `go.mod` and build your own
binary:

```shell
go get github.com/invakid404/caddy-tailscale-upstreams
go build
```

### Configuration

Reference the module inside a Caddy `reverse_proxy` definition.

This module has two required options:

- `tag`: should match the Tailscale ACL tag (the `tag:` prefix is automatically
  prepended).
- `port` is the port of the upstream service port you expect on those peers.

```caddyfile
example.com {
  reverse_proxy {
    dynamic tailscale {
      tag my-service
      port 8080
    }
  }
}
```

### How It Works

When Caddy provisions the module it:

1. Connects to the local Tailscale daemon via the official client API.
2. Fetches the full status and filters peers that are online and carry the
   configured tag.
3. Builds a reverse-proxy upstream list from each peer's Tailscale IP and the
   configured port.
4. Refreshes the list every ten seconds to track changes.

### License

This project is released into the public domain via the Unlicense. See `LICENSE`
for details.
