package order

import (
	"strconv"
)

// Order is official type to store data needed by Worker
type Order struct {
	Info           *Info
	OriginDistance int
}

// Info will store data of each incoming order
// It is needed to avoid duplication of data of order information
// that is held by each Worker
type Info struct {
	OrderID             int
	Origin              string
	Destination         string
	Timestamp           int64
	DestinationDistance int
}

// ToString creates readable string of Order
func (o Order) ToString() string {
	originDistance := "OriginDistance: " + strconv.Itoa(o.OriginDistance)
	orderID := "OrderID: " + strconv.Itoa(o.Info.OrderID)
	timestamp := "Timestamp: " + strconv.FormatInt(o.Info.Timestamp, 10)
	destinationDistance := "DestinationDistance: " + strconv.Itoa(o.Info.DestinationDistance)
	origin := "Origin: " + o.Info.Origin
	destination := "Destination: " + o.Info.Destination
	return "{" + orderID + "," + originDistance + "," + destinationDistance + "," + origin + "," + destination + "," + timestamp + "}"
}
