package worker

import (
	"goridepay-driverworker/common"
	"goridepay-driverworker/invalidator"
	"goridepay-driverworker/model/order"
	"time"
)

func isValid(orderInfo order.Info) bool {
	return time.Now().Unix()-orderInfo.Timestamp < common.MaxOrderWaitingTime &&
		invalidator.IsValid(orderInfo.OrderID)
}

func (d *Worker) runCleaner() {
	d.isCleaning = true
	for len(d.orderQueue) > 0 {
		d.queueLock.Lock()
		clean(&d.orderQueue)
		d.queueLock.Unlock()
		time.Sleep(common.MaxOrderWaitingTime * time.Second)
	}
	d.isCleaning = false
}

func clean(q *[]*order.Order) {
	i := 0
	for i < len(*q) {
		order := (*q)[i]
		if !isValid(*order.Info) {
			(*q) = append((*q)[:i], (*q)[i+1:]...)
		} else {
			i++
		}
	}
}
