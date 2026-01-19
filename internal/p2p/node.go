package p2p

import (
	"context"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/multiformats/go-multiaddr"
)

// ProtocolID is our custom protocol identifier
const ProtocolID = "/clipp2p/1.0.0"

// Node is a peer
type Node struct {
	host   host.Host
	ctx    context.Context
	cancel context.CancelFunc
}

// NewNode creates and starts a node
func NewNode(ctx context.Context) (*Node, error) {
	nodeCtx, cancel := context.WithCancel(ctx)

	h, err := libp2p.New(
		libp2p.ListenAddrStrings(
			"/ip4/0.0.0.0/tcp/0",
			"/ip6/::/tcp/0",
		),
	)
	if err != nil {
		cancel()
		return nil, err
	}

	return &Node{
		host:   h,
		ctx:    nodeCtx,
		cancel: cancel,
	}, nil
}

func (n *Node) ID() peer.ID {
	return n.host.ID()
}

func (n *Node) Addrs() []multiaddr.Multiaddr {
	return n.host.Addrs()
}

func (n *Node) Host() host.Host {
	return n.host
}

func (n *Node) Close() error {
	n.cancel()
	return n.host.Close()
}

func (n *Node) AddrInfo() peer.AddrInfo {
	return peer.AddrInfo{
		ID:    n.host.ID(),
		Addrs: n.host.Addrs(),
	}
}
