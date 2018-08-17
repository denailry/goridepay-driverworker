package reject

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
)

// Request is used to store data of request coming to rejectHandler
type Request struct {
	DriverID int
	OrderID  int
}

// NewRequest is constructor of reject.Request which converts request body coming to rejectHandler to reject.Request
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
