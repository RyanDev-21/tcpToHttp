package server

import (
	"fmt"
	"io"
	"net"
	"tests/internal/request"
	"tests/internal/response"
)
type Server struct{
	close bool
}

type Handler func(w io.Writer,req *request.Request) *HandlerError

type HandlerError struct{
	statusCode int
	message string
}

func runServer(s *Server,l net.Listener){
	for {
		conn, err:= l.Accept()
		if s.close {
			return
		}
		if err !=nil{
			return 
		}
		go handle(conn)	
		
	}
}

func handle(conn io.ReadWriteCloser){
	defer conn.Close()
	response.WriteStatusLine(conn,response.OK)	
	response.WriteHeaders(conn,response.GetDefaultHeaders(0))
	conn.Write([]byte("\r\n"))
}


func Serve(port uint16,handler Handler)(*Server,error){
	listener, err:= listen(port)	
	if err !=nil{
		
		return nil,err
	}
	server := &Server{close: false}	
	go runServer(server,listener)


	return server,nil
}


func listen(port uint16)(net.Listener,error){
	lisetener , err:= net.Listen("tcp",fmt.Sprintf(":%d",port))
	if err !=nil{
		return nil,err
	}
	return lisetener,nil
}


func (s *Server)Close()error{
	s.close = true
	return nil
}
