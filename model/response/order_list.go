package response

import (
	"encoding/json"
	"goridepay-driverworker/model/order"
)

// OrderList is the response of /get-order-list
type OrderList struct {
	Order []orderData
}

type orderData struct {
	OrderID             int
	Origin              string
	Destination         string
	OriginDistance      int
	DestinationDistance int
}

// NewOrderList is constructor of OrderList
func NewOrderList(list []*order.Order) *OrderList {
	ol := OrderList{}
	for _, data := range list {
		od := orderData{
			OrderID:             data.Info.OrderID,
			Origin:              data.Info.Origin,
			Destination:         data.Info.Destination,
			OriginDistance:      data.OriginDistance,
			DestinationDistance: data.Info.DestinationDistance,
		}
		ol.Order = append(ol.Order, od)
	}
	return &ol
}

// ToJSON will convert OrderList to json
// Later, it will also responsible to handling the json marshalling error
func (d OrderList) ToJSON() []byte {
	result, _ := json.Marshal(d.Order)
	return result
}
