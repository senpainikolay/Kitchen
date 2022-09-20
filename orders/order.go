package orders

import "sync"

type Order struct {
	OrderId    int     `json:"order_id"`
	Items      []int   `json:"items"`
	Priority   int     `json:"priority"`
	MaxWait    float64 `json:"max_wait"`
	PickUpTime int64   `json:"pick_up_time"`
	TableId    int     `json:"table_id"`
	WaiterId   int     `json:"waiter_id"`
	Wg         sync.WaitGroup
}

type OrderList struct {
	Orders []Order
}

func (ol *OrderList) Append(o *Order) {
	ol.Orders = append(ol.Orders, *o)
}

func (ol *OrderList) IsEmpty() bool {
	if len(ol.Orders) == 0 {
		return true
	}
	return false
}

func (ol *OrderList) PickUp() *Order {
	elem := ol.Orders[0]
	sliced := ol.Orders[1:]
	ol.Orders = sliced
	return &elem
}
