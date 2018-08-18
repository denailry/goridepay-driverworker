package worker

import (
	"fmt"
	"goridepay-driverworker/common"
	"goridepay-driverworker/invalidator"
	"goridepay-driverworker/model/order"
	"strconv"
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
		fmt.Println("Create new driver with ID " + strconv.Itoa(driverID))
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
		workerList[getWorkerIndex(driverID)] = &w
		return &w
	}
	fmt.Println("Return instance driver with ID " + strconv.Itoa(driverID))
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

func (d *Worker) queue(order order.Order) {
	// if !d.isOffering {
	// 	go d.startOfferingDriver()
	// }
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

func pop(lock *sync.Mutex, pq *[]*order.Order) order.Order {
	lock.Lock()
	order := *(*pq)[0]
	*pq = (*pq)[1:]
	lock.Unlock()
	return order
}

func (d *Worker) insert(idx int, o order.Order) {
	d.queueLock.Lock()
	if idx >= len(d.orderQueue) {
		fmt.Println("Length before append in queue of " + strconv.Itoa(d.DriverID) + " is " + strconv.Itoa(len(d.orderQueue)))
		d.orderQueue = append(d.orderQueue, &o)
		fmt.Println("Length after append in queue of " + strconv.Itoa(d.DriverID) + " is " + strconv.Itoa(len(d.orderQueue)))
		fmt.Println("Appended new element in queue of " + strconv.Itoa(d.DriverID))
	} else {
		d.orderQueue = append(d.orderQueue[:idx], append([]*order.Order{&o}, d.orderQueue[idx:]...)...)
		fmt.Println("Inserted new element in queue of " + strconv.Itoa(d.DriverID) +
			" at index " + strconv.Itoa(idx))
	}
	d.queueLock.Unlock()
}

func (d *Worker) prioritize() {
	d.isPrioritizing = true
	for len(d.orderPending) > 0 {
		order := pop(d.pendingLock, &d.orderPending)
		d.insert(findSmallerIndex(d.orderQueue, order), order)
		fmt.Println("Worker of " + strconv.Itoa(d.DriverID) + " has queue length " + strconv.Itoa(len(d.orderQueue)))
	}
	d.isPrioritizing = false
}

func findSmallerIndex(q []*order.Order, o order.Order) int {
	i := len(q) - 1
	result := -1
	for result == -1 && i >= 0 {
		fmt.Println("Iter " + strconv.Itoa(i))
		c := *q[i]
		if c.OriginDistance < o.OriginDistance {
			fmt.Println("Choose " + strconv.Itoa(i))
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

func printlist(list []*order.Order) {
	for _, data := range list {
		fmt.Println(data.Info.OrderID)
	}
}
