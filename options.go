package http

const (
	defaultPort    = 80
	defaultAddress = ""
)

type Option func(s *Server) error

func DefaultOption() Option {
	return func(s *Server) error {
		s.Address = defaultAddress
		s.Port = defaultPort
		return nil
	}
}

func WithAddress(addr string) Option {
	return func(s *Server) error {
		s.Address = addr
		return nil
	}
}

func WithPort(port int) Option {
	return func(s *Server) error {
		s.Port = port
		return nil
	}
}

func WithHandler(h Handler) Option {
	return func(s *Server) error {
		s.Handler = h
		return nil
	}
}
