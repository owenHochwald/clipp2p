package clipboard

import (
	"golang.design/x/clipboard"
)

// SystemClipboard implements the Clipboard interface using the OS clipboard.
type SystemClipboard struct{}

// NewSystemClipboard creates a new SystemClipboard.
// It initializes the underlying clipboard library.
func NewSystemClipboard() (*SystemClipboard, error) {
	if err := clipboard.Init(); err != nil {
		return nil, err
	}
	return &SystemClipboard{}, nil
}

// Read returns the current OS clipboard content as text.
func (s *SystemClipboard) Read() (string, error) {
	data := clipboard.Read(clipboard.FmtText)
	return string(data), nil
}

// Write sets the OS clipboard content to the given text.
func (s *SystemClipboard) Write(content string) error {
	clipboard.Write(clipboard.FmtText, []byte(content))
	return nil
}
