package headers_test

import (
	"testing"
	"httpfromtcp/internal/headers"
	"httpfromtcp/internal/request"
	"io"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	//"fmt"
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

	n := r.numBytesPerRead
	if n > len(p) {
		n = len(p)
	}
	if n > len(r.data)-r.pos {
		n = len(r.data) - r.pos
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
test , ok:=  h.Get("host")
require.NoError(t, err)
require.NotNil(t, h)
assert.Equal(t, "localhost:42069", test)
assert.Equal(t, 25, n)
assert.True(t, done)
assert.True(t ,ok)

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
	test_host , ok := r.Headers.Get("host")
	test_host1 , ok1:= r.Headers.Get("User-Agent")
	test_host2 , ok2:= r.Headers.Get("Accept")
    // Use the .Get() method you defined earlier
    assert.Equal(t, "localhost:42069", test_host )
    assert.Equal(t, "curl/7.81.0", test_host1)
    assert.Equal(t, "*/*", test_host2)
	assert.True(t, ok)
	assert.True(t , ok1)
	assert.True(t , ok2)

    // Test: Malformed Header (Missing Colon)
    reader = &chunkReader{
        data:            "GET / HTTP/1.1\r\nHost localhost:42069\r\n\r\n",
        numBytesPerRead: 3,
    }
    r, err = request.RequestFromReader(reader)
    require.Error(t, err, "Should error when colon is missing")
}

func TestParseBody(t *testing.T) {
	// Test: Standard Body
reader := &chunkReader{
	data: "POST /submit HTTP/1.1\r\n" +
		"Host: localhost:42069\r\n" +
		"Content-Length: 13\r\n" +
		"\r\n" +
		"hello world!\n",
	numBytesPerRead: 3,
}
r, err := request.RequestFromReader(reader)

require.NoError(t, err)
require.NotNil(t, r)
assert.Len(t, r.Body, 13)
assert.Equal(t, "hello world!\n", string(r.Body))



// Test: Body shorter than reported content length
reader = &chunkReader{
	data: "POST /submit HTTP/1.1\r\n" +
		"Host: localhost:42069\r\n" +
		"Content-Length: 20\r\n" +
		"\r\n" +
		"partial content",
	numBytesPerRead: 3,
}
r, err = request.RequestFromReader(reader)
require.Error(t, err)
}