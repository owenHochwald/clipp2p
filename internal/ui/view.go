package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

const version = "v1.0"

var (
	primaryColor   = lipgloss.Color("#00FF00") // Green
	secondaryColor = lipgloss.Color("#00AAAA") // Cyan
	dimColor       = lipgloss.Color("#666666") // Gray
	errorColor     = lipgloss.Color("#FF0000") // Red
	localColor     = lipgloss.Color("#FFFF00") // Yellow
	remoteColor    = lipgloss.Color("#00FFFF") // Cyan

	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(primaryColor)

	connectedStyle = lipgloss.NewStyle().
			Foreground(primaryColor).
			Bold(true)

	disconnectedStyle = lipgloss.NewStyle().
				Foreground(errorColor).
				Bold(true)

	dividerStyle = lipgloss.NewStyle().
			Foreground(dimColor)

	timestampStyle = lipgloss.NewStyle().
			Foreground(dimColor)

	localTagStyle = lipgloss.NewStyle().
			Foreground(localColor).
			Bold(true)

	remoteTagStyle = lipgloss.NewStyle().
			Foreground(remoteColor).
			Bold(true)

	contentStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF"))

	footerStyle = lipgloss.NewStyle().
			Foreground(dimColor)

	keyStyle = lipgloss.NewStyle().
			Foreground(secondaryColor).
			Bold(true)

	syncOnStyle = lipgloss.NewStyle().
			Foreground(primaryColor).
			Bold(true)

	syncOffStyle = lipgloss.NewStyle().
			Foreground(errorColor).
			Bold(true)
)

func (m Model) View() string {
	if m.quitting {
		return "Goodbye!\n"
	}

	var b strings.Builder

	// Title
	title := titleStyle.Render(fmt.Sprintf("CLIP-P2P [%s]", version))
	b.WriteString(title)
	b.WriteString("\n")

	// Divider
	divider := dividerStyle.Render(strings.Repeat("─", 50))
	b.WriteString(divider)
	b.WriteString("\n")

	// Connection status
	b.WriteString(m.renderStatus())
	b.WriteString("\n\n")

	// History section
	b.WriteString(m.renderHistory())

	// Footer divider
	b.WriteString(divider)
	b.WriteString("\n")

	// Footer with controls
	b.WriteString(m.renderFooter())
	b.WriteString("\n")

	return b.String()
}

func (m Model) renderStatus() string {
	peerCount := len(m.Peers)

	var statusIcon, statusText string
	var style lipgloss.Style

	if peerCount > 0 {
		statusIcon = "[●]"
		style = connectedStyle
		peerNames := m.getPeerNames()
		statusText = fmt.Sprintf("CONNECTED (%d Peers: %s)", peerCount, peerNames)
	} else {
		statusIcon = "[○]"
		style = disconnectedStyle
		statusText = "SEARCHING..."
	}

	return style.Render(statusIcon + " " + statusText)
}

func (m Model) getPeerNames() string {
	if len(m.Peers) == 0 {
		return ""
	}

	names := make([]string, 0, len(m.Peers))
	for _, p := range m.Peers {
		name := p.Name
		if name == "" {
			idStr := p.ID.String()
			if len(idStr) > 8 {
				name = idStr[:8] + "..."
			} else {
				name = idStr
			}
		}
		names = append(names, name)
	}

	// Limit display to first 3 peers
	if len(names) > 3 {
		return strings.Join(names[:3], ", ") + fmt.Sprintf(" +%d more", len(names)-3)
	}
	return strings.Join(names, ", ")
}

func (m Model) renderHistory() string {
	var b strings.Builder

	b.WriteString("HISTORY:\n")

	if len(m.History) == 0 {
		emptyStyle := lipgloss.NewStyle().Foreground(dimColor)
		b.WriteString(emptyStyle.Render("  No sync history yet..."))
		b.WriteString("\n\n")
		return b.String()
	}

	// Show last 10 entries (most recent at bottom)
	start := 0
	if len(m.History) > 10 {
		start = len(m.History) - 10
	}

	for _, entry := range m.History[start:] {
		b.WriteString(m.renderHistoryEntry(entry))
		b.WriteString("\n")
	}

	b.WriteString("\n")
	return b.String()
}

func (m Model) renderHistoryEntry(entry ClipEntry) string {
	// Timestamp
	ts := timestampStyle.Render(entry.Timestamp.Format("3:04 PM"))

	// Tag
	var tag string
	if entry.IsLocal {
		tag = localTagStyle.Render("[Local]")
	} else {
		tag = remoteTagStyle.Render("[Remote]")
	}

	// Content
	content := truncateContent(entry.Content, 35)
	contentRendered := contentStyle.Render(content)

	return fmt.Sprintf("  %s  %s  %s", ts, tag, contentRendered)
}

func truncateContent(content string, maxLen int) string {
	content = strings.ReplaceAll(content, "\n", " ")
	content = strings.ReplaceAll(content, "\r", "")

	content = strings.TrimSpace(content)

	if len(content) > maxLen {
		return content[:maxLen-3] + "..."
	}
	return content
}

func (m Model) renderFooter() string {
	quit := keyStyle.Render("(q)") + " Quit"
	toggle := keyStyle.Render("(s)") + " Toggle Sync "

	var syncStatus string
	if m.SyncActive {
		syncStatus = syncOnStyle.Render("[ON]")
	} else {
		syncStatus = syncOffStyle.Render("[OFF]")
	}

	clear := keyStyle.Render("(c)") + " Clear History"

	return footerStyle.Render(fmt.Sprintf("%s  %s%s  %s", quit, toggle, syncStatus, clear))
}
