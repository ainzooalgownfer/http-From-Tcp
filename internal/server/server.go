package server

import (
	"context"
	"fmt"
	"httpfromtcp/internal/request"
	"httpfromtcp/internal/response"
	"io"
	"net"
	"sync"
)

type Server struct{
	listener net.Listener  
	closed bool
	handler Handler
	activeConns int64        
    wg          sync.WaitGroup
    closeOnce   sync.Once

}

type HandlerError struct {
	StatusCode response.StatusCode
	Message string 
}

type Handler func(w *response.Writer, req *request.Request) 


func runConnection(s *Server, conn io.ReadWriteCloser) {
		defer conn.Close() 

		responseWriter := response.NewWriter(conn)
		r , err := request.RequestFromReader(conn)

		if err != nil {
			responseWriter.WriteStatusLine(response.StatusBadRequest)
			responseWriter.WriteHeaders(*response.GetDefaultHeaders(0))
			return
		}

		 s.handler(responseWriter, r)

	}


func runServer(s *Server, listener net.Listener) {
	for {
		conn, err := listener.Accept()
		if err != nil {
			if s.closed {
				return
			}
			continue
		}
	
		s.wg.Add(1)
		go func(c net.Conn) {
			defer s.wg.Done()
			runConnection(s, c)
		}(conn)
	}
}

func Serve(port uint16, handler Handler) (*Server, error) {
    listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
    if err != nil {
        return nil, err
    }
    s := &Server{
        listener: listener,
        handler:  handler,
		closed: false,
    }
    go runServer(s, listener)
    return s, nil
}


func (s *Server) Shutdown(ctx context.Context) error {
	if err := s.Close(); err != nil {
		return err
	}
	done := make(chan struct{})
	go func() {
		s.wg.Wait()
		close(done)
	}()
	select {
	case <-done:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (s *Server) Close() error {
	var err error
	s.closeOnce.Do(func() {
		s.closed = true
		if s.listener != nil {
			err = s.listener.Close()
		}
	})
	return err
}