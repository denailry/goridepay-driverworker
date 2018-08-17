package response

import (
	"encoding/json"
)

// Simple is the most simple response structure
type Simple struct {
	Error   bool
	Message string
}

// ToJSON will convert OrderResponse to json
// Later, it will also responsible to handling the json marshalling error
func (d Simple) ToJSON() []byte {
	result, _ := json.Marshal(d)
	return result
}
