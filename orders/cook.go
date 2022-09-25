package orders

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"sync"
	"time"
)

type Cook struct {
	Id          int    `json:"id"`
	Rank        int    `json:"rank"`
	Proficiency int    `json:"proficiency"`
	Name        string `json:"name"`
	CatchPhrase string `json:"catch_phrase"`
	CookChan    chan *CookingDetails
}

type Cooks struct {
	Cook []Cook `json:"cooks"`
}
type CookingDetails struct {
	CookId int `json:"cook_id"`
	FoodId int `json:"food_id"`
	wg     *sync.WaitGroup
}

func GetCooks() *Cooks {
	jsonFile, err := os.Open("configurations/Cooks.json")
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	var cooks Cooks
	json.Unmarshal(byteValue, &cooks)
	return &cooks

}

func (c *Cook) PickUpOrder(orderList *OrderList, cooks *Cooks) {
	rand.Seed(time.Now().UnixNano())
	orderList.Mutex.Lock()
	order, orderBool := orderList.PickUp()
	orderList.Mutex.Unlock()
	if orderBool {
		//fmt.Printf("NO ORDER TO Recieve for this cook id %v\n", c.Id)
		return
	}
	fmt.Printf("OREEDRRR:     %+v\n", order)
	var wg sync.WaitGroup
	wg.Add(len(order.Items))

	payload := Payload{Order: *order}
	for i, foodId := range payload.Items {
		payload.CookingDetails = append(payload.CookingDetails, CookingDetails{FoodId: foodId})
		payload.CookingDetails[i].wg = &wg
	}

	tempOrders := order.Items
	idCounter := 0
	for {
		if len(tempOrders) == 0 {
			break
		}
		rnd := rand.Intn(2)
		switch rnd {

		case 0:

			for i := 0; i < len(cooks.Cook); i++ {
				if cooks.Cook[i].Rank == Menu.Foods[tempOrders[0]-1].Complexity || cooks.Cook[i].Rank == Menu.Foods[tempOrders[0]-1].Complexity-1 {
					idx := i
					tIdC := idCounter
					go func() { cooks.Cook[idx].CookChan <- &payload.CookingDetails[tIdC] }()
					idCounter += 1

					if len(tempOrders) <= 1 {
						tempOrders = make([]int, 0)
						break
					} else {
						tempOrders = remove(tempOrders)
						break
					}

				}
			}

		case 1:
			for i := len(cooks.Cook) - 1; i >= 0; i-- {
				if cooks.Cook[i].Rank == Menu.Foods[tempOrders[0]-1].Complexity || cooks.Cook[i].Rank == Menu.Foods[tempOrders[0]-1].Complexity-1 {
					idx := i
					tIdC := idCounter
					go func() { cooks.Cook[idx].CookChan <- &payload.CookingDetails[tIdC] }()
					idCounter += 1
					if len(tempOrders) <= 1 {
						tempOrders = make([]int, 0)
						break
					} else {
						tempOrders = remove(tempOrders)
						break
					}

				}
			}
		}

	}

	wg.Wait()
	fmt.Printf("%+v\n", payload)

}

func (c *Cook) Work(orderList *OrderList, cooks *Cooks) {

	for {

		select {
		case cd := <-c.CookChan:
			time.Sleep(time.Duration(Menu.Foods[cd.FoodId-1].PreparationTime) * 50 * time.Millisecond)
			cd.CookId = c.Id
			cd.wg.Done()
		default:
			go func() {
				c.PickUpOrder(orderList, cooks)
			}()

			time.Sleep(5 * time.Second)
		}
	}
}

func remove(slice []int) []int {
	return slice[1:]
}
