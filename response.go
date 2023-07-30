package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strconv"
)

const (
	protoVersion = "HTTP/1.1"
)

type ResponseWriter struct {
	conn net.Conn

	Headers Headers
}

func (w *ResponseWriter) Header(key, value string) {
	if w.Headers == nil {
		w.Headers = map[string]string{key: value}
		return
	}

	w.Headers[key] = value
}

func (w *ResponseWriter) WriteStatus(status int) error {
	return w.WriteBody(status, nil)
}

func (w *ResponseWriter) WriteJSON(status int, v any) error {
	body, err := json.Marshal(v)
	if err != nil {
		return fmt.Errorf("json: %w", err)
	}

	w.Header("Content-Type", "application/json")
	return w.WriteBody(status, body)
}

func (w *ResponseWriter) WriteText(status int, text string) error {
	w.Header("Content-Type", "text/plain")
	return w.WriteBody(status, []byte(text))
}

func (w *ResponseWriter) WriteBody(status int, body []byte) error {
	buf := &bytes.Buffer{}

	// Status line
	buf.WriteString(protoVersion + " ")
	buf.WriteString(strconv.Itoa(status) + " ")
	buf.WriteString(http.StatusText(status))
	buf.Write(crlf)

	// Header section
	w.addContentLength(body)

	for v, k := range w.Headers {
		buf.WriteString(v + ": " + k)
		buf.Write(crlf)
	}
	buf.Write(crlf)

	if len(body) == 0 {
		_, err := w.Write(buf.Bytes())
		return err
	}

	// Body section
	buf.Write(body)

	_, err := w.Write(buf.Bytes())
	return err
}

func (w *ResponseWriter) Write(b []byte) (int, error) {
	return w.conn.Write(b)
}

func (w *ResponseWriter) addContentLength(body []byte) {
	w.Header("Content-Length", strconv.Itoa(len(body)))
}
