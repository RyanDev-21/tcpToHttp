package headers

import (
	"bytes"
	"fmt"
	"strings"
)


var SEPARATOR =[]byte("\r\n") 
var colon = []byte(":")
var ErrNotValidHeader = fmt.Errorf("not valid header")
var ErrNotValidToken = fmt.Errorf("not valid token")

type Headers struct{
	headers	map[string]string
}
func NewHeaders()*Headers{
	return &Headers{headers:map[string]string{}}
}


func (h Headers)Get(name string)(string,bool){
	str , ok := h.headers[strings.ToLower(name)]
	return str,ok	
}

func (h Headers) Set(name ,value string){
	h.headers[strings.ToLower(name)] = value
}

func (h Headers)ForEach(cb func(k,v string)){
	for k,v:= range h.headers{
		cb(k,v)
	}
}


func (h Headers)combineHeadersValue(header,value string)string{
	oldValue,_ := h.Get(header)
	return fmt.Sprintf("%s,%s",oldValue,value)	
}


//this apply the constraints based on the token rule of the rfc
func (h Headers)isValidToken(token string)bool{
	valid := false 
	for _, ch := range token{
		if ch >= 'A' &&  ch <='Z' ||
			ch >= 'a' && ch<='z'||
			ch >='0'&& ch<='9'{
				valid = true	
		}
		switch ch{

		case '!', '#', '$', '%', '&', '\'', '*', '+','-','.', '^', '_', '`', '|', '~':
			valid = true	

		}
		if !valid {
			return false
		}
	}

	return valid
}


//this returns the header and value if the fieldLine and value are valid 
func parseHeader(fieldLine []byte)(string,string,error){
	parts := bytes.SplitN(fieldLine,colon,2)
	if len(parts) != 2{
		return "","",nil
	}
	name := parts[0]
	if bytes.HasSuffix(name,[]byte(" ")){
		return "","",ErrNotValidHeader
	}

	name = bytes.TrimLeft(name," ")
	value := bytes.TrimSpace(parts[1])


	return string(name),string(value),nil
}



func (h Headers)Parse(data []byte)(int,bool,error){
	read := 0
	done := false
	for {
		index := bytes.Index(data[read:],SEPARATOR)
		if index == -1{
			return read,done,nil
		}

		//End of the crlf
		if index ==0{
			read += len(SEPARATOR)
			done = true
			break		
		}
		header, value ,err := parseHeader(data[read:index+read])
		if err !=nil{
			return 0,done,err
		}
		if !h.isValidToken(header){
			return 0,done,ErrNotValidToken
		}
		
		_,ok := h.Get(header)
		if ok{
			value =h.combineHeadersValue(header,value)	
		}
		h.Set(header,value)	
		read +=index+len(SEPARATOR)
	}

	return read,done,nil

}

