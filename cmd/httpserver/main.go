package main

import (
	"context"
	"crypto/sha256"
	"fmt"
	"httpfromtcp/cmd/middleware"
	"httpfromtcp/internal/request"
	"httpfromtcp/internal/response"
	"httpfromtcp/internal/server"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

//# In WSL
//ip addr show eth0 | grep -oP '(?<=inet\s)\d+(\.\d+){3}'

const port = 42069

func ToString(bytes []byte) string {
	out := ""
	for _, b := range bytes {
		out += fmt.Sprintf("%02x", b)
	}
	return out
}

func respond400() []byte {
	return []byte(`<html>
  <head>
    <title>400 Bad Request</title>
  </head>
  <body>
    <h1>Bad Request</h1>
    <p>Your request honestly kinda sucked.</p>
  </body>
</html>`)
}

func respond500() []byte {
	return []byte(`<html>
  <head>
    <title>500 Internal Server Error</title>
  </head>
  <body>
    <h1>Internal Server Error</h1>
    <p>Okay, you know what? This one is on me.</p>
  </body>
</html>`)
}

func respond200() []byte {
	return []byte(`<html>
  <head>
    <title>200 OK</title>
  </head>
  <body>
    <h1>Success!</h1>
    <p>Your request was an absolute banger.</p>
  </body>
</html>	`)
}	


func main() {

	mainHandler := server.Handler(func(w *response.Writer, req *request.Request) {
		h := response.GetDefaultHeaders(0)
		body := respond200()
		status := response.StatusOk

		target := req.RequestLine.RequestTarget

		switch {
		case target == "/yourproblem":
			body = respond400()
			status = response.StatusBadRequest
			h.Replace("Content-length", fmt.Sprintf("%d", len(body)))
			h.Replace("Content-type", "text/html")
			w.WriteStatusLine(status)
			w.WriteHeaders(*h)
			w.WriteBody(body)

		case target == "/myproblem":
			body = respond500()
			status = response.StatusInternalServerError
			h.Replace("Content-length", fmt.Sprintf("%d", len(body)))
			h.Replace("Content-type", "text/html")
			w.WriteStatusLine(status)
			w.WriteHeaders(*h)
			w.WriteBody(body)
			
		case target == "/httpbin/html":
	        		res, err := http.Get("https://httpbin.org/" + target[len("/httpbin/"):])
            if err != nil {
                body = respond500()
                status = response.StatusInternalServerError
                h.Replace("Content-length", fmt.Sprintf("%d", len(body)))
                h.Replace("Content-type", "text/html")
                w.WriteStatusLine(status)
                w.WriteHeaders(*h)
                w.WriteBody(body)
                return
            }
            defer res.Body.Close()
        
            // Read the entire response body
            fullBody, err := io.ReadAll(res.Body)
            if err != nil {
                body = respond500()
                status = response.StatusInternalServerError
                h.Replace("Content-length", fmt.Sprintf("%d", len(body)))
                h.Replace("Content-type", "text/html")
                w.WriteStatusLine(status)
                w.WriteHeaders(*h)
                w.WriteBody(body)
                return
            }
        
            // Wrap in HTML
            htmlResponse := fmt.Sprintf("<html><body><pre>%s</pre></body></html>", fullBody)
        
            // Send regular HTTP response (no chunked, no trailers)
            h.Delete("Content-length")
            h.Set("Content-Type", "text/html")
            h.Set("Content-Length", fmt.Sprintf("%d", len(htmlResponse)))
            w.WriteStatusLine(response.StatusOk)
            w.WriteHeaders(*h)
            w.WriteBody([]byte(htmlResponse))
            return

		case strings.HasPrefix(target, "/httpbin/"):
			res, err := http.Get("https://httpbin.org/" + target[len("/httpbin/"):])	
			if err != nil {
				body = respond500()
				status = response.StatusInternalServerError
				h.Replace("Content-length", fmt.Sprintf("%d", len(body)))
				h.Replace("Content-type", "text/html")
				w.WriteStatusLine(status)
				w.WriteHeaders(*h)
				w.WriteBody(body)
				return
			}
			defer res.Body.Close()

			w.WriteStatusLine(response.StatusOk)
			h.Delete("Content-length")
			h.Set("transfer-encoding", "chunked")
			h.Set("Content-Type", "text/html")
			h.Set("Trailer", "X-Content-SHA256, X-Content-Length")

			w.WriteHeaders(*h)

			fullBody := []byte{}
            data := make([]byte, 32)
	    	for {
             n, err := res.Body.Read(data)
             if n > 0 {
                fullBody = append(fullBody, data[:n]...)
                w.WriteBody([]byte(fmt.Sprintf("%x\r\n", n)))
                w.WriteBody(data[:n])
                w.WriteBody([]byte("\r\n"))
              }
              if err != nil {
                break
               }
        }
        
			//TODO: fixed a the chunked but...
			//curl -v -D - http://localhost:42069/httpbin/stream/2
             w.WriteBody([]byte("0\r\n"))
             
         
             out := sha256.Sum256(fullBody)
             shaStr := ToString(out[:])
             lenStr := fmt.Sprintf("%d", len(fullBody))
             
             
             trailerBytes := []byte{}
             trailerBytes = fmt.Appendf(trailerBytes, "X-Content-SHA256: %s\r\n", shaStr)
             trailerBytes = fmt.Appendf(trailerBytes, "X-Content-Length: %s\r\n", lenStr)
             
             
             trailerBytes = append(trailerBytes, "\r\n"...)
             
             w.WriteBody(trailerBytes)
             return

		case target == "/video":
			f, _ := os.ReadFile("assets/vim.mp4")
			h.Replace("Content-type", "video/mp4")
			h.Replace("content-length", fmt.Sprintf("%d", len(f)))
			w.WriteStatusLine(response.StatusOk)
			w.WriteHeaders(*h)
			w.WriteBody(f)
			
		case target == "/panic":
            panic("test panic")

		case target == "/echo":
			//forgettttt the \n registred nurse
			body = []byte(req.Body  + "\n")
			status = response.StatusOk
			h.Replace("Content-length", fmt.Sprintf("%d", len(body)))
			h.Replace("Content-type", "text/html")
			w.WriteStatusLine(status)
			w.WriteHeaders(*h)
			log.Printf("Body: %s", string(body))
			w.WriteBody(body )

	        return

		default:
			body = []byte("<html><body><h1>404 Not Found</h1></body></html>")
			status = response.StatusNotFound
			h.Replace("Content-length", fmt.Sprintf("%d", len(body)))
			h.Replace("Content-type", "text/html")
			w.WriteStatusLine(status)
			w.WriteHeaders(*h)
			w.WriteBody(body)
		}
	})

	finalhandler := middleware.Chain(mainHandler,
									 middleware.Recovery,
									 middleware.Logging,
									 )	

	s, err := server.Serve(port, finalhandler )

	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}

	log.Println("Server started on port", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server gracefully stopped")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := s.Shutdown(ctx); err != nil {
    log.Printf("Forced shutdown: %v", err)
	} else {
    log.Println("All connections finished, server stopped")
	}
}