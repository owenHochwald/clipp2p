package p2p

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewNode_CreatesSuccessfully(t *testing.T) {
	ctx := context.Background()
	node, err := NewNode(ctx)

	assert.NoError(t, err, "NewNode() should succeed")
	assert.NotNil(t, node, "NewNode() should return a non-nil node")

	err = node.Close()

	assert.NoError(t, err, "Node should close successfully after creation")
}

func TestNode_HasListenAddresses(t *testing.T) {
	ctx := context.Background()
	node, err := NewNode(ctx)

	assert.NoError(t, err, "NewNode() should succeed")

	defer node.Close()

	addrs := node.Addrs()

	assert.NotEmpty(t, addrs, "node should have at least one listen address")

	for _, addr := range addrs {
		t.Logf("Listen address: %s", addr.String())
	}
}

func TestNode_HostAccessible(t *testing.T) {
	ctx := context.Background()
	node, err := NewNode(ctx)
	assert.NoError(t, err, "NewNode() should succeed")
	defer node.Close()

	h := node.Host()
	assert.NotNil(t, h, "node.Host() should return a non-nil host")

	// Verify host ID matches node ID
	assert.Equal(t, h.ID(), node.ID(), "Host().ID() should match node.ID()")
}
