package ui

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/libp2p/go-libp2p/core/peer"
)

// ClipEntry is a sync event
type ClipEntry struct {
	Content   string
	Timestamp time.Time
	IsLocal   bool
	PeerName  string
}

// PeerInfo is a connected peer
type PeerInfo struct {
	ID   peer.ID
	Name string
}

type Model struct {
	History    []ClipEntry
	Peers      []PeerInfo
	SyncActive bool
	MaxHistory int
	PeerName   string
	quitting   bool
}

type ClipReceivedMsg struct {
	Content   string
	Timestamp time.Time
	PeerName  string
	PeerID    peer.ID
}

type ClipSentMsg struct {
	Content   string
	Timestamp time.Time
}

type PeerConnectedMsg struct {
	ID   peer.ID
	Name string
}

type PeerDisconnectedMsg struct {
	ID peer.ID
}

type PeersUpdatedMsg struct {
	Peers []PeerInfo
}

type ToggleSyncMsg struct{}

type ClearHistoryMsg struct{}

func NewModel(peerName string) Model {
	return Model{
		History:    make([]ClipEntry, 0),
		Peers:      make([]PeerInfo, 0),
		SyncActive: true,
		MaxHistory: 50,
		PeerName:   peerName,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

// Update handles messages and updates the model
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			m.quitting = true
			return m, tea.Quit
		case "s":
			m.SyncActive = !m.SyncActive
			return m, nil
		case "c":
			m.History = make([]ClipEntry, 0)
			return m, nil
		}

	case ClipReceivedMsg:
		entry := ClipEntry{
			Content:   msg.Content,
			Timestamp: msg.Timestamp,
			IsLocal:   false,
			PeerName:  msg.PeerName,
		}
		m.History = append(m.History, entry)
		if len(m.History) > m.MaxHistory {
			m.History = m.History[1:]
		}
		return m, nil

	case ClipSentMsg:
		entry := ClipEntry{
			Content:   msg.Content,
			Timestamp: msg.Timestamp,
			IsLocal:   true,
			PeerName:  m.PeerName,
		}
		m.History = append(m.History, entry)
		if len(m.History) > m.MaxHistory {
			m.History = m.History[1:]
		}
		return m, nil

	case PeerConnectedMsg:
		for i, p := range m.Peers {
			if p.ID == msg.ID {
				m.Peers[i].Name = msg.Name
				return m, nil
			}
		}
		m.Peers = append(m.Peers, PeerInfo{
			ID:   msg.ID,
			Name: msg.Name,
		})
		return m, nil

	case PeerDisconnectedMsg:
		for i, p := range m.Peers {
			if p.ID == msg.ID {
				m.Peers = append(m.Peers[:i], m.Peers[i+1:]...)
				break
			}
		}
		return m, nil

	case PeersUpdatedMsg:
		m.Peers = msg.Peers
		return m, nil

	case ToggleSyncMsg:
		m.SyncActive = !m.SyncActive
		return m, nil

	case ClearHistoryMsg:
		m.History = make([]ClipEntry, 0)
		return m, nil
	}

	return m, nil
}

// IsQuitting returns whether the user has requested to quit
func (m Model) IsQuitting() bool {
	return m.quitting
}

// IsSyncActive returns the current sync state
func (m Model) IsSyncActive() bool {
	return m.SyncActive
}
