package worker

import (
	"goridepay-driverworker/model/order"
	"sync"
)

func pop(lock *sync.Mutex, pq *[]*order.Order) *order.Order {
	lock.Lock()
	if len(*pq) == 0 {
		return nil
	}
	order := (*pq)[0]
	*pq = (*pq)[1:]
	lock.Unlock()
	return order
}

func (d *Worker) insert(idx int, o order.Order) {
	d.queueLock.Lock()
	if idx >= len(d.orderQueue) {
		d.orderQueue = append(d.orderQueue, &o)
	} else {
		d.orderQueue = append(d.orderQueue[:idx], append([]*order.Order{&o}, d.orderQueue[idx:]...)...)
	}
	d.queueLock.Unlock()
}

func findSmallerIndex(q []*order.Order, o order.Order) int {
	i := len(q) - 1
	result := -1
	for result == -1 && i >= 0 {
		c := *q[i]
		if c.OriginDistance < o.OriginDistance {
			result = i
		} else {
			i--
		}
	}
	return result + 1
}
