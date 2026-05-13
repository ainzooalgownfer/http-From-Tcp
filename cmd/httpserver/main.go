package main

import (
	"io"
	"log"
	"os"

	"httpfromtcp/internal/request"
	"httpfromtcp/internal/response"
	"httpfromtcp/internal/server"
	"os/signal"
	"syscall"
)


const port = 42069

func main() {
	s, err := server.Serve(port, func(w io.Writer, req *request.Request) *server.HandlerError {
			switch req.RequestLine.RequestTarget {
							case "/yourproblem":
						return &server.HandlerError{
							StatusCode: response.StatusBadRequest,
							Message: "ur pb not mine\r\n" ,
						}
			case "/myproblem":
						return &server.HandlerError{
							StatusCode: response.StatusInternalServerError,
							Message: "Woopsie, my bad \r\n" ,
						}
			default:
				w.Write([]byte("ALL goo fr fr \n"))
			}
		return nil
	
		} )
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer s.Close()
	log.Println("Server started on port", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server gracefully stopped")
}