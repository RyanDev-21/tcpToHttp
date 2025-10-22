package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"tests/internal/request"
	"tests/internal/response"
	"tests/internal/server"
)

const port= 42069

func BadReq400()[]byte{
	return []byte(
`<html>
<head>
<title>400 Bad Request</title>
</head>
<body>
<h1>Bad Request</h1>
<p>Your request honestly kinda sucked.</p>
</body>
</html>`)
}

func ServerErr500()[]byte{
	return []byte(`<html>
  <head>
    <title>500 Internal Server Error</title>
  </head>
  <body>
    <h1>Internal Server Error</h1>
    <p>Okay, you know what? This one is on me.</p>
  </body>
</html>`)
}

func OK200()[]byte{
	return []byte(
`<html>
  <head>
    <title>200 OK</title>
  </head>
  <body>
    <h1>Success!</h1>
    <p>Your request was an absolute banger.</p>
  </body>
</html>`)
	
}

func main(){
	
	server , err:= server.Serve(port,func(w *response.Writer, req *request.Request)  {
		headers := response.GetDefaultHeaders(0)
		var statusCode response.StatusCode
		body := []byte{}
		switch req.RequestLine.Target{
		case "/yourproblem":
			statusCode =  response.BadReq
			body  = BadReq400()
		case "myproblem":
			statusCode = response.ServerError
			body = ServerErr500()
		default : 
			statusCode = response.OK
			body = OK200()
		}
		headers.Set("Content-length",fmt.Sprintf("%d",len(body)))	
		headers.Set("Content-Type","text/html")	
		w.WriteStatusLine(statusCode)
		w.WriteHeaders(headers)
		w.WriteBody(body)
		
			
	})
     
	if err !=nil{
		log.Fatalf("Error starting server :%v",err)
	}

	defer server.Close()
	log.Println("Server has started on port",port)

	sigChan := make(chan os.Signal,1)
	signal.Notify(sigChan,syscall.SIGINT,syscall.SIGTERM)
	<-sigChan
	log.Println("Sever gracefully stopped")
}


