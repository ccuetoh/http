package http

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
)

type Handler interface {
	Handle(r *Request, w *ResponseWriter)
}

type Server struct {
	Handler Handler
	Address string
	Port    int

	socket *socket
}

func New(opts ...Option) (*Server, error) {
	s := &Server{}
	_ = DefaultOption()(s)

	for _, opt := range opts {
		err := opt(s)
		if err != nil {
			return nil, err
		}
	}

	return s, nil
}

func (s *Server) ListenAndServe() error {
	if s.Handler == nil {
		return errors.New("no handler set")
	}

	sock, err := listen(s.Address, s.Port)
	if err != nil {
		return fmt.Errorf("socket: %w", err)
	}

	s.socket = sock
	for {
		conn, err := sock.Accept()
		if err != nil {
			log.Printf("invalid connection: %v\n", err)
			continue
		}

		go s.Serve(conn)
	}
}

func (s *Server) Close() error {
	if s.socket == nil {
		return errors.New("server is not running")
	}

	return s.socket.Close()
}

func (s *Server) Serve(conn net.Conn) {
	w := &ResponseWriter{conn: conn}
	for {
		r, err := s.socket.receive(conn)
		if err != nil {
			if errors.Is(err, io.EOF) || errors.Is(err, net.ErrClosed) {
				// No need to do anything
				return
			}

			log.Println(err)

			_ = w.WriteStatus(http.StatusBadRequest)
			_ = w.conn.Close()
			return
		}

		log.Println(r)

		w.Header("Server", "ccueto-http")
		w.Header("Connection", "keep-alive")

		s.Handler.Handle(r, w)
	}
}
