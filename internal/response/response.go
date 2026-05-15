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




 func GetDefaultHeaders(contentLen int) *headers.Headers {
	h := headers.NewHeaders()
	h.Set("Content-Length", fmt.Sprintf("%d", contentLen))
	h.Set("Connection", "closed")
	h.Set("Content-Type", "text/plain")

	return h
 }

 

 type Writer struct {
	writer io.Writer
	headersStore *headers.Headers
	wroteHeader bool
	
}

func NewWriter(writer io.Writer)  *Writer {
	return  &Writer{writer: writer,
					headersStore: headers.NewHeaders(),
				}
}

    func (w *Writer) WriteStatusLine(statusCode StatusCode) error {
	statusLine := []byte{}
	switch statusCode {
	case StatusOk:  statusLine = []byte("HTTP/1.1 200 OK\r\n")
	case StatusBadRequest: statusLine = []byte("HTTP/1.1 400 BAD REQUEST\r\n")
	case StatusInternalServerError: statusLine= []byte("HTTP/1.1 500 INTERNAL SERVER ERROR\r\n")
	default: return fmt.Errorf("unowkn error code")
	}
	
	_, err := w.writer.Write(statusLine)
	return err
	}

func (w *Writer) SetHeader(key, value string) {
    if w.wroteHeader {
        return
    }

	if w.headersStore == nil {
        w.headersStore = headers.NewHeaders()
    }
	
    w.headersStore.Set(key, value)
}
	
	func (w *Writer) WriteHeaders(h headers.Headers) error{
		if w.wroteHeader {
        	return nil
    	}

		 if w.headersStore != nil {
      		  w.headersStore.Foreach(func(k, v string) {
            h.Set(k, v)
       		 })
    	}

		b := []byte{}
	    h.Foreach(func(n, v string) {
		b = fmt.Appendf(b, "%s: %s\r\n", n, v)
	     
	      })
	     b = append(b, "\r\n"...)
	     _ , err := w.writer.Write(b)

		  if err == nil {
        w.wroteHeader = true   
   		}
	    return err	
	}


	func (w *Writer) WriteBody(p []byte) (int, error){
		n , err := w.writer.Write(p)
		if err != nil {
			return 0 , err
		}
		return n , err
	}
