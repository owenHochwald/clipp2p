package clipboard

import "sync"

// MockClipboard is a mock clipboard that stores content in memory
type MockClipboard struct {
	mu      sync.RWMutex
	content string
}

func NewMockClipboard() *MockClipboard {
	return &MockClipboard{}
}

func (m *MockClipboard) Read() (string, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.content, nil
}

func (m *MockClipboard) Write(content string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.content = content
	return nil
}

// SetContent is a test helper to simulate external clipboard changes.
func (m *MockClipboard) SetContent(content string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.content = content
}
