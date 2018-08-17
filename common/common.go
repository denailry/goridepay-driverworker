package common

const (
	MaxWorker = 1000000
)

type Order struct {
	OrderId  int
	Distance int
}

var ServiceId int
