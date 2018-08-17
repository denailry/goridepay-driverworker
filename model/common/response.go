package response

import (
	"encoding/json"
)

// Response is type to store data for response of orderHandler
type Response struct {
	Error   bool
	Message string
}

// ToJSON will convert OrderResponse to json
// Later, it will also responsible to handling the json marshalling error
func (d Response) ToJSON() []byte {
	result, _ := json.Marshal(d)
	return result
}
