package model

// Order is official type to store data needed by Worker
type Order struct {
	Info                *OrderInfo
	OriginDistance      int
	DestinationDistance int
}

// OrderInfo will store data of each incoming order
// It is needed to avoid duplication of data of order information
// that is held by each Worker
type OrderInfo struct {
	OrderID     int
	Origin      string
	Destination string
}
