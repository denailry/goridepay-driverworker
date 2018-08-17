package invalidate

import (
	"time"
)

// InvalidOrder will store data used by invalidator
type InvalidOrder struct {
	OrderID   int
	Timestamp int64
}

// NewInvalidOrder is official way to create new variable with InvalidOrder type
func NewInvalidOrder(orderID int) *InvalidOrder {
	return &InvalidOrder{
		OrderID:   orderID,
		Timestamp: time.Now().Unix(),
	}
}
