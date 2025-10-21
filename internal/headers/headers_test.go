package headers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)


func TestHeaders(t *testing.T){
	
	//Test: Valid single header
	h := *NewHeaders()
	data := []byte("Host: localhost:42069\r\nFoo: foofoo\r\n\r\n")
	n , done , err := h.Parse(data)
	require.NoError(t,err)
	require.NotNil(t,h.headers)
	assert.Equal(t,"localhost:42069",h.Get("Host"))
	assert.Equal(t,"foofoo",h.Get("Foo"))
	assert.Equal(t,"",h.Get("missing one"))
	assert.Equal(t,36,n)
	assert.True(t,done)

	//Test: Invalid  header space
	h =* NewHeaders()
	data = []byte("     Host : localhost:42069      \r\n\r\n")
	n, done, err = h.Parse(data)
	require.Error(t,err)
	assert.Equal(t,0,n)
	assert.False(t,done)

	//Test: With two valid header
	h= *NewHeaders()
	data = []byte("Host: localhost:42069\r\nUser-Agent: curl/7.53.0\r\n\r\n")
	_, done ,err= h.Parse(data)

	require.NoError(t,err)
	require.NotNil(t,h.headers)
	assert.Equal(t,"localhost:42069",h.Get("Host"))
	assert.Equal(t,"curl/7.53.0",h.Get("User-Agent"))
	assert.Equal(t,"",h.Get("missing one"))
	assert.True(t,done)

	//Test: Invalid  token 
	h =* NewHeaders()
	data = []byte("     H©st : localhost:42069      \r\n\r\n")
	n, done, err = h.Parse(data)
	require.Error(t,err)
	assert.Equal(t,0,n)
	assert.False(t,done)

	//Test: With multiple value for same field-line
	h= *NewHeaders()
	data = []byte("Host: localhost:42069\r\nHost: curl/7.53.0\r\n\r\n")
	_, done ,err= h.Parse(data)

	require.NoError(t,err)
	require.NotNil(t,h.headers)
	assert.Equal(t,"localhost:42069,curl/7.53.0",h.Get("Host"))
	assert.Equal(t,"",h.Get("missing one"))
	assert.True(t,done)
	
		
}

