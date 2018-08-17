package model

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
)

// OrderRequest is used to store data of request coming to orderHandler
type OrderRequest struct {
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

// NewOrderRequest is constructor of OrderRequest which converts request body coming to orderHandler to OrderRequest
func NewOrderRequest(requestBody io.Reader) *OrderRequest {
	s, _ := ioutil.ReadAll(requestBody)
	requestBody = ioutil.NopCloser(bytes.NewBuffer(s))
	decoder := json.NewDecoder(requestBody)
	var t *OrderRequest
	err := decoder.Decode(t)
	if err != nil {
		panic(err)
	}
	return t
}
