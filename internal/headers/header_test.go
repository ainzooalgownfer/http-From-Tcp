package headers_test 

import (
	"testing"
	"MODULE_NAME/internal/headers"
	"MODULE_NAME/internal/request"
	"io"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type chunkReader struct {
	data            string
	numBytesPerRead int
	pos             int
}

func (r *chunkReader) Read(p []byte) (int, error) {
	if r.pos >= len(r.data) {
		return 0, io.EOF
	}

	// determine how many bytes to read this time
	n := r.numBytesPerRead
	remaining := len(r.data) - r.pos

	if n > remaining {
		n = remaining
	}

	if n > len(p) {
		n = len(p)
	}

	copy(p, r.data[r.pos:r.pos+n])
	r.pos += n

	return n, nil
}


func TestHeaderParse(t *testing.T) {
	// Test: Valid single header
h := headers.NewHeaders()
data := []byte("Host: localhost:42069\r\n\r\n")
n, done, err := h.Parse(data)
require.NoError(t, err)
require.NotNil(t, h)
assert.Equal(t, "localhost:42069", h.Get("Host"))
assert.Equal(t, 25, n)
assert.True(t, done)

// Test: Invalid spacing header
h = headers.NewHeaders()
data = []byte("       Host: localhost:42069\r\n\r\n")
n, done, err = h.Parse(data)
require.Error(t, err)
assert.Equal(t, 0, n)
assert.False(t, done)
}

func TestParseHeaders(t *testing.T) {
    // Test: Standard Headers
    reader := &chunkReader{
        data:            "GET / HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n",
        numBytesPerRead: 3,
    }
    r, err := request.RequestFromReader(reader)
    
    require.NoError(t, err)
    require.NotNil(t, r)

    // Use the .Get() method you defined earlier
    assert.Equal(t, "localhost:42069", r.Headers.Get("Host"))
    assert.Equal(t, "curl/7.81.0", r.Headers.Get("User-Agent"))
    assert.Equal(t, "*/*", r.Headers.Get("Accept"))

    // Test: Malformed Header (Missing Colon)
    reader = &chunkReader{
        data:            "GET / HTTP/1.1\r\nHost localhost:42069\r\n\r\n",
        numBytesPerRead: 3,
    }
    r, err = request.RequestFromReader(reader)
    require.Error(t, err, "Should error when colon is missing")
}