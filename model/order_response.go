package model

import (
	"encoding/json"
)

// OrderResponse is type to store data for response of orderHandler
type OrderResponse struct {
	Error   bool
	Message string
}

// ToJSON will convert OrderResponse to json
// Later, it will also responsible to handling the json marshalling error
func (d OrderResponse) ToJSON() []byte {
	result, _ := json.Marshal(d)
	return result
}
