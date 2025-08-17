package request

import (
	"fmt"
	"io"
	"strings"
)

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

type Request struct {
	RequestLine RequestLine
}

var ERR_BAD_REQUEST = fmt.Errorf("this an error line")
var SEPERATOR = "\r\n"

func parseRequestLine(b string) (*RequestLine, string, error) {
	idx := strings.Index(b, SEPERATOR)
	if idx == -1 {
		return nil, b, nil
	}

	starLine := b[:idx]
	restOfMsg := b[idx+len(SEPERATOR):]
	parts := strings.Split(starLine, " ")
	if len(parts) != 3 {
		return nil, restOfMsg, ERR_BAD_REQUEST
	}

	return &RequestLine{
		Method:        parts[0],
		RequestTarget: parts[1],
		HttpVersion:   parts[2][len(parts[2])-3:],
	}, restOfMsg, nil

}

func RequestFromReader(reader io.Reader) (*Request, error) {
	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	str := string(data)
	rl, str, err := parseRequestLine(str)
	if err != nil {
		return nil, err
	}
	r := &Request{
		RequestLine: *rl,
	}
	return r, nil
}
