package request

import (
	"bytes"
	"fmt"
	"io"
)

type parseState string

const (
	ParseInit  parseState = "init"
	ParseDone  parseState = "done"
	ParseError parseState = "error"
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
		case ParseError:
			return 0, fmt.Errorf("Error at request parsing")
		case ParseInit:
			rl, n, err := parseRequestLine(b[read:])

			if err != nil {
				r.ParseState = ParseError
				return 0, err
			}

			if n == 0 {
				break outer
			}

			r.RequestLine = *rl
			read += n

			r.ParseState = ParseDone

		case ParseDone:
			break outer
		}
	}
	return read, nil
}

func (r *Request) done() bool {
	return r.ParseState == ParseDone || r.ParseState == ParseError
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

	return request, nil
}
