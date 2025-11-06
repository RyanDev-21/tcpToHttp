package response

import (
	"crypto/sha256"
	"fmt"
	"io"
	"log"
	"tests/internal/headers"
)

type StatusCode int 

const (
	OK StatusCode = 200
	BadReq StatusCode = 400
	ServerError StatusCode = 500
)



var  ErrNotFoundStatusCode = fmt.Errorf("cannot find the status code")


func GetDefaultHeaders(contentLen int)headers.Headers{
	headers := headers.NewHeaders()
	headers.Set("Content-Length",fmt.Sprintf("%d",contentLen))
	headers.Set("Connection","close")
	headers.Set("Content-Type","text/plain")

	return *headers
}

type Writer struct{ writer io.Writer }

func NewWriter(writer io.Writer)*Writer{
	return &Writer{writer: writer}
}


func (w *Writer)WriteChunkedBody(buff []byte){
	w.WriteBody([]byte(fmt.Sprintf("%x\r\n",len(buff))))
	w.WriteBody(buff)
	w.WriteBody([]byte("\r\n"))
}

func (w *Writer)WriteChunkedBodyDone(body []byte){
	w.WriteBody([]byte("0\r\n"))
	trailers := headers.NewHeaders()
	trailers.Set("X-Content-SHA256",fmt.Sprintf("%02x",sha256.Sum224(body)))
	trailers.Set("X-Content-Length",fmt.Sprintf("%v",len(body)))
	w.WriteBody([]byte("\r\n"))
}


func (w *Writer)WriteStatusLine(statusCode StatusCode)error{
	var b []byte
	switch statusCode{
		case OK:
			b = fmt.Appendf(b,"HTTP/1.1 %d OK\r\n",statusCode)
		case BadReq:
			b = fmt.Appendf(b,"HTTP/1.1 %d Bad Request\r\n",statusCode)
		case ServerError:
			b = fmt.Appendf(b,"HTTP/1.1 %d Internal Server Error\r\n",statusCode)
		default:
			return ErrNotFoundStatusCode
	}
	log.Printf("StatusLine /%v/",string(b))
	_,err:=w.writer.Write(b)
	return err 
	
}

func (w *Writer)WriteTrailers(h headers.Headers)error{
	h.Delete("Content-length")
	h.Set("Trailer","X-Content-SHA256")
	h.Set("Trailer","X-Content-Length")
	
	return nil
}

func (w *Writer)WriteHeaders(h headers.Headers) error{
	var err error = nil
	var b = []byte{}
	h.ForEach(func(k,v string){
		if err !=nil{
			return
		}	
		b=fmt.Appendf(b,"%s: %s\r\n",k,v)
	})
	b=fmt.Append(b,"\r\n")
	log.Printf("Headers /%v/",string(b))
	_,err= w.writer.Write(b)
	return err
}

func (w *Writer)WriteBody(p []byte)(int,error){
	n,err := w.writer.Write(p)
	return n,err
}


