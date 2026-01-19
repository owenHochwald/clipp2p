package app

import (
	"context"
	"os"
	"sync"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/libp2p/go-libp2p/core/peer"

	"github.com/owenHochwald/clipp2p/internal/clipboard"
	"github.com/owenHochwald/clipp2p/internal/p2p"
	"github.com/owenHochwald/clipp2p/internal/ui"
)

type Config struct {
	PeerName     string
	PollInterval time.Duration
}

// DefaultConfig returns sensible defaults
func DefaultConfig() Config {
	hostname, _ := os.Hostname()
	if hostname == "" {
		hostname = "ClipP2P"
	}
	return Config{
		PeerName:     hostname,
		PollInterval: 500 * time.Millisecond,
	}
}

// App orchestrates the clipboard, p2p, and UI components
type App struct {
	config        Config
	ctx           context.Context
	cancel        context.CancelFunc
	clipboard     clipboard.Clipboard
	watcher       *clipboard.Watcher
	node          *p2p.Node
	discovery     *p2p.Discovery
	streamHandler *p2p.StreamHandler
	program       *tea.Program
	model         ui.Model

	mu              sync.Mutex
	lastClipContent string
	ignoreNextClip  bool
}

func New(cfg Config) *App {
	return &App{
		config: cfg,
		model:  ui.NewModel(cfg.PeerName),
	}
}

func (a *App) Start(ctx context.Context) error {
	a.ctx, a.cancel = context.WithCancel(ctx)

	// Initialize clipboard
	cb, err := clipboard.NewSystemClipboard()
	if err != nil {
		return err
	}
	a.clipboard = cb

	// Initialize P2P node
	a.node, err = p2p.NewNode(a.ctx)
	if err != nil {
		return err
	}

	a.node.SetupConnectionNotifier(a.handlePeerConnected, a.handlePeerDisconnected)
	a.streamHandler = p2p.NewStreamHandler(a.node, a.handleIncomingClip)
	a.discovery, err = a.node.SetupDiscovery(a.handlePeerFound)

	if err != nil {
		a.node.Close()
		return err
	}

	a.watcher = clipboard.NewWatcher(a.clipboard, a.config.PollInterval, a.handleClipboardChange)

	go a.watcher.Start(a.ctx)

	return nil
}

func (a *App) handleClipboardChange(change clipboard.ClipboardChange) {
	a.mu.Lock()
	if a.ignoreNextClip {
		a.ignoreNextClip = false
		a.mu.Unlock()
		return
	}

	if !a.model.IsSyncActive() {
		a.mu.Unlock()
		return
	}

	a.lastClipContent = change.Content
	a.mu.Unlock()

	msg := p2p.ClipMessage{
		Content:   change.Content,
		Timestamp: change.Timestamp,
		PeerName:  a.config.PeerName,
	}

	a.streamHandler.Broadcast(a.ctx, msg)

	if a.program != nil {
		a.program.Send(ui.ClipSentMsg{
			Content:   change.Content,
			Timestamp: change.Timestamp,
		})
	}
}

func (a *App) handleIncomingClip(from peer.ID, msg p2p.ClipMessage) {
	a.mu.Lock()
	if !a.model.IsSyncActive() {
		a.mu.Unlock()
		return
	}

	if msg.Content == a.lastClipContent {
		a.mu.Unlock()
		return
	}

	a.ignoreNextClip = true
	a.mu.Unlock()

	a.clipboard.Write(msg.Content)

	// Update UI
	if a.program != nil {
		a.program.Send(ui.ClipReceivedMsg{
			Content:   msg.Content,
			Timestamp: msg.Timestamp,
			PeerName:  msg.PeerName,
			PeerID:    from,
		})
	}
}

func (a *App) handlePeerFound(info peer.AddrInfo) {
	// mDNS discovery - peer found but not necessarily connected yet
	// The actual connection event will be handled by handlePeerConnected
}

func (a *App) handlePeerConnected(peerID peer.ID) {
	if a.streamHandler == nil {
		return
	}

	name := a.streamHandler.GetPeerName(peerID)

	if a.program != nil {
		a.program.Send(ui.PeerConnectedMsg{
			ID:   peerID,
			Name: name,
		})
	}
}

func (a *App) handlePeerDisconnected(peerID peer.ID) {
	if a.program != nil {
		a.program.Send(ui.PeerDisconnectedMsg{
			ID: peerID,
		})
	}
}

func (a *App) SetProgram(p *tea.Program) {
	a.program = p
}

func (a *App) GetModel() ui.Model {
	return a.model
}

func (a *App) Stop() {
	if a.watcher != nil {
		a.watcher.Stop()
	}
	if a.discovery != nil {
		a.discovery.Close()
	}
	if a.node != nil {
		a.node.Close()
	}
	if a.cancel != nil {
		a.cancel()
	}
}
