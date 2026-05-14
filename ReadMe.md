# HTTP Server from Scratch

A custom HTTP/1.1 server built in Go without using the standard `net/http` package.  
Implements low‑level TCP handling, request parsing, response writing, and supports chunked encoding with trailers.

## Features

- **Pure TCP listener** – accepts connections and spawns goroutines per request.
- **Custom request parser** – reads and parses HTTP/1.1 requests.
- **Custom response writer** – builds status line, headers, and body.
- **Multiple endpoints**:
  - `/yourproblem` → 400 Bad Request (HTML)
  - `/myproblem`  → 500 Internal Server Error (HTML)
  - `/httpbin/*`  → proxies requests to `https://httpbin.org`, streams response with chunked encoding and SHA‑256 trailer.
  - `/video`      → serves a local MP4 file (`assets/vim.mp4`)
  - any other path → 404 Not Found (HTML)
- **Chunked transfer encoding** – streams large responses without pre‑computing `Content-Length`.
- **Trailer headers** – sends `X-Content-SHA256` and `X-Content-Length` after the body.


