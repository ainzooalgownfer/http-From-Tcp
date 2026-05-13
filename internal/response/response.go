package response

import (
	"fmt"
	"httpfromtcp/internal/headers"
	"io"
)


type Response struct {

}


type StatusCode int 
const (
	StatusOk StatusCode = 200
	StatusBadRequest StatusCode = 400
	StatusNotFound StatusCode = 404
	StatusInternalServerError StatusCode = 500
)

 func WriteStatusLine(w io.Writer, statusCode StatusCode) error {
	statusLine := []byte{}
	switch statusCode {
	case StatusOk:  statusLine = []byte("HTTP/1.1 200 OK\r\n")
	case StatusBadRequest: statusLine = []byte("HTTP/1.1 400 BAD REQUEST\r\n")
	case StatusInternalServerError: statusLine= []byte("HTTP/1.1 500 INTERNAL SERVER ERROR\r\n")
	default: return fmt.Errorf("unowkn error code")
	}
	
	_, err := w.Write(statusLine)
	return err

 }

 func GetDefaultHeaders(contentLen int) *headers.Headers {
	h := headers.NewHeaders()
	h.Set("Content-Length", fmt.Sprintf("%d", contentLen))
	h.Set("Connection", "closed")
	h.Set("Content-Type", "text/plain")

	return h
 }

 func WriteHeaders(w io.Writer, h *headers.Headers) error {
	b := []byte{}
	h.Foreach(func(n, v string) {
		b = fmt.Appendf(b, "%s: %s\r\n", n, v)
	     
	})
	b = append(b, "\r\n"...)
	_ , err := w.Write(b) 
	return err 
 }

