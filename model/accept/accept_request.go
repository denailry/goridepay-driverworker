package accept

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
)

// Request is used to store data of request coming to acceptHandler
type Request struct {
	DriverID int
	OrderID  int
}

// NewRequest is constructor of accept.Request which converts request body coming to acceptHandler to accept.Request
func NewRequest(requestBody io.Reader) *Request {
	s, _ := ioutil.ReadAll(requestBody)
	requestBody = ioutil.NopCloser(bytes.NewBuffer(s))
	decoder := json.NewDecoder(requestBody)
	var t *Request
	err := decoder.Decode(t)
	if err != nil {
		panic(err)
	}
	return t
}
