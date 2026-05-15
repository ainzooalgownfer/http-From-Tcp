package middleware

import (
	"httpfromtcp/internal/request"
	"httpfromtcp/internal/response"
	"httpfromtcp/internal/server"
	"log"
	"time"
)

type Middleware func(server.Handler) server.Handler

func LoggingMiddleware(next server.Handler) server.Handler {
		return func(w *response.Writer , req *request.Request) {
			start := time.Now()
			next(w,req)
			log.Printf("%s %s - %v" , req.RequestLine.Method , 
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



