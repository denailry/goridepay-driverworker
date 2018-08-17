package invalidate

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
)

// Request is used to store data of request coming to invalidateHandler
type Request struct {
	OrderID int
}

// NewRequest is constructor of invalidate.Request which converts request body coming to invalidateHandler to invalidate.Request
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
