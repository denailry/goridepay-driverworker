package request

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
)

// Invalidate is used to store data of request coming to invalidateHandler
type Invalidate struct {
	OrderID int
}

// NewInvalidate is constructor of request.Invalidate which converts request body coming to invalidateHandler to request.Invalidate
func NewInvalidate(requestBody io.Reader) *Invalidate {
	s, _ := ioutil.ReadAll(requestBody)
	requestBody = ioutil.NopCloser(bytes.NewBuffer(s))
	decoder := json.NewDecoder(requestBody)
	var t Invalidate
	err := decoder.Decode(&t)
	if err != nil {
		panic(err)
	}
	return &t
}

// ToJSON will convert Invalidate to json
// Later, it will also responsible to handling the json marshalling error
func (d Invalidate) ToJSON() []byte {
	result, _ := json.Marshal(d)
	return result
}
