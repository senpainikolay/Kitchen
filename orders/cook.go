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

const TIME_UNIT = 50

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
		return
	}
	var wg sync.WaitGroup
	wg.Add(len(order.Items))

	payload := Payload{Order: *order}
	for i, foodId := range payload.Items {
		payload.CookingDetails = append(payload.CookingDetails, CookingDetails{FoodId: foodId})
		payload.CookingDetails[i].wg = &wg
	}

	oldTime := time.Now().Unix()

	tempOrders := order.Items
	idCounter1 := 0
	idCounter2 := len(order.Items) - 1

	for {
		if len(tempOrders) == 0 {
			break
		}
		rnd := rand.Intn(2)
		switch rnd {

		case 0:

			for i := 0; i < len(cooks.Cook); i++ {
				if len(cooks.Cook[i].CookChan) < 5 && (cooks.Cook[i].Rank == Menu.Foods[tempOrders[0]-1].Complexity || cooks.Cook[i].Rank == Menu.Foods[tempOrders[0]-1].Complexity-1) {
					idx := i
					tIdC := idCounter1
					go func() { cooks.Cook[idx].CookChan <- &payload.CookingDetails[tIdC] }()
					idCounter1 += 1

					if len(tempOrders) <= 1 {
						tempOrders = make([]int, 0)
						break
					} else {
						tempOrders = popFront(tempOrders)
						continue
					}
				}
				time.Sleep(2 * TIME_UNIT * time.Millisecond)

			}

		case 1:

			for i := len(cooks.Cook) - 1; i >= 0; i-- {
				if len(cooks.Cook[i].CookChan) < 5 && (cooks.Cook[i].Rank == Menu.Foods[tempOrders[0]-1].Complexity || cooks.Cook[i].Rank == Menu.Foods[tempOrders[0]-1].Complexity-1) {
					idx := i
					tIdC := idCounter2
					go func() { cooks.Cook[idx].CookChan <- &payload.CookingDetails[tIdC] }()
					idCounter2 -= 1
					if len(tempOrders) <= 1 {
						tempOrders = make([]int, 0)
						break
					} else {
						tempOrders = popBack(tempOrders)
						continue
					}

				}

				time.Sleep(2 * TIME_UNIT * time.Millisecond)

			}
		}

	}

	wg.Wait()
	payload.CookingTime = time.Now().Unix() - oldTime
	fmt.Printf("%+v\n", payload)

}

func (c *Cook) Work(orderList *OrderList, cooks *Cooks) {

	for {

		select {
		case cd := <-c.CookChan:
			if len(c.CookChan) >= 1 {
				go func() {
					time.Sleep(time.Duration(Menu.Foods[cd.FoodId-1].PreparationTime) * TIME_UNIT * time.Millisecond)
					cd.CookId = c.Id
					cd.wg.Done()
				}()
				for i := 0; i < c.Proficiency-1; i++ {
					if len(c.CookChan) == 0 {
						break
					}
					cdExtra := <-c.CookChan
					go func() {
						time.Sleep(time.Duration(Menu.Foods[cdExtra.FoodId-1].PreparationTime) * TIME_UNIT * time.Millisecond)
						cdExtra.CookId = c.Id
						cdExtra.wg.Done()
					}()

				}
			} else {
				time.Sleep(time.Duration(Menu.Foods[cd.FoodId-1].PreparationTime) * TIME_UNIT * time.Millisecond)
				cd.CookId = c.Id
				cd.wg.Done()
			}

		default:
			go func() {
				c.PickUpOrder(orderList, cooks)
			}()
			time.Sleep(TIME_UNIT * 3 * time.Millisecond)
		}
	}
}

func popFront(slice []int) []int {
	return slice[1:]
}

func popBack(slice []int) []int {
	return slice[:len(slice)-1]
}
