package p2p

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
)

// ClipMessage is the packet sent between peers
type ClipMessage struct {
	Content   string    `json:"content"`
	Timestamp time.Time `json:"timestamp"`
	PeerName  string    `json:"peer_name"`
}

// StreamHandler manages protocol streams for messages
type StreamHandler struct {
	node      *Node
	onReceive func(from peer.ID, msg ClipMessage)
	mu        sync.RWMutex
	peerNames map[peer.ID]string
}

func NewStreamHandler(node *Node, onReceive func(from peer.ID, msg ClipMessage)) *StreamHandler {
	sh := &StreamHandler{
		node:      node,
		onReceive: onReceive,
		peerNames: make(map[peer.ID]string),
	}

	node.host.SetStreamHandler(ProtocolID, sh.handleStream)

	return sh
}

func (sh *StreamHandler) handleStream(stream network.Stream) {
	defer stream.Close()

	reader := bufio.NewReader(stream)
	remotePeer := stream.Conn().RemotePeer()

	for {
		line, err := reader.ReadBytes('\n')
		if err != nil {
			return
		}

		var msg ClipMessage
		if err := json.Unmarshal(line, &msg); err != nil {
			continue
		}

		if msg.PeerName != "" {
			sh.mu.Lock()
			sh.peerNames[remotePeer] = msg.PeerName
			sh.mu.Unlock()
		}

		if sh.onReceive != nil {
			sh.onReceive(remotePeer, msg)
		}
	}
}

func (sh *StreamHandler) SendClip(ctx context.Context, peerID peer.ID, msg ClipMessage) error {
	stream, err := sh.node.host.NewStream(ctx, peerID, ProtocolID)
	if err != nil {
		return fmt.Errorf("failed to open stream: %w", err)
	}
	defer stream.Close()

	data, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	data = append(data, '\n')
	if _, err := stream.Write(data); err != nil {
		return fmt.Errorf("failed to write message: %w", err)
	}

	return nil
}

func (sh *StreamHandler) Broadcast(ctx context.Context, msg ClipMessage) []error {
	peers := sh.node.host.Network().Peers()
	var errs []error

	for _, peerID := range peers {
		if err := sh.SendClip(ctx, peerID, msg); err != nil {
			errs = append(errs, err)
		}
	}

	return errs
}

func (sh *StreamHandler) GetPeerName(peerID peer.ID) string {
	sh.mu.RLock()
	defer sh.mu.RUnlock()

	if name, ok := sh.peerNames[peerID]; ok {
		return name
	}
	idStr := peerID.String()
	if len(idStr) > 8 {
		return idStr[:8] + "..."
	}
	return idStr
}

func (sh *StreamHandler) ConnectedPeers() []peer.ID {
	return sh.node.host.Network().Peers()
}
