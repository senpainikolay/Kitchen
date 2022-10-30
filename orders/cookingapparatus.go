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

func GetApparatus(Menu *Foods, NR_OF_OVENS int, NR_OF_STOVES int) (*CookingApparatus, *CookingApparatus) {
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

func (ca *CookingApparatus) Use(cd *CookingDetails, cook *Cook, olController *OrderListPickUpController) {
	ca.borrow()
	if cd.TempPreparationTime <= 10 {
		time.Sleep(time.Duration(int64(cd.TempPreparationTime) * TIME_UNIT * int64(time.Millisecond)))
		cd.CookId = cook.Id
		cd.wg.Done()
		go FoodCounterDecreaser(olController)
	} else {
		time.Sleep(time.Duration(10 * TIME_UNIT * int64(time.Millisecond)))
		cd.TempPreparationTime -= 10
		cook.Queue <- cd
	}
	ca.C.L.Lock()
	ca.Counter -= 1
	ca.C.Signal()
	ca.C.L.Unlock()

}
