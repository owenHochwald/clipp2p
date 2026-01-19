package clipboard

import (
	"context"
	"github.com/stretchr/testify/assert"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestWatcher_CallbackOnChange(t *testing.T) {
	mock := NewMockClipboard()
	mock.SetContent("initial")

	var mu sync.Mutex
	var changes []ClipboardChange

	watcher := NewWatcher(mock, 10*time.Millisecond, func(change ClipboardChange) {
		mu.Lock()
		changes = append(changes, change)
		mu.Unlock()
	})

	ctx, cancel := context.WithCancel(context.Background())
	go watcher.Start(ctx)

	// Wait for watcher to start
	time.Sleep(20 * time.Millisecond)

	// Simulate clipboard change
	mock.SetContent("changed content")

	// Wait for poll to detect change
	time.Sleep(30 * time.Millisecond)

	cancel()
	time.Sleep(20 * time.Millisecond)

	mu.Lock()
	defer mu.Unlock()

	assert.Equal(t, 1, len(changes))
	assert.Equal(t, "changed content", changes[0].Content)
}

func TestWatcher_NoCallbackWhenUnchanged(t *testing.T) {
	mock := NewMockClipboard()
	mock.SetContent("same content")

	var ops atomic.Uint64

	watcher := NewWatcher(mock, 10*time.Millisecond, func(change ClipboardChange) {
		ops.Add(1)
	})

	ctx, cancel := context.WithCancel(context.Background())
	go watcher.Start(ctx)

	// Let several polls occur with unchanged content
	time.Sleep(50 * time.Millisecond)

	cancel()
	time.Sleep(20 * time.Millisecond)

	assert.Equal(t, 0, int(ops.Load()))
}

func TestWatcher_MultipleChanges(t *testing.T) {
	mock := NewMockClipboard()
	mock.SetContent("start")

	var mu sync.Mutex
	var changes []string

	watcher := NewWatcher(mock, 10*time.Millisecond, func(change ClipboardChange) {
		mu.Lock()
		changes = append(changes, change.Content)
		mu.Unlock()
	})

	ctx, cancel := context.WithCancel(context.Background())
	go watcher.Start(ctx)

	time.Sleep(20 * time.Millisecond)

	// First change
	mock.SetContent("first")
	time.Sleep(30 * time.Millisecond)

	// Second change
	mock.SetContent("second")
	time.Sleep(30 * time.Millisecond)

	cancel()
	time.Sleep(20 * time.Millisecond)

	mu.Lock()
	defer mu.Unlock()

	assert.Equal(t, 2, len(changes))
	assert.Equal(t, "first", changes[0])
	assert.Equal(t, "second", changes[1])
}

func TestWatcher_Stop(t *testing.T) {
	mock := NewMockClipboard()

	watcher := NewWatcher(mock, 10*time.Millisecond, func(change ClipboardChange) {})

	ctx := context.Background()
	go watcher.Start(ctx)

	// Wait for watcher to start
	time.Sleep(20 * time.Millisecond)

	assert.True(t, watcher.IsRunning())

	watcher.Stop()

	assert.False(t, watcher.IsRunning())
}
