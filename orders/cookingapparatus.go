package orders

import (
	"sync"
	"time"
)

type CookingApparatus struct {
	Counter  int
	Quantity int
	C        sync.Cond
}

func (ca *CookingApparatus) borrow() {
	ca.C.L.Lock()
	for ca.Counter >= ca.Quantity {
		ca.C.Wait()
	}
	ca.Counter += 1
	ca.C.L.Unlock()
}

func (ca *CookingApparatus) Use(cd *CookingDetails, cookId int) {
	ca.borrow()
	time.Sleep(time.Duration(int64(Menu.Foods[cd.FoodId-1].PreparationTime) * TIME_UNIT * int64(time.Millisecond)))
	ca.C.L.Lock()
	ca.Counter -= 1
	ca.C.Signal()
	cd.CookId = cookId
	cd.wg.Done()
	ca.C.L.Unlock()

}
