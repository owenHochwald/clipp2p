package p2p

import (
	"context"
	"sync"

	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/p2p/discovery/mdns"
)

const discoveryServiceTag = "clipp2p"

// DiscoveryNotifee handles peer discovery events
type DiscoveryNotifee struct {
	mu      sync.RWMutex
	peers   map[peer.ID]peer.AddrInfo
	onFound func(peer.AddrInfo)
}

// NewDiscoveryNotifee creates a notifee that tracks discovered peers
func NewDiscoveryNotifee(onFound func(peer.AddrInfo)) *DiscoveryNotifee {
	return &DiscoveryNotifee{
		peers:   make(map[peer.ID]peer.AddrInfo),
		onFound: onFound,
	}
}

// HandlePeerFound is called when discovers a new peer
func (n *DiscoveryNotifee) HandlePeerFound(info peer.AddrInfo) {
	n.mu.Lock()
	_, exists := n.peers[info.ID]
	n.peers[info.ID] = info
	n.mu.Unlock()

	if !exists && n.onFound != nil {
		n.onFound(info)
	}
}

func (n *DiscoveryNotifee) Peers() []peer.AddrInfo {
	n.mu.RLock()
	defer n.mu.RUnlock()

	result := make([]peer.AddrInfo, 0, len(n.peers))
	for _, info := range n.peers {
		result = append(result, info)
	}
	return result
}

func (n *DiscoveryNotifee) RemovePeer(id peer.ID) {
	n.mu.Lock()
	defer n.mu.Unlock()
	delete(n.peers, id)
}

// Discovery wraps mDNS peer discovery
type Discovery struct {
	service mdns.Service
	notifee *DiscoveryNotifee
	node    *Node
	ctx     context.Context
	cancel  context.CancelFunc
}

// SetupDiscovery initializes mDNS discovery for the node
func (n *Node) SetupDiscovery(onPeerFound func(peer.AddrInfo)) (*Discovery, error) {
	ctx, cancel := context.WithCancel(n.ctx)

	notifee := NewDiscoveryNotifee(func(info peer.AddrInfo) {
		// Don't discover ourselves
		if info.ID == n.ID() {
			return
		}

		// Try to connect to new peer
		if err := n.host.Connect(ctx, info); err == nil {
			if onPeerFound != nil {
				onPeerFound(info)
			}
		}
	})

	service := mdns.NewMdnsService(n.host, discoveryServiceTag, notifee)
	if err := service.Start(); err != nil {
		cancel()
		return nil, err
	}

	return &Discovery{
		service: service,
		notifee: notifee,
		node:    n,
		ctx:     ctx,
		cancel:  cancel,
	}, nil
}

func (d *Discovery) Peers() []peer.AddrInfo {
	return d.notifee.Peers()
}

func (d *Discovery) Close() error {
	d.cancel()
	return d.service.Close()
}
