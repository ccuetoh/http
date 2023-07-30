package http

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"net"
	"strconv"
)

const (
	tcp = "tcp"
)

const (
	cr = 0xD
	lf = 0xA
)

var crlf = []byte{cr, lf}

type socket struct {
	net.Listener
}

func (s socket) receive(conn net.Conn) (*Request, error) {
	r := &Request{}
	reader := bufio.NewReader(conn)

	// RFC9112 2.1
	// HTTP-message   = start-line CRLF
	//                  *( field-line CRLF )
	//                  CRLF
	//                  [ message-body ]

	// Request line (for example, "POST / HTTP/1.1")
	//
	requestLine, err := reader.ReadBytes('\n')
	if err != nil {
		return nil, fmt.Errorf("start: %w", err)
	}

	err = r.parseRequestLine(requestLine)
	if err != nil {
		return nil, fmt.Errorf("request: %w", err)
	}

	// Header section
	// field-line CRLF
	for {
		fieldLine, err := reader.ReadBytes('\n')
		if err != nil {
			return nil, fmt.Errorf("headers: %w", err)
		}

		if bytes.Equal(fieldLine, crlf) {
			break
		}

		err = r.parseHeader(fieldLine)
		if err != nil {
			return nil, fmt.Errorf("headers: %w", err)
		}
	}

	bodyLenStr, contains := r.Headers["Content-Length"]
	if !contains {
		// No body expected
		return r, nil
	}

	bodyLen, err := strconv.Atoi(bodyLenStr)
	if err != nil {
		return nil, fmt.Errorf("content length: %w", err)
	}

	// Body section
	body := make([]byte, bodyLen)
	n, err := reader.Read(body)
	if err != nil {
		return nil, fmt.Errorf("start: %w", err)
	}

	if n != bodyLen {
		return nil, errors.New("body is shorter than expected")
	}

	r.Body = body
	return r, nil
}

func listen(host string, port int) (*socket, error) {
	addr := net.JoinHostPort(host, strconv.Itoa(port))
	lt, err := net.Listen(tcp, addr)
	if err != nil {
		return nil, err
	}

	return &socket{lt}, nil
}
