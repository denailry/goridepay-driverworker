package request

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
)

// Reject is used to store data of request coming to rejectHandler
type Reject struct {
	DriverID int
	OrderID  int
}

// NewReject is constructor of request.Reject which converts request body coming to rejectHandler to request.Reject
func NewReject(requestBody io.Reader) *Reject {
	s, _ := ioutil.ReadAll(requestBody)
	requestBody = ioutil.NopCloser(bytes.NewBuffer(s))
	decoder := json.NewDecoder(requestBody)
	var t Reject
	err := decoder.Decode(&t)
	if err != nil {
		panic(err)
	}
	return &t
}
