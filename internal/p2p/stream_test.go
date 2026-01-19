package p2p

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/stretchr/testify/assert"
)

func TestStreamHandler_GetPeerName_Unknown(t *testing.T) {
	ctx := context.Background()
	node, err := NewNode(ctx)
	assert.NoError(t, err)
	defer node.Close()

	handler := NewStreamHandler(node, nil)

	// Unknown peer should return shortened ID
	name := handler.GetPeerName("12D3KooWTestPeer123456789")
	assert.NotEmpty(t, name)
}

func TestTwoNodes_SendReceiveMessage(t *testing.T) {
	ctx := context.Background()

	// Create two nodes
	node1, err := NewNode(ctx)
	assert.NoError(t, err)
	defer node1.Close()

	node2, err := NewNode(ctx)
	assert.NoError(t, err)
	defer node2.Close()

	// Track received messages on node2
	var mu sync.Mutex
	var receivedMsgs []ClipMessage
	var receivedFrom peer.ID

	handler1 := NewStreamHandler(node1, nil)
	_ = handler1

	handler2 := NewStreamHandler(node2, func(from peer.ID, msg ClipMessage) {
		mu.Lock()
		receivedMsgs = append(receivedMsgs, msg)
		receivedFrom = from
		mu.Unlock()
	})
	_ = handler2

	// Connect node1 to node2
	err = node1.Host().Connect(ctx, node2.AddrInfo())
	assert.NoError(t, err)

	time.Sleep(100 * time.Millisecond)

	// Send message from node1 to node2
	testMsg := ClipMessage{
		Content:   "Test clipboard content",
		Timestamp: time.Now(),
		PeerName:  "Node1",
	}

	err = handler1.SendClip(ctx, node2.ID(), testMsg)
	assert.NoError(t, err)

	time.Sleep(100 * time.Millisecond)

	mu.Lock()
	defer mu.Unlock()

	assert.Equal(t, 1, len(receivedMsgs))
	assert.Equal(t, "Test clipboard content", receivedMsgs[0].Content)
	assert.Equal(t, "Node1", receivedMsgs[0].PeerName)
	assert.Equal(t, node1.ID(), receivedFrom)
}

func TestTwoNodes_Broadcast(t *testing.T) {
	ctx := context.Background()

	// Create three nodes
	node1, err := NewNode(ctx)
	assert.NoError(t, err)
	defer node1.Close()

	node2, err := NewNode(ctx)
	assert.NoError(t, err)
	defer node2.Close()

	node3, err := NewNode(ctx)
	assert.NoError(t, err)
	defer node3.Close()

	// Track received messages
	var mu sync.Mutex
	var node2Msgs, node3Msgs []ClipMessage

	handler1 := NewStreamHandler(node1, nil)

	NewStreamHandler(node2, func(from peer.ID, msg ClipMessage) {
		mu.Lock()
		node2Msgs = append(node2Msgs, msg)
		mu.Unlock()
	})

	NewStreamHandler(node3, func(from peer.ID, msg ClipMessage) {
		mu.Lock()
		node3Msgs = append(node3Msgs, msg)
		mu.Unlock()
	})

	// Connect node1 to node2 and node3
	err = node1.Host().Connect(ctx, node2.AddrInfo())
	assert.NoError(t, err)
	err = node1.Host().Connect(ctx, node3.AddrInfo())
	assert.NoError(t, err)

	time.Sleep(100 * time.Millisecond)

	// Broadcast from node1
	testMsg := ClipMessage{
		Content:   "Broadcast message",
		Timestamp: time.Now(),
		PeerName:  "Broadcaster",
	}

	errs := handler1.Broadcast(ctx, testMsg)
	assert.Empty(t, errs)

	time.Sleep(100 * time.Millisecond)

	mu.Lock()
	defer mu.Unlock()

	assert.Equal(t, 1, len(node2Msgs))
	assert.Equal(t, 1, len(node3Msgs))
	assert.Equal(t, "Broadcast message", node2Msgs[0].Content)
	assert.Equal(t, "Broadcast message", node3Msgs[0].Content)
}
