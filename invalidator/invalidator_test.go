package invalidator

import (
	"goridepay-driverworker/model/invalidate"
	"testing"
	"time"
)

var validOrder = invalidate.InvalidOrder{
	OrderID:   1,
	Timestamp: time.Now().Unix(),
}

var invalidOrder = invalidate.InvalidOrder{
	OrderID:   2,
	Timestamp: time.Now().Unix(),
}

func TestValidOrder(t *testing.T) {
	valid := IsValid(validOrder.OrderID)
	if valid == false {
		t.Fatalf("Expected: %t; Got: %t", true, valid)
	}
}

func TestInvalidOrder(t *testing.T) {
	Invalidate(&invalidOrder)
	valid := IsValid(invalidOrder.OrderID)
	if valid == true {
		t.Fatalf("Expected: %t; Got: %t", false, valid)
	}
}
