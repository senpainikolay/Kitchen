package orders

import (
	"sync"
)

type Order struct {
	OrderId    int     `json:"order_id"`
	Items      []int   `json:"items"`
	Priority   int     `json:"priority"`
	MaxWait    float64 `json:"max_wait"`
	PickUpTime int64   `json:"pick_up_time"`
	TableId    int     `json:"table_id"`
	WaiterId   int     `json:"waiter_id"`
}

type Payload struct {
	OrderId        int              `json:"order_id"`
	Items          []int            `json:"items"`
	Priority       int              `json:"priority"`
	MaxWait        float64          `json:"max_wait"`
	PickUpTime     int64            `json:"pick_up_time"`
	TableId        int              `json:"table_id"`
	WaiterId       int              `json:"waiter_id"`
	CookingTime    int64            `json:"cooking_time"`
	CookingDetails []CookingDetails `json:"cooking_details"`
}

type OrderList struct {
	Orders []Order
	Mutex  sync.Mutex
}

type OrderListPickUpController struct {
	Mutex                 sync.Mutex
	CounterOrdersPickedUp int
}

func (ol *OrderList) Append(o *Order) {
	ol.Mutex.Lock()
	defer ol.Mutex.Unlock()
	if ol.IsEmpty() {
		ol.Orders = append(ol.Orders, *o)
	} else {
		for i := 0; i < len(ol.Orders); i++ {
			if o.Priority > ol.Orders[i].Priority && i < len(ol.Orders) {
				ol.Insert(o, i)
				break
			}
			if o.Priority <= ol.Orders[i].Priority && len(ol.Orders)-1 == i {
				ol.Orders = append(ol.Orders, *o)
				break
			}
		}
	}

}

func (ol *OrderList) Insert(o *Order, index int) {
	ol.Orders = append(ol.Orders[:index+1], ol.Orders[index:]...)
	ol.Orders[index] = *o
}

func (ol *OrderList) IsEmpty() bool {
	if len(ol.Orders) == 0 {
		return true
	}
	return false
}

func (ol *OrderList) PickUp() (*Order, bool) {
	if len(ol.Orders) == 0 {
		return nil, true
	}
	elem := ol.Orders[0]
	sliced := ol.Orders[1:]
	ol.Orders = sliced
	return &elem, false
}
