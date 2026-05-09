package request

import (
	"bytes"
	"fmt"
	"io"
)
type Request struct {
	RequestLine RequestLine
	state parserState
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

type parserState string
const (
	StateInit parserState = "init"
	StateDone parserState = "done"
	StateErr parserState = "Error"
)

func (r *RequestLine) ValidHttp() bool {
	return r.HttpVersion == "HTTP/1.1"
}

func newRequest() *Request{
	return &Request{
		state: StateInit,
	}
}
 var ERROR_BAD_START_LINE = fmt.Errorf("malformed request Line!")
 var ERROR_UNSOPPORTED_HTTP= fmt.Errorf("HTTP IS NOT 1.1")
 var ERROR_State=fmt.Errorf("state err")
 var SEPARATOR = []byte("\r\n")


func parserequestLine(b []byte) (*RequestLine, int, error) {
    idx := bytes.Index(b, SEPARATOR)
    if idx == -1 {
        return nil, 0, nil
    }
    
    startLine := b[:idx]
    read := idx + len(SEPARATOR)
    
    // Split into 3 parts: [GET] [/] [HTTP/1.1]
    parts := bytes.Split(startLine, []byte(" "))
    if len(parts) != 3 {	
        return nil, 0, ERROR_BAD_START_LINE
    }

    rl := &RequestLine{
        Method:        string(parts[0]),
        RequestTarget: string(parts[1]),
        HttpVersion:   string(parts[2]), // Changed from parts[3] to parts[2]
    }

    if !rl.ValidHttp() {
        return nil, 0, ERROR_UNSOPPORTED_HTTP
    }

    return rl, read, nil
}



func (r *Request) parse(data []byte) (int,error) {

	read := 0
outer :
	for{

	 switch r.state{
	 case StateErr:
		return 0,ERROR_State
	 case StateInit:
		rl , n , err := parserequestLine(data[read:])
		if err != nil {
			r.state=StateErr
			return 0,err
		}
		if n== 0 {
			break outer
		}
		r.RequestLine = *rl
		read += n

		r.state = StateDone


	 case StateDone:
		 break outer
	}
  }
  return read,nil
}

func (r *Request) done() bool {
	return r.state == StateDone
}

func (r *Request) Error() bool {
	return r.state == StateErr
}

func RequestFromReader(reader io.Reader) (*Request, error){

	request := newRequest()
	// buf could be overrun
	buf := make([]byte, 1024)
	bufLen:=0
	for !request.done() && !request.Error(){
		n , err := reader.Read(buf[bufLen:])
		if err != nil {
			return nil ,err
		}

		bufLen += n

		readN, err := request.parse(buf[:bufLen + n])
		if err != nil {
			return nil , err
		}

		copy(buf,buf[readN:bufLen])
		bufLen -= readN
		
		
	}

	return  request ,nil
	
}