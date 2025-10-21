package main

import (
	"fmt"
	"log"
	"net"

	"tests/internal/request"
)

// we pass the conn object whcih is the same as io Reader interface
func main() {
	listener, err := net.Listen("tcp", ":42069")
	if err != nil {
		log.Fatal(err)
	}
	for {

		conn, err := listener.Accept()
		if err != nil {
			log.Fatal("error", "error", err)
		}
		r, err := request.RequestFromReader(conn)
		if err != nil {
			log.Fatal("error", "error", err)
		}

		fmt.Printf("Request Line :\n- Method: %v\n- Target: %v\n- Version: %v\n", r.RequestLine.Method, r.RequestLine.Target, r.RequestLine.HttpVersion)
		fmt.Println("Headers: ")
		r.Headers.ForEach(func(k,v string){
			fmt.Printf("- %v:%v\n",k,v)	
		})
		fmt.Println("Body :")
		fmt.Printf("BODY_STRING\n%s",string(r.Body))

	}
}
