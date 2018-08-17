package worker

import (
	"goridepay-driverworker/common"
	"sync"
	"time"
)

type Worker struct {
	DriverID       int
	ready          bool
	isPrioritizing bool
	orderQueue     []*common.Order
	orderPending   []*common.Order
	pendingLock    *sync.Mutex
	queueLock      *sync.Mutex
}

var workerList = make([]*Worker, common.MaxWorker)

func AddOrder(driverID int, order common.Order) {
	w := NewWorker(driverID)
	w.queue(order)
}

func NewWorker(driverID int) *Worker {
	if workerList[driverID] == nil {
		w := Worker{
			DriverID:       driverID,
			ready:          true,
			isPrioritizing: false,
			pendingLock:    &sync.Mutex{},
			queueLock:      &sync.Mutex{},
		}
		go w.startOfferingDriver()
		return &w
	}
	return workerList[driverID]
}

func (d Worker) startOfferingDriver() {
	for d.ready && len(d.orderQueue) > 0 {
		go d.pushNotification()
		time.Sleep(5000 * time.Millisecond)
		if !d.ready {
			d.orderQueue = nil
		}
	}
}

func (d Worker) queue(order common.Order) {
	if d.ready {
		d.pendingLock.Lock()
		d.orderPending = append(d.orderPending, &order)
		d.pendingLock.Unlock()
		if !d.isPrioritizing {
			go d.prioritize()
		}
	}
}

func pop(lock *sync.Mutex, pq *[]*common.Order) common.Order {
	lock.Lock()
	q := *pq
	order := *q[0]
	temp := q[1:]
	pq = &temp
	lock.Unlock()
	return order
}

func (d Worker) insert(idx int, order common.Order) {
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

func findSmallerIndex(q []*common.Order, o common.Order) int {
	i := 0
	result := -1
	for result != -1 && i < len(q) {
		c := *q[i]
		if c.Distance < o.Distance {
			result = i
		} else {
			i++
		}
	}
	return result
}

func (d Worker) pushNotification() {
	// order := pop(&d.queueLock, &d.orderQueue)
	// Push notification
}
