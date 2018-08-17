package invalidator

import (
	"goridepay-driverworker/common"
	"goridepay-driverworker/model/invalidate"
	"sync"
	"time"
)

// InvalidOrderList is list of cancelled or taken order
var InvalidOrderList []*invalidate.InvalidOrder
var listLock sync.RWMutex
var isRunning bool

// Invalidate is the official way to invalidate cancelled or taken order
func Invalidate(order *invalidate.InvalidOrder) {
	listLock.Lock()
	InvalidOrderList = append(InvalidOrderList, order)
	if !isRunning {
		go run()
	}
	listLock.Unlock()
}

// IsValid returns true if orderID is not in the InvalidOrderList
func IsValid(orderID int) bool {
	listLock.RLock()
	for _, order := range InvalidOrderList {
		if order.OrderID == orderID {
			return false
		}
	}
	listLock.RUnlock()
	return true
}

func run() {
	isRunning = true
	for len(InvalidOrderList) > 0 {
		listLock.Lock()
		i := 0
		for i < len(InvalidOrderList) {
			order := InvalidOrderList[i]
			if order.Timestamp-time.Now().Unix() >= common.MaxOrderWaitingTime {
				InvalidOrderList = append(InvalidOrderList[:i], InvalidOrderList[i+1:]...)
			} else {
				i++
			}
		}
		listLock.Unlock()
		time.Sleep(common.MaxOrderWaitingTime * time.Second)
	}
	isRunning = false
}
