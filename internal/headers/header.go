package headers

import (
	"bytes"
	"fmt"
)


type Headers map[string]string

var rn=[]byte("\r\n")

func NewHeaders() Headers {
	return make(Headers)
}

func parseHeader(fieldLine []byte) (string, string, error) {
  
    if len(fieldLine) > 0 && (fieldLine[0] == ' ' || fieldLine[0] == '\t') {
        return "", "", fmt.Errorf("malformed header: leading whitespace")
    }

    parts := bytes.SplitN(fieldLine, []byte(":"), 2)
    if len(parts) != 2 {
        return "", "", fmt.Errorf("malformed header: missing colon")
    }

    name := parts[0]
    value := bytes.TrimSpace(parts[1])

    
    if bytes.HasSuffix(name, []byte(" ")) || bytes.HasSuffix(name, []byte("\t")) {
        return "", "", fmt.Errorf("malformed header: whitespace before colon")
    }

    return string(name), string(value), nil
}

func (h Headers) Parse(data []byte) (n int, done bool, err error){

	read := 0 
	done = false	
	for{
		idx := bytes.Index(data[read:],rn)
		if idx == -1 {
			return read, false, nil
		}

		if idx == 0 {
			read += len(rn)
			return read, true, nil
		}

		name , value , err := parseHeader(data[read:read+idx])
		if err != nil {
			return 0,false,err
		}

		h[name] = value 
		read += idx + len(rn)
		


	}

}