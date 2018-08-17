package order

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
)

// Request is used to store data of request coming to orderHandler
type Request struct {
	OrderID     int
	Origin      string
	Destination string
	DriverData  []driverData
}

type driverData struct {
	DriverID            int
	OriginDistance      int
	DestinationDistance int
}

// NewRequest is constructor of order.Request which converts request body coming to orderHandler to order.Request
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
