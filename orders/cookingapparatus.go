package orders

import (
	"sync"
	"time"
)

type CookingApparatus struct {
	Counter  int
	Quantity int
	C        sync.Cond
	Menu     *Foods
}

func GetApparatus(Menu *Foods) (*CookingApparatus, *CookingApparatus) {
	Oven := CookingApparatus{0, NR_OF_OVENS, *sync.NewCond(&sync.Mutex{}), Menu}
	Stove := CookingApparatus{0, NR_OF_STOVES, *sync.NewCond(&sync.Mutex{}), Menu}
	return &Oven, &Stove
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
	time.Sleep(time.Duration(int64(ca.Menu.Foods[cd.FoodId-1].PreparationTime) * TIME_UNIT * int64(time.Millisecond)))
	ca.C.L.Lock()
	ca.Counter -= 1
	ca.C.Signal()
	cd.CookId = cookId
	cd.wg.Done()
	ca.C.L.Unlock()

}
