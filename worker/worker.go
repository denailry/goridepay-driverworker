package worker

import (
	"goridepay-driverworker/common"
	"goridepay-driverworker/invalidator"
	"goridepay-driverworker/model/order"
	"sync"
	"time"
)

// Worker is responsible for:
//   - queueing orders coming for driver
//   - sending order notification to driver by taking the nearest originDestination first
// It handles the threading, so queuing and notification process won't block the entire service
type Worker struct {
	DriverID       int
	isPrioritizing bool
	isOffering     bool
	isNotifying    bool
	isConfirming   bool
	isCleaning     bool
	orderQueue     []*order.Order
	orderPending   []*order.Order
	pendingLock    *sync.Mutex
	queueLock      *sync.Mutex
	rejectChan     chan bool
	confirmChan    chan bool
	previosOrderID int
}

// Always use getWorkerIndex to get element from workerList
var workerList = make([]*Worker, common.MaxWorker)

func getWorkerIndex(driverID int) int {
	return (driverID % common.MaxWorker) + 1
}

// GetOrderList will return order queue of the certain driver
func GetOrderList(driverID int) []*order.Order {
	w := NewWorker(driverID)
	w.queueLock.Lock()
	clean(&w.orderQueue)
	w.queueLock.Unlock()
	return w.orderQueue
}

// AddOrder is official way to add order to certain driver
func AddOrder(driverID int, order order.Order) {
	w := NewWorker(driverID)
	w.queue(order)
}

// RejectOrder is official way to reject order for certain driver
// func RejectOrder(driverID int, orderID int) bool {
// 	w := NewWorker(driverID)
// 	if orderID == w.previosOrderID {
// 		w.rejectChan <- true
// 		return true
// 	}
// 	return false
// }

// AcceptOrder is official way to accept order for certain driver
// func AcceptOrder(driverID int, orderID int) bool {
// 	w := NewWorker(driverID)
// 	if orderID == w.previosOrderID && !w.isNotifying {
// 		w.isConfirming = true
// 		w.confirmChan <- confirmOrder(orderID)
// 		w.isConfirming = false
// 		return true
// 	}
// 	return false
// }

// Confirm the order wheter it is taken or not
// func confirmOrder(orderID int) bool {
// 	return true
// }

// NewWorker return worker stored in workerList or create the new one if worker pointer is nil
func NewWorker(driverID int) *Worker {
	if workerList[getWorkerIndex(driverID)] == nil {
		w := Worker{
			DriverID:       driverID,
			isPrioritizing: false,
			isOffering:     false,
			isNotifying:    false,
			isConfirming:   false,
			isCleaning:     false,
			pendingLock:    &sync.Mutex{},
			queueLock:      &sync.Mutex{},
			rejectChan:     make(chan bool),
			confirmChan:    make(chan bool),
			previosOrderID: -1,
		}
		return &w
	}
	return workerList[getWorkerIndex(driverID)]
}

// func (d Worker) startOfferingDriver() {
// 	d.isOffering = true
// 	for len(d.orderQueue) > 0 {
// 		if d.isConfirming {
// 			accepted := <-d.confirmChan
// 			if accepted {
// 				d.orderQueue = nil
// 				break
// 			}
// 		}
// 		success := d.pushNotification()
// 		if success {
// 			d.isNotifying = false
// 			select {
// 			case <-d.rejectChan:
// 			case <-time.After(5000 * time.Millisecond):
// 			}
// 			d.isNotifying = true
// 		}
// 	}
// 	d.isNotifying = false
// 	d.isOffering = false
// }

func (d Worker) queue(order order.Order) {
	// if !d.isOffering {
	// 	go d.startOfferingDriver()
	// }
	d.pendingLock.Lock()
	d.orderPending = append(d.orderPending, &order)
	d.pendingLock.Unlock()
	if !d.isPrioritizing {
		go d.prioritize()
	}
}

func pop(lock *sync.Mutex, pq *[]*order.Order) order.Order {
	lock.Lock()
	q := *pq
	order := *q[0]
	temp := q[1:]
	pq = &temp
	lock.Unlock()
	return order
}

func (d Worker) insert(idx int, order order.Order) {
	d.queueLock.Lock()
	if idx >= len(d.orderQueue) {
		d.orderQueue = append(d.orderQueue, &order)
	} else {
		temp := append(d.orderQueue[:idx], &order)
		d.orderQueue = append(temp, d.orderQueue[idx:]...)
	}
	d.queueLock.Unlock()
}

func (d Worker) prioritize() {
	d.isPrioritizing = true
	for len(d.orderPending) > 0 {
		order := pop(d.pendingLock, &d.orderPending)
		if len(d.orderQueue) == 0 {
			d.insert(-1, order)
		} else {
			d.insert(findSmallerIndex(d.orderQueue, order), order)
		}
	}
	d.isPrioritizing = false
}

func findSmallerIndex(q []*order.Order, o order.Order) int {
	i := len(q) - 1
	result := -1
	for result != -1 && i >= 0 {
		c := *q[i]
		if c.OriginDistance < o.OriginDistance {
			result = i
		} else {
			i--
		}
	}
	return result + 1
}

// func (d Worker) pushNotification() bool {
// 	order := pop(d.queueLock, &d.orderQueue)
// 	valid := isValid(*order.Info)
// 	if valid {
// 		// Push notification
// 		d.previosOrderID = order.Info.OrderID
// 	}
// 	return valid
// }

func isValid(orderInfo order.Info) bool {
	return orderInfo.Timestamp-time.Now().Unix() < common.MaxOrderWaitingTime &&
		invalidator.IsValid(orderInfo.OrderID)
}

func (d Worker) runCleaner() {
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
		if isValid(*order.Info) {
			(*q) = append((*q)[:i], (*q)[i+1:]...)
		} else {
			i++
		}
	}
}
