package request

import (
	"bytes"
	"fmt"
	"httpfromtcp/internal/headers"
	"io"
	"strconv"
)

type Request struct {
	RequestLine RequestLine
	Headers *headers.Headers
	Body string
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
	StateBody parserState = "Body"
	StateHeaders parserState = "Headers"
	StateErr parserState = "Error"
)

func (r *RequestLine) ValidHttp() bool {
	return r.HttpVersion == "HTTP/1.1"
}

func getInt(headers *headers.Headers ,name string, defaultValue int) int {
	valueStr , exists := headers.Get(name)
	if !exists {
		return defaultValue
	}

	value , err := strconv.Atoi(valueStr) 
	if err != nil {
		return defaultValue
	}
	return value
}

func newRequest() *Request{
	return &Request{
		state: StateInit,
		Headers: headers.NewHeaders(),
		Body: "",
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
        HttpVersion:   string(parts[2]), 
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
		currentData := data[read:]
	 switch r.state{
	 case StateErr:
		return 0,ERROR_State

	 case StateHeaders:
		n , done , err := r.Headers.Parse(currentData)
		if err != nil {
			r.state=StateErr
			return 0,err
		}
		
		if n== 0 {
			break outer
		}

		read += n
		if done{
			r.state = StateBody
		}

	 case StateInit:
		rl , n , err := parserequestLine(currentData)
		if err != nil {
			r.state=StateErr
			return 0,err
		}
		if n== 0 {
			break outer
		}
		r.RequestLine = *rl
		read += n
		r.state = StateHeaders

	case StateBody:
	length := getInt(r.Headers, "content-length", 0)
	if length == 0 {
		r.state = StateDone
		break outer
	}

	toRead := length - len(r.Body)
	if toRead <= 0 {
		r.state = StateDone
		break
	}

	if len(currentData) == 0 {
		break outer 
	}

	remaining := min(toRead, len(currentData))

	r.Body += string(currentData[:remaining])
	read += remaining

	if len(r.Body) == length {
		r.state = StateDone
	}

	 case StateDone:
		 break outer
	 default:
	    panic("Somhw we fkd")
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

func RequestFromReader(reader io.Reader) (*Request, error) {
    request := newRequest()
    buf := make([]byte, 1024)
    bufLen := 0

    for !request.done() && !request.Error() {
        n, err := reader.Read(buf[bufLen:])
        bufLen += n // Update buffer length immediately

		if n == 0 && err == nil {
        return nil, fmt.Errorf("no progress from reader")
    }

  
        if bufLen > 0 {
            readN, parseErr := request.parse(buf[:bufLen])
            if parseErr != nil {
                return nil, parseErr
            }
            if readN > 0 {
                copy(buf, buf[readN:bufLen])
                bufLen -= readN
            }
        }

    
        if err != nil {
            if err == io.EOF {
                if !request.done() {
                    return nil, io.ErrUnexpectedEOF
                }
                break
            }
            return nil, err
        }

        if bufLen == len(buf) {
            return nil, fmt.Errorf("request header or body chunk too large")
        }
    }

    if !request.done() {
        return nil, fmt.Errorf("incomplete request")
    }

    return request, nil
}

// 2:20:00