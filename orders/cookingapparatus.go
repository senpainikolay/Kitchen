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
	Queue    chan *CookingDetails
	OrdChan  chan *CookingDetails
}

type CAOrdController struct {
	Counter int
	C       sync.Cond
}

func GetApparatus(Menu *Foods, NR_OF_OVENS int, NR_OF_STOVES int) (*CookingApparatus, *CookingApparatus) {
	Oven := CookingApparatus{0, NR_OF_OVENS, *sync.NewCond(&sync.Mutex{}), Menu, make(chan *CookingDetails, 10), make(chan *CookingDetails, 10)}
	Stove := CookingApparatus{0, NR_OF_STOVES, *sync.NewCond(&sync.Mutex{}), Menu, make(chan *CookingDetails, 10), make(chan *CookingDetails, 10)}
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

func (ca *CookingApparatus) Use(cd *CookingDetails, olController *OrderListPickUpController, foodBuffer *CAOrdController) {
	ca.borrow()
	if cd.TempPreparationTime <= 5 {
		time.Sleep(time.Duration(int64(cd.TempPreparationTime) * TIME_UNIT * int64(time.Millisecond)))
		cd.wg.Done()
		go FoodCounterDecreaser(olController)
		go CountOrdersPrepared(olController)
		foodBuffer.C.L.Lock()
		foodBuffer.Counter -= 1
		foodBuffer.C.Signal()
		foodBuffer.C.L.Unlock()

	} else {
		time.Sleep(time.Duration(5 * TIME_UNIT * int64(time.Millisecond)))
		cd.TempPreparationTime -= 5
		go func() { ca.OrdChan <- cd }()
	}
	ca.C.L.Lock()
	ca.Counter -= 1
	ca.C.Signal()
	ca.C.L.Unlock()

}

func (ca *CookingApparatus) Work(olController *OrderListPickUpController, foodBuffer *CAOrdController) {
	for {
		select {
		case cd := <-ca.OrdChan:
			tempCd := cd
			go func() { ca.Use(tempCd, olController, foodBuffer) }()

		// Thread Controller
		case cda := <-ca.Queue:
			tempCd := cda
			go func() {
				foodBuffer.C.L.Lock()
				for foodBuffer.Counter >= 3 {
					foodBuffer.C.Wait()
				}
				foodBuffer.Counter += 1
				go func() { ca.OrdChan <- tempCd }()
				foodBuffer.C.L.Unlock()
			}()
		default:

			time.Sleep(1 * time.Millisecond)

		}
	}
}
