package p2p

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/stretchr/testify/assert"
)

func TestConnectionNotifier_OnConnect(t *testing.T) {
	ctx := context.Background()

	node1, err := NewNode(ctx)
	assert.NoError(t, err)
	defer node1.Close()

	node2, err := NewNode(ctx)
	assert.NoError(t, err)
	defer node2.Close()

	var mu sync.Mutex
	var connectedPeers []peer.ID

	node1.SetupConnectionNotifier(
		func(peerID peer.ID) {
			mu.Lock()
			connectedPeers = append(connectedPeers, peerID)
			mu.Unlock()
		},
		nil,
	)

	// Connect node1 to node2
	err = node1.Host().Connect(ctx, node2.AddrInfo())
	assert.NoError(t, err)

	// Wait for connection event
	time.Sleep(100 * time.Millisecond)

	mu.Lock()
	defer mu.Unlock()

	assert.Equal(t, 1, len(connectedPeers))
	assert.Equal(t, node2.ID(), connectedPeers[0])
}

func TestConnectionNotifier_BothCallbacks(t *testing.T) {
	ctx := context.Background()

	node1, err := NewNode(ctx)
	assert.NoError(t, err)
	defer node1.Close()

	node2, err := NewNode(ctx)
	assert.NoError(t, err)

	var mu sync.Mutex
	var connected, disconnected bool

	node1.SetupConnectionNotifier(
		func(peerID peer.ID) {
			mu.Lock()
			connected = true
			mu.Unlock()
		},
		func(peerID peer.ID) {
			mu.Lock()
			disconnected = true
			mu.Unlock()
		},
	)

	// Connect
	err = node1.Host().Connect(ctx, node2.AddrInfo())
	assert.NoError(t, err)
	time.Sleep(100 * time.Millisecond)

	mu.Lock()
	assert.True(t, connected, "should have received connect event")
	assert.False(t, disconnected, "should not have received disconnect event yet")
	mu.Unlock()

	// Disconnect
	node2.Close()
	time.Sleep(200 * time.Millisecond)

	mu.Lock()
	assert.True(t, disconnected, "should have received disconnect event")
	mu.Unlock()
}
