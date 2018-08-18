package request

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
)

// Accept is used to store data of request coming to acceptHandler
type Accept struct {
	DriverID int
	OrderID  int
}

// NewAccept is constructor of request.Accept which converts request body coming to acceptHandler to request.Accept
func NewAccept(requestBody io.Reader) *Accept {
	s, _ := ioutil.ReadAll(requestBody)
	requestBody = ioutil.NopCloser(bytes.NewBuffer(s))
	decoder := json.NewDecoder(requestBody)
	var t Accept
	err := decoder.Decode(&t)
	if err != nil {
		panic(err)
	}
	return &t
}
