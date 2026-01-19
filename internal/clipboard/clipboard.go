package clipboard

// Clipboard defines the interface for OS clipboard operations
type Clipboard interface {
	// Read returns the current clipboard content
	Read() (string, error)

	// Write sets the clipboard content
	Write(content string) error
}
