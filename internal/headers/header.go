package headers

import (
	"bytes"
	"fmt"
	"strings"
)


type Headers struct {
	headers map[string]string
}

func isToken(str []byte) bool {
    if len(str) == 0 {
        return false
    }

    for _, ch := range str {
        switch {
      
        case ch >= 'a' && ch <= 'z':
        case ch >= 'A' && ch <= 'Z':
        case ch >= '0' && ch <= '9':
        
        case ch == '!' || ch == '#' || ch == '$' || ch == '%' || ch == '&' || 
             ch == '\'' || ch == '*' || ch == '+' || ch == '-' || ch == '.' || 
             ch == '^' || ch == '_' || ch == '`' || ch == '|' || ch == '~':
       
        default:
            return false
        }
    }
    return true
}

var rn=[]byte("\r\n")


func NewHeaders() *Headers {
	return &Headers {
		headers : map[string]string{},
	}
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

func (h *Headers) Get(name string) (string,bool) {
	
	 str,ok := h.headers[strings.ToLower(name)]
	 return str ,ok
}

func (h *Headers) Replace(name , value string) {
		name = strings.ToLower(name)
		h.headers[name] = value
} 

func (h *Headers) Set(name , value string) {
		name = strings.ToLower(name)

		if v ,ok := h.headers[name]; ok {
			h.headers[name] = fmt.Sprintf("%s,%s",v,value)
		} else  {
				h.headers[name] = value

		}
	
} 

func (h *Headers) Foreach(cb func(n,v string)) {
	for n, v := range h.headers{
		cb(n,v)
	}
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

		if !isToken([]byte(name)){
			return 0,false,fmt.Errorf("malformed header name")
		}


		
		read += idx + len(rn)
		h.Set(name,value) 


	}

}

// 1:52:00