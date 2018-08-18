package order

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
