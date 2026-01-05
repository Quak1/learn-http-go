package headers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseFieldLine(t *testing.T) {
	// Test: Valid single header
	headers := NewHeaders()
	data := []byte("Host: localhost:42069\r\n\r\n")
	n, done, err := headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:42069", headers.Get("Host"))
	assert.Equal(t, 23, n)
	assert.False(t, done)

	// Test: Valid single header with extra whitespace
	headers = NewHeaders()
	data = []byte("       Host:         localhost:42069            \r\n\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:42069", headers.Get("Host"))
	assert.Equal(t, 50, n)
	assert.False(t, done)

	// Test: Valid header multiple value
	headers = NewHeaders()
	headers.Parse([]byte("Host: localhost:42069\r\n"))
	n, done, err = headers.Parse([]byte("host: :6969\r\n"))
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:42069, :6969", headers.Get("host"))
	assert.Equal(t, 13, n)
	assert.False(t, done)

	// Test: Valid done
	headers = NewHeaders()
	headers.Parse([]byte("Host: localhost:42069\r\n\r\n"))
	n, done, err = headers.Parse([]byte("\r\n"))
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:42069", headers.Get("Host"))
	assert.Equal(t, 2, n)
	assert.True(t, done)

	// Test: Invalid spacing header
	headers = NewHeaders()
	data = []byte("       Host : localhost:42069       \r\n\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)

	// Test: Invalid key characters
	headers = NewHeaders()
	data = []byte("Host;: localhost:42069\r\n\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)
}
