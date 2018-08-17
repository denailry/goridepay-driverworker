package worker

import (
	"goridepay-driverworker/common"
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
	ready          bool
	isPrioritizing bool
	orderQueue     []*order.Order
	orderPending   []*order.Order
	pendingLock    *sync.Mutex
	queueLock      *sync.Mutex
	rejectChannel  chan bool
	previosOrderID int
}

// Always use getWorkerIndex to get element from workerList
var workerList = make([]*Worker, common.MaxWorker)

func getWorkerIndex(driverID int) int {
	return (driverID % common.MaxWorker) + 1
}

// AddOrder is official way to add order to certain driver
func AddOrder(driverID int, order order.Order) {
	w := NewWorker(driverID)
	w.queue(order)
}

// RejectOrder is official way to reject order for certain driver
func RejectOrder(driverID int, orderID int) bool {
	w := NewWorker(driverID)
	if orderID == w.previosOrderID {
		w.rejectChannel <- true
		return true
	}
	return false
}

// NewWorker return worker stored in workerList or create the new one if worker pointer is nil
func NewWorker(driverID int) *Worker {
	if workerList[getWorkerIndex(driverID)] == nil {
		w := Worker{
			DriverID:       driverID,
			ready:          true,
			isPrioritizing: false,
			pendingLock:    &sync.Mutex{},
			queueLock:      &sync.Mutex{},
			rejectChannel:  make(chan bool),
			previosOrderID: -1,
		}
		go w.startOfferingDriver()
		return &w
	}
	return workerList[getWorkerIndex(driverID)]
}

func (d Worker) startOfferingDriver() {
	for d.ready && len(d.orderQueue) > 0 {
		go d.pushNotification()
		select {
		case <-d.rejectChannel:
		case <-time.After(5000 * time.Millisecond):
		}
		if !d.ready {
			d.orderQueue = nil
		}
	}
}

func (d Worker) queue(order order.Order) {
	if d.ready {
		d.pendingLock.Lock()
		d.orderPending = append(d.orderPending, &order)
		d.pendingLock.Unlock()
		if !d.isPrioritizing {
			go d.prioritize()
		}
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
	if idx == -1 {
		d.orderQueue = append(d.orderQueue, &order)
	} else {
		temp := append(d.orderQueue[:idx], &order)
		d.orderQueue = append(temp, d.orderQueue[idx:]...)
	}
	d.queueLock.Unlock()
}

func (d Worker) prioritize() {
	d.isPrioritizing = true
	for d.ready && len(d.orderPending) > 0 {
		order := pop(d.pendingLock, &d.orderPending)
		if len(d.orderQueue) == 0 {
			d.insert(-1, order)
		} else {
			d.insert(findSmallerIndex(d.orderQueue, order), order)
		}
		if !d.ready {
			d.orderPending = nil
		}
	}
	d.isPrioritizing = false
}

func findSmallerIndex(q []*order.Order, o order.Order) int {
	i := 0
	result := -1
	for result != -1 && i < len(q) {
		c := *q[i]
		if c.OriginDistance < o.OriginDistance {
			result = i
		} else {
			i++
		}
	}
	return result
}

func (d Worker) pushNotification() {
	order := pop(d.queueLock, &d.orderQueue)
	// Push notification
	d.previosOrderID = order.Info.OrderID
}
