package p2p

import (
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/multiformats/go-multiaddr"
)

// ConnectionNotifier handles events
type ConnectionNotifier struct {
	onConnect    func(peer.ID)
	onDisconnect func(peer.ID)
}

func NewConnectionNotifier(onConnect, onDisconnect func(peer.ID)) *ConnectionNotifier {
	return &ConnectionNotifier{
		onConnect:    onConnect,
		onDisconnect: onDisconnect,
	}
}

func (n *ConnectionNotifier) Listen(network.Network, multiaddr.Multiaddr)      {}
func (n *ConnectionNotifier) ListenClose(network.Network, multiaddr.Multiaddr) {}

func (n *ConnectionNotifier) Connected(net network.Network, conn network.Conn) {
	peerID := conn.RemotePeer()

	// Only trigger callback if this is the first connection to this peer
	conns := net.ConnsToPeer(peerID)
	if len(conns) == 1 && n.onConnect != nil {
		n.onConnect(peerID)
	}
}

func (n *ConnectionNotifier) Disconnected(net network.Network, conn network.Conn) {
	peerID := conn.RemotePeer()

	// Only trigger callback if this was the last connection to this peer
	conns := net.ConnsToPeer(peerID)
	if len(conns) == 0 && n.onDisconnect != nil {
		n.onDisconnect(peerID)
	}
}
