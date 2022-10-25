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

func (ca *CookingApparatus) Use(cd *CookingDetails, cook *Cook) {
	ca.borrow()
	if cd.TempPreparationTime <= 5 {
		time.Sleep(time.Duration(int64(cd.TempPreparationTime) * TIME_UNIT * int64(time.Millisecond)))
		cd.CookId = cook.Id
		cd.wg.Done()
	} else {
		time.Sleep(time.Duration(5 * TIME_UNIT * int64(time.Millisecond)))
		cd.TempPreparationTime -= 5
		cook.Queue <- cd
	}
	ca.C.L.Lock()
	ca.Counter -= 1
	ca.C.Signal()
	ca.C.L.Unlock()

}
