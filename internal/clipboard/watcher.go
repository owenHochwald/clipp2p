package clipboard

import (
	"context"
	"sync"
	"time"
)

// ClipboardChange is a message to indicate clipboard content changed
type ClipboardChange struct {
	Content   string
	Timestamp time.Time
}

// Watcher listens to clipboard changes with a callback function
type Watcher struct {
	clipboard    Clipboard
	pollInterval time.Duration
	onChange     func(ClipboardChange)

	mu          sync.Mutex
	lastContent string
	running     bool
	stopCh      chan struct{}
	doneCh      chan struct{}
}

func NewWatcher(cb Clipboard, interval time.Duration, onChange func(ClipboardChange)) *Watcher {
	return &Watcher{
		clipboard:    cb,
		pollInterval: interval,
		onChange:     onChange,
		stopCh:       make(chan struct{}),
		doneCh:       make(chan struct{}),
	}
}

// Start polls the clipboard for changes.
func (w *Watcher) Start(ctx context.Context) error {
	w.mu.Lock()
	if w.running {
		w.mu.Unlock()
		return nil
	}
	w.running = true
	w.stopCh = make(chan struct{})
	w.doneCh = make(chan struct{})
	w.mu.Unlock()

	// Initialize lastContent with current clipboard state
	if content, err := w.clipboard.Read(); err == nil {
		w.mu.Lock()
		w.lastContent = content
		w.mu.Unlock()
	}

	defer func() {
		w.mu.Lock()
		w.running = false
		close(w.doneCh)
		w.mu.Unlock()
	}()

	ticker := time.NewTicker(w.pollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-w.stopCh:
			return nil
		case <-ticker.C:
			w.poll()
		}
	}
}

// poll checks for clipboard changes and fires the callback if changed
func (w *Watcher) poll() {
	content, err := w.clipboard.Read()
	if err != nil {
		return
	}

	w.mu.Lock()
	lastContent := w.lastContent
	w.mu.Unlock()

	if content != lastContent {
		w.mu.Lock()
		w.lastContent = content
		w.mu.Unlock()

		if w.onChange != nil {
			w.onChange(ClipboardChange{
				Content:   content,
				Timestamp: time.Now(),
			})
		}
	}
}

// Stop waits and stops the watcher
func (w *Watcher) Stop() {
	w.mu.Lock()
	if !w.running {
		w.mu.Unlock()
		return
	}
	stopCh := w.stopCh
	doneCh := w.doneCh
	w.mu.Unlock()

	close(stopCh)
	<-doneCh
}

func (w *Watcher) IsRunning() bool {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.running
}
