package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

func main(){
	addr,err := net.ResolveUDPAddr("udp","localhost:42069")

	if err !=nil{
		log.Fatal("error","error",err)
	}

	conn, err := net.DialUDP("udp",nil,addr)
	if err !=nil{
		log.Fatal("Error establishing the udp connection")
	}

	defer conn.Close()
	

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print(">")
		lines,err:= reader.ReadString('\n')
		if err !=nil{
			log.Println("error reading input:",err)
			continue
		}	
	
		_,err= conn.Write([]byte(lines))
		if err !=nil{
			log.Println("error reading data: ",err)
			continue	
		}



	}





}
