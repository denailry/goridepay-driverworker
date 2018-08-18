package request

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
)

// Order is used to store data of request coming to orderHandler
type Order struct {
	OrderID             int
	Origin              string
	Destination         string
	DestinationDistance int
	TransactionID       int
	DriverData          []driverData
}

type driverData struct {
	DriverID       int
	OriginDistance int
}

// NewOrder is constructor of request.Order which converts request body coming to orderHandler to request.Order
func NewOrder(requestBody io.Reader) *Order {
	s, _ := ioutil.ReadAll(requestBody)
	requestBody = ioutil.NopCloser(bytes.NewBuffer(s))
	decoder := json.NewDecoder(requestBody)
	var t *Order
	err := decoder.Decode(t)
	if err != nil {
		panic(err)
	}
	return t
}
