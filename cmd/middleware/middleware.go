package middleware

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"httpfromtcp/internal/request"
	"httpfromtcp/internal/response"
	"httpfromtcp/internal/server"
	"log"
	"time"
)

type Middleware func(server.Handler) server.Handler

func Logging(next server.Handler) server.Handler {
    return func(w *response.Writer, req *request.Request) {
        start := time.Now()
        reqID := generateID()
		
		w.SetHeader("X-Request-ID", reqID)
        next(w, req)
        
        log.Printf("[%s] %s %s — %v",
            reqID,
            req.RequestLine.Method,
            req.RequestLine.RequestTarget,
            time.Since(start))
    }
}


func Recovery(next server.Handler) server.Handler {
	return func(w *response.Writer , req *request.Request){
		defer func() {
			if r := recover(); r != nil {
				log.Printf("PANIC recovered: %v", r)
				w.WriteStatusLine(response.StatusInternalServerError)
                h := response.GetDefaultHeaders(0)
                h.Set("Content-Type", "text/html")
                w.WriteHeaders(*h)
                w.WriteBody([]byte("<html><body><h1>500 Internal Server Error</h1></body></html>"))
			}		
		}()
		next(w , req)
	}
}


func Chain(handler server.Handler, middlewares ...func(server.Handler) server.Handler) server.Handler {
	for i := len(middlewares) - 1 ; i >= 0 ; i-- {
		handler = middlewares[i](handler)
	}
	return  handler
}

func generateID() string {
	bytes := make([]byte,16)
	if _, err := rand.Read(bytes); err != nil {
		return fmt.Sprintf("%d%x", time.Now().UnixNano(), time.Now().UnixNano())
	}
	return hex.EncodeToString(bytes)
}

