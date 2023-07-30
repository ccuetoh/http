package http

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
)

type Request struct {
	Method       string
	URL          string
	ProtoVersion string
	Headers      Headers
	Body         []byte
}

type Headers map[string]string

func (r *Request) String() string {
	return r.Method + " " + r.URL
}

func (r *Request) parseRequestLine(requestLine []byte) error {
	sections := bytes.Split(requestLine, []byte(" "))
	if len(sections) != 3 {
		return errors.New("invalid section count")
	}

	r.Method = string(sections[0])
	r.URL = string(sections[1])
	r.ProtoVersion = string(sections[2])
	return nil
}

func (r *Request) parseHeader(fieldLine []byte) error {
	fieldLine = bytes.TrimSuffix(fieldLine, crlf)

	// field-line = field-name ":" OWS field-value OWS
	sections := bytes.SplitN(fieldLine, []byte(":"), 2)
	if len(sections) != 2 {
		return fmt.Errorf("invalid length")
	}

	if r.Headers == nil {
		r.Headers = map[string]string{
			string(sections[0]): strings.TrimPrefix(string(sections[1]), " "),
		}

		return nil
	}

	r.Headers[string(sections[0])] = strings.TrimPrefix(string(sections[1]), " ")
	return nil
}
