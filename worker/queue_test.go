package worker

import (
	"goridepay-driverworker/model/order"
	"testing"
	"time"
)

var info1 = &order.Info{
	OrderID:             1,
	Origin:              "somewhere",
	Destination:         "nowhere",
	DestinationDistance: 1000,
	Timestamp:           time.Now().Unix(),
}
var o1 = &order.Order{
	OriginDistance: 5,
	Info:           info1,
}
var info2 = &order.Info{
	OrderID:             2,
	Origin:              "nowhere",
	Destination:         "somewhere",
	DestinationDistance: 1000,
	Timestamp:           time.Now().Unix(),
}
var o2 = &order.Order{
	OriginDistance: 10,
	Info:           info2,
}
var info3 = &order.Info{
	OrderID:             3,
	Origin:              "everywhere",
	Destination:         "anywhere",
	DestinationDistance: 80,
	Timestamp:           time.Now().Unix(),
}
var o3 = &order.Order{
	OriginDistance: 100,
	Info:           info3,
}
var info4 = &order.Info{
	OrderID:             4,
	Origin:              "everytime",
	Destination:         "anytime",
	DestinationDistance: 1000,
	Timestamp:           time.Now().Unix(),
}
var o4 = &order.Order{
	OriginDistance: 8,
	Info:           info4,
}

func TestPopZeroLengthQueue(t *testing.T) {
	w := NewWorker(1)
	o := pop(w.queueLock, &w.orderQueue)
	if o != nil {
		t.Fatalf("Expected: nil; Got: %s", o.ToString())
	}
}

func TestPopNonZeroLengthQueue(t *testing.T) {
	w := NewWorker(2)
	info := &order.Info{
		OrderID:             1,
		Origin:              "somewhere",
		Destination:         "nowhere",
		DestinationDistance: 1000,
		Timestamp:           time.Now().Unix(),
	}
	o := &order.Order{
		OriginDistance: 5,
		Info:           info,
	}
	w.orderQueue = append(w.orderQueue, o)
	op := pop(w.queueLock, &w.orderQueue)
	if op != o {
		t.Fatalf("Expected: op = %s; Got: op = %s", o.ToString(), op.ToString())
	}
}

func TestInsertInTheBeginning(t *testing.T) {
	w := NewWorker(3)
	w.orderQueue = append(w.orderQueue, o1)
	w.insert(0, *o2)
	ot := w.orderQueue[0]
	if ot.Info.OrderID != o2.Info.OrderID {
		t.Fatalf("Expected: ot = %s; Got: ot = %s", o2.ToString(), ot.ToString())
	}
}

func TestInsertInTheEnd(t *testing.T) {
	w := NewWorker(4)
	w.orderQueue = append(w.orderQueue, o1)
	w.insert(1, *o2)
	ot := w.orderQueue[1]
	if ot.Info.OrderID != o2.Info.OrderID {
		t.Fatalf("Expected: ot = %s; Got: ot = %s", o2.ToString(), ot.ToString())
	}
}

func TestInsertInTheMiddle(t *testing.T) {
	w := NewWorker(5)
	w.orderQueue = append(w.orderQueue, o1)
	w.orderQueue = append(w.orderQueue, o2)
	w.insert(1, *o3)
	ot := w.orderQueue[1]
	if ot.Info.OrderID != o3.Info.OrderID {
		t.Fatalf("Expected: ot = %s; Got: ot = %s", o3.ToString(), ot.ToString())
	}
}

func TestFindingSmallerIndex(t *testing.T) {
	w := NewWorker(6)
	w.orderQueue = append(w.orderQueue, o1)
	w.orderQueue = append(w.orderQueue, o2)
	w.orderQueue = append(w.orderQueue, o3)
	idx := findSmallerIndex(w.orderQueue, *o4)
	if idx != 1 {
		t.Fatalf("Expected: 1; Got: %d", idx)
	}
}
