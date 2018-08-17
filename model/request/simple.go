package request

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
)

// Simple is the most simple request structure
type Simple struct {
	DriverID int
}

// NewSimple is constructor of request.Simple which converts request body to equest.Simple
func NewSimple(requestBody io.Reader) *Simple {
	s, _ := ioutil.ReadAll(requestBody)
	requestBody = ioutil.NopCloser(bytes.NewBuffer(s))
	decoder := json.NewDecoder(requestBody)
	var t *Simple
	err := decoder.Decode(t)
	if err != nil {
		panic(err)
	}
	return t
}
