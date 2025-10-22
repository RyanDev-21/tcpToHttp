package server

import (
	"fmt"
	"io"
	"net"

	"tests/internal/request"
	"tests/internal/response"
)

type Server struct {
	close   bool
	handler Handler
}

type Handler func(w *response.Writer, req *request.Request)

type HandlerError struct {
	StatusCode response.StatusCode
	Message    string
}

func runServer(s *Server, l net.Listener) {
	for {
		conn, err := l.Accept()
		if s.close {
			return
		}
		if err != nil {
			return
		}
		go handle(s, conn)

	}
}

func handle(s *Server, conn io.ReadWriteCloser) {
	defer conn.Close()
	r, err := request.RequestFromReader(conn)
	responseWriter := response.NewWriter(conn)
	if err != nil {
		responseWriter.WriteStatusLine(response.BadReq)
		responseWriter.WriteHeaders(response.GetDefaultHeaders(0))
		return
	}

	s.handler(responseWriter, r)
}

func Serve(port uint16, handler Handler) (*Server, error) {
	listener, err := listen(port)
	if err != nil {
		return nil, err
	}
	server := &Server{
		close:   false,
		handler: handler,
	}
	go runServer(server, listener)

	return server, nil
}

func listen(port uint16) (net.Listener, error) {
	lisetener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}
	return lisetener, nil
}

func (s *Server) Close() error {
	s.close = true
	return nil
}
