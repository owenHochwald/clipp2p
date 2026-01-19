package p2p

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/stretchr/testify/assert"
)

func TestDiscoveryNotifee_HandlePeerFound(t *testing.T) {
	var mu sync.Mutex
	var foundPeers []peer.AddrInfo

	notifee := NewDiscoveryNotifee(func(info peer.AddrInfo) {
		mu.Lock()
		foundPeers = append(foundPeers, info)
		mu.Unlock()
	})

	// Create a fake peer info
	fakeInfo := peer.AddrInfo{
		ID: "12D3KooWTestPeer1",
	}

	notifee.HandlePeerFound(fakeInfo)

	mu.Lock()
	assert.Equal(t, 1, len(foundPeers))
	mu.Unlock()

	notifee.HandlePeerFound(fakeInfo)

	mu.Lock()
	assert.Equal(t, 1, len(foundPeers), "duplicate peer should not trigger callback")
	mu.Unlock()
}

func TestDiscoveryNotifee_Peers(t *testing.T) {
	notifee := NewDiscoveryNotifee(nil)

	peer1 := peer.AddrInfo{ID: "12D3KooWTestPeer1"}
	peer2 := peer.AddrInfo{ID: "12D3KooWTestPeer2"}

	notifee.HandlePeerFound(peer1)
	notifee.HandlePeerFound(peer2)

	peers := notifee.Peers()
	assert.Equal(t, 2, len(peers))
}

func TestDiscoveryNotifee_RemovePeer(t *testing.T) {
	notifee := NewDiscoveryNotifee(nil)

	peerInfo := peer.AddrInfo{ID: "12D3KooWTestPeer1"}
	notifee.HandlePeerFound(peerInfo)

	assert.Equal(t, 1, len(notifee.Peers()))

	notifee.RemovePeer(peerInfo.ID)

	assert.Equal(t, 0, len(notifee.Peers()))
}

func TestNode_SetupDiscovery(t *testing.T) {
	ctx := context.Background()
	node, err := NewNode(ctx)
	assert.NoError(t, err)
	defer node.Close()

	discovery, err := node.SetupDiscovery(nil)
	assert.NoError(t, err)
	assert.NotNil(t, discovery)

	defer discovery.Close()
}

func TestTwoNodes_DiscoverEachOther(t *testing.T) {
	ctx := context.Background()

	node1, err := NewNode(ctx)
	assert.NoError(t, err)
	defer node1.Close()

	node2, err := NewNode(ctx)
	assert.NoError(t, err)
	defer node2.Close()

	var mu sync.Mutex
	var node1Found, node2Found bool

	// discovery on node1
	discovery1, err := node1.SetupDiscovery(func(info peer.AddrInfo) {
		mu.Lock()
		if info.ID == node2.ID() {
			node1Found = true
		}
		mu.Unlock()
	})
	assert.NoError(t, err)
	defer discovery1.Close()

	// discovery on node2
	discovery2, err := node2.SetupDiscovery(func(info peer.AddrInfo) {
		mu.Lock()
		if info.ID == node1.ID() {
			node2Found = true
		}
		mu.Unlock()
	})
	assert.NoError(t, err)
	defer discovery2.Close()

	// Wait for discovery
	time.Sleep(2 * time.Second)

	mu.Lock()
	defer mu.Unlock()
	assert.True(t, node1Found || node2Found)
}
