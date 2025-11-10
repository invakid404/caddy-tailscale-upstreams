package tailscale

import (
	"context"
	"fmt"
	"strings"

	"tailscale.com/client/local"
	"tailscale.com/ipn/ipnstate"
)

type Client struct {
	instance *local.Client
}

func NewClient() *Client {
	return &Client{
		instance: &local.Client{},
	}
}

func (c *Client) Status(ctx context.Context) (*ipnstate.Status, error) {
	return c.instance.StatusWithoutPeers(ctx)
}

func (c *Client) GetPeersByTag(ctx context.Context, tag string) ([]*ipnstate.PeerStatus, error) {
	if !strings.HasPrefix(tag, tagPrefix) {
		tag = tagPrefix + tag
	}

	status, err := c.instance.Status(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get tailscale status: %w", err)
	}

	var matchingPeers []*ipnstate.PeerStatus
	for _, peer := range status.Peer {
		if !peer.Online {
			continue
		}

		if !hasTag(peer, tag) {
			continue
		}

		matchingPeers = append(matchingPeers, peer)
	}

	return matchingPeers, nil
}

func hasTag(peer *ipnstate.PeerStatus, tag string) bool {
	if peer.Tags == nil {
		return false
	}

	return peer.Tags.ContainsFunc(func(current string) bool {
		return current == tag
	})
}

const (
	tagPrefix = "tag:"
)
