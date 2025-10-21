package request

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strconv"
	H "tests/internal/headers"
)


type parseState string 
const  (
	StartState parseState = "InitRequestLine"
	HeadersState parseState = "InitHeaders"
	BodyState parseState = "InitBody"
	EndState parseState= "End"
	ErrorState parseState = "error"
)


func GetInt(headers *H.Headers,name string,defaultValue int)int{
	values , exists := headers.Get(name)
	if !exists{
		return defaultValue

	}
	value , err:= strconv.Atoi(values)
	if err !=nil{
		return defaultValue	
	}
	return value
}

type Request struct{
	RequestLine RequestLine
	Headers *H.Headers
	Body  []byte 
	State parseState
}


func(r *Request)hasBody()bool{
	length := GetInt(r.Headers,"Content-length",0)
	return length >0
}


//This will parse and return whether we have successfully parsed the requestline or 
//need to be parse more
//and if we have successfully parsed the requestline then we move to the next state which is parsing
//the method
//this parsing and regulating the headers and RequestLine follows the rfc9110 & rfc9112
func (r *Request)parse(data []byte)(int,error){
	read := 0
	outer :

		for {
		currentData := data[read:]
		//after each loop clear the currentData so that it will request the next data
		if len(currentData) == 0{
			break outer
		}
		switch r.State{
		case StartState:
			rl,n,err:= parseRequestLine(currentData)
			if n == 0{
				break outer	
			}
			if err !=nil{
				r.State = ErrorState
				return 0,err 
			}
			read += n
			
			r.RequestLine = *rl
			r.State = HeadersState 

		case HeadersState:
			n,done, err:= r.Headers.Parse(currentData)
			if err !=nil{
				return n,err
			}
			if n==0{
				break outer	
			}
			
			read +=n	
			if done{
				
				if r.hasBody(){
					r.State = BodyState
				}else {
					r.State = EndState
				}
			}
		case BodyState:
			length := GetInt(r.Headers,"Content-length",0)
			if length == 0{
				panic("chunked not implemented")
			}

					

			//this one will read as long as the length is not sastified
			remaining:= min(length-len(r.Body),len(currentData))
			r.Body = append(r.Body,currentData[:remaining]...)
			read += remaining	
			if len(r.Body) == length{
				r.State = EndState
				break outer
			}
		case EndState:
			break outer	
		default: 
			panic("somehow we did some stupid shit")
		}
	}
	return read, nil


}


func (r *Request)checkEndState()bool{
	return r.State == EndState
}

func (r *Request)checkErrorState()bool{
	return r.State == ErrorState
}

type RequestLine struct{
	Method string
	HttpVersion string
	Target   string
}
//Thsese two funtion are just check for some specific format for parsing RequestLine
func (rl *RequestLine) validHttpVersion()bool{
	return rl.HttpVersion == "1.1"
}

func (rl *RequestLine) validMethod()bool{
	valid := true
	if rl.Method == " "{
		valid = false
		return valid
	}		

	for _,v := range rl.Method{
		if v >='A' && v<='Z'{
			valid = true
		}else{
			valid = false
			break
		}
	}
	return valid
}


var ErrFailedToRead  = fmt.Errorf("cannot read all the  passed in reader")
var ErrInvalidRequestLine = fmt.Errorf("invalid request line format")
var ErrFailedToParse = fmt.Errorf("unable to parse the request line")
var SEPARATOR = []byte("\r\n")



//This function parse the RequestLine by searching the SEPARATOR 
//and returns number of bytes read and an error if there is any
func parseRequestLine(data []byte)(*RequestLine,int,error){
	index := bytes.Index(data,SEPARATOR)
	if index == -1{
		return nil,0,nil
	}
	
	startLine := data[:index]
	read := index+len(SEPARATOR)

	

	parts := bytes.Split(startLine,[]byte(" "))
	if len(parts) != 3{
		return nil,read,ErrInvalidRequestLine
	}	

	httpParts := bytes.Split(parts[2],[]byte("/"))
	if len(httpParts) != 2{
		return nil,read,ErrInvalidRequestLine
	}
	
	rl :=RequestLine{
		Method: string(parts[0]),
		HttpVersion: string(httpParts[1]),
		Target: string(parts[1]),
	}	
	
	if !rl.validHttpVersion(){
		return nil,read,ErrInvalidRequestLine
	}

	if !rl.validMethod(){
		return nil,read,ErrInvalidRequestLine
	}	

	return &rl,read,nil
}


func RequestFromReader(reader io.Reader)(*Request,error){
	request := &Request{
		State: StartState,
		Headers: H.NewHeaders(),
	}
	//this will only handle the bytes array which len is <=1kb 
	//any other greater than that can cause problem need to think more efficient way
	data := make([]byte,1024)
	index := 0
	for !request.checkEndState() && !request.checkErrorState(){
		//this one read into the buffer(data)
		n, err := reader.Read(data[index:])
		if err !=nil{
			return nil,errors.Join(
				ErrFailedToParse,
				err,
				)
		}
		index +=n
		//this one read from buffer(data) and interpret it
		//in other words, it parse from buffer(data)
		readN , err := request.parse(data[:index])
		if err !=nil{
			return nil,errors.Join(
				ErrFailedToParse,
				err,
				)
		}
		//update the data if we have read any
		//for example if the parse function returns smth then we know it is has read smth 
		//for example it might have paresed requestline or one of the header
		//if so the readN will have value ,not 0, then we do that what ever left so whcih means no more data left 
		//and if the data left then we just simply update that  data with the current one 
		//if not
		// the data becomes empty and then index become 0 again and then same thing 
		//and after hitting the crlf line at the end the parese funtion will switch the state of the requestline to the endstate
		//then the loop will exist
		copy(data,data[readN:index])
		index -= readN
	}


	
	return request,nil
}
