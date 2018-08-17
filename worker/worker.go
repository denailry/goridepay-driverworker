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
	isOffering     bool
	isNotifying    bool
	isConfirming   bool
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

// AddOrder is official way to add order to certain driver
func AddOrder(driverID int, order order.Order) {
	w := NewWorker(driverID)
	w.queue(order)
}

// RejectOrder is official way to reject order for certain driver
func RejectOrder(driverID int, orderID int) bool {
	w := NewWorker(driverID)
	if orderID == w.previosOrderID {
		w.rejectChan <- true
		return true
	}
	return false
}

// AcceptOrder is official way to accept order for certain driver
func AcceptOrder(driverID int, orderID int) bool {
	w := NewWorker(driverID)
	if orderID == w.previosOrderID && !w.isNotifying {
		w.isConfirming = true
		w.confirmChan <- confirmOrder(orderID)
		w.isConfirming = false
		return true
	}
	return false
}

func confirmOrder(orderID int) bool {
	// Confirm the order wheter it is taken or not
	return true
}

// NewWorker return worker stored in workerList or create the new one if worker pointer is nil
func NewWorker(driverID int) *Worker {
	if workerList[getWorkerIndex(driverID)] == nil {
		w := Worker{
			DriverID:       driverID,
			ready:          true,
			isPrioritizing: false,
			isOffering:     false,
			isNotifying:    false,
			isConfirming:   false,
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

func (d Worker) startOfferingDriver() {
	d.isOffering = true
	for d.ready && len(d.orderQueue) > 0 {
		if d.isConfirming {
			accepted := <-d.confirmChan
			if accepted {
				d.ready = false
				d.orderQueue = nil
				break
			}
		}
		d.pushNotification()
		d.isNotifying = false
		select {
		case <-d.rejectChan:
		case <-time.After(5000 * time.Millisecond):
		}
		d.isNotifying = true
		if !d.ready {
			d.orderQueue = nil
		}
	}
	d.isNotifying = false
	d.isOffering = false
}

func (d Worker) queue(order order.Order) {
	if d.ready {
		if !d.isOffering {
			go d.startOfferingDriver()
		}
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
