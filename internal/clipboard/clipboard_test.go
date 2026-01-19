package clipboard

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMockClipboard_ReadWrite(t *testing.T) {
	mock := NewMockClipboard()

	content, err := mock.Read()

	assert.NoError(t, err, "Read() should succeed on an empty clipboard")
	assert.Empty(t, content, "Read() should return empty string on an empty clipboard")

	// Write and read back
	testContent := "Hello, ClipP2P!"
	assert.NoError(t, mock.Write(testContent), "Write() should succeed")

	content, err = mock.Read()
	assert.NoError(t, err, "Read() should succeed after Write()")
	assert.Equal(t, testContent, content, "Read() should return the same content as Write()")
}

func TestMockClipboard_SetContent(t *testing.T) {
	mock := NewMockClipboard()

	externalContent := "Copied from another app"
	mock.SetContent(externalContent)

	content, err := mock.Read()
	assert.NoError(t, err, "Read() should succeed after SetContent()")
	assert.Equal(t, externalContent, content, "Read() should return the same content as SetContent()")
}

func TestMockClipboard_Overwrite(t *testing.T) {
	mock := NewMockClipboard()

	assert.NoError(t, mock.Write("first"), "Write() should succeed")
	assert.NoError(t, mock.Write("second"), "Write() should succeed")

	content, err := mock.Read()
	assert.NoError(t, err, "Read() should succeed after Write()")
	assert.Equal(t, "second", content, "Read() should return the last written content")
}
