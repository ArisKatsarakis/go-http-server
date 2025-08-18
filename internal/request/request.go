package request

import (
	"bytes"
	"fmt"
	"io"
	"strings"
)

type parseState string

const (
	ParseInit parseState = "init"
	ParseDone parseState = "done"
)

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

type Request struct {
	RequestLine RequestLine
	ParseState  parseState
}

func (r *Request) parse(b []byte) (int, error) {
	read := 0
outer:
	for {

		switch r.ParseState {
		case ParseInit:
		case ParseDone:
			break outer
		}
	}
	return read, nil
}

func (r *Request) done() bool {
	return r.ParseState == ParseDone
}
func newRequest() *Request {
	return &Request{
		ParseState: ParseInit,
	}
}

var ERR_BAD_REQUEST = fmt.Errorf("this an error line")
var SEPERATOR = []byte("\r\n")

func parseRequestLine(b []byte) (*RequestLine, int, error) {
	idx := bytes.Index(b, SEPERATOR)
	if idx == -1 {
		return nil, 0, nil
	}

	starLine := b[:idx]
	read := idx + len(SEPERATOR)
	parts := bytes.Split(starLine, []byte(" "))
	if len(parts) != 3 {
		return nil, read, ERR_BAD_REQUEST
	}

	return &RequestLine{
		Method:        string(parts[0]),
		RequestTarget: string(parts[1]),
		HttpVersion:   string(parts[2][len(parts[2])-3:]),
	}, read, nil

}

func RequestFromReader(reader io.Reader) (*Request, error) {
	request := newRequest()

	buf := make([]byte, 1024)
	bufIdx := 0
	for !request.done() {
		n, err := reader.Read(buf[bufIdx:])
		if err != nil {
			return nil, err
		}

		readN, err := request.parse(buf[:bufIdx+n])
		if err != nil {
			return nil, err
		}

		bufIdx += n
		copy(buf, buf[readN:bufIdx])
		bufIdx -= readN

	}
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
