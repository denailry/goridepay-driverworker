package worker

import (
	"fmt"
	"goridepay-driverworker/common"
	"goridepay-driverworker/model/order"
	"strconv"
	"sync"
)

// Worker is responsible for:
//   - queueing orders coming for driver
//   - sending order notification to driver by taking the nearest originDestination first
// It handles the threading, so queuing and notification process won't block the entire service
type Worker struct {
	DriverID       int
	isPrioritizing bool
	isConfirming   bool
	isCleaning     bool
	orderQueue     []*order.Order
	orderPending   []*order.Order
	pendingLock    *sync.Mutex
	queueLock      *sync.Mutex
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

// NewWorker return worker stored in workerList or create the new one if worker pointer is nil
func NewWorker(driverID int) *Worker {
	if workerList[getWorkerIndex(driverID)] == nil {
		fmt.Println("Create new driver with ID " + strconv.Itoa(driverID))
		w := Worker{
			DriverID:       driverID,
			isPrioritizing: false,
			isConfirming:   false,
			isCleaning:     false,
			pendingLock:    &sync.Mutex{},
			queueLock:      &sync.Mutex{},
		}
		workerList[getWorkerIndex(driverID)] = &w
		return &w
	}
	fmt.Println("Return instance driver with ID " + strconv.Itoa(driverID))
	return workerList[getWorkerIndex(driverID)]
}

func (d *Worker) queue(order order.Order) {
	d.pendingLock.Lock()
	d.orderPending = append(d.orderPending, &order)
	d.pendingLock.Unlock()
	if !d.isPrioritizing {
		go d.prioritize()
	}
	if !d.isCleaning {
		go d.runCleaner()
	}
}

func (d *Worker) prioritize() {
	d.isPrioritizing = true
	for len(d.orderPending) > 0 {
		order := pop(d.pendingLock, &d.orderPending)
		d.insert(findSmallerIndex(d.orderQueue, *order), *order)
	}
	d.isPrioritizing = false
}
