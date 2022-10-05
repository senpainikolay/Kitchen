package orders

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

const (
	TIME_UNIT    = 100
	NR_OF_OVENS  = 2
	NR_OF_STOVES = 1
)

type Cook struct {
	Id               int    `json:"id"`
	Rank             int    `json:"rank"`
	Proficiency      int    `json:"proficiency"`
	Name             string `json:"name"`
	CatchPhrase      string `json:"catch_phrase"`
	CookChan         chan *CookingDetails
	CondVar          sync.Cond
	Queue            chan *CookingDetails
	CounterAvailable int
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
	orderList.Mutex.Lock()
	order, orderBool := orderList.PickUp()
	orderList.Mutex.Unlock()
	if orderBool {
		return
	}
	var wg sync.WaitGroup
	wg.Add(len(order.Items))

	payload := Payload{
		OrderId:    order.OrderId,
		Items:      order.Items,
		Priority:   order.Priority,
		MaxWait:    order.MaxWait,
		PickUpTime: order.PickUpTime,
		TableId:    order.TableId,
		WaiterId:   order.WaiterId,
	}
	for i, foodId := range payload.Items {
		payload.CookingDetails = append(payload.CookingDetails, CookingDetails{FoodId: foodId})
		payload.CookingDetails[i].wg = &wg
	}

	oldTime := time.Now().UnixMilli()
	tempOrders := order.Items
	FoodIdCounter := 0
	i := 0
	for {
		if len(tempOrders) == 0 {
			break
		}
		if i == len(cooks.Cook) {
			i = 0
		}
		if len(cooks.Cook[i].Queue) < cooks.Cook[i].Proficiency+4 && cooks.Cook[i].Rank == Menu.Foods[tempOrders[0]-1].Complexity || cooks.Cook[i].Rank == Menu.Foods[tempOrders[0]-1].Complexity-1 {
			tFIC := FoodIdCounter
			idx := i
			go func() {
				cooks.Cook[idx].Queue <- &payload.CookingDetails[tFIC]
			}()
			FoodIdCounter += 1
			if len(tempOrders) <= 1 {
				tempOrders = make([]int, 0)
			} else {
				tempOrders = popFront(tempOrders)
			}
		}
		i += 1
	}
	wg.Wait()
	payload.CookingTime = (time.Now().UnixMilli() - oldTime) / int64(TIME_UNIT)
	SendOrder(&payload)
	log.Printf("Order id %v sent back to dining hall", payload.OrderId)

}

func (c *Cook) Work(orderList *OrderList, cooks *Cooks, Oven *CookingApparatus, Stove *CookingApparatus) {

	for {
		select {
		case cd := <-c.CookChan:
			tempCd := cd
			go func() {
				switch Menu.Foods[tempCd.FoodId-1].CookingApparatus {
				case "oven":
					go func() { Oven.Use(tempCd, c.Id) }()
					c.CondVar.L.Lock()
					c.CounterAvailable -= 1
					c.CondVar.Signal()
					c.CondVar.L.Unlock()
				case "stove":
					go func() { Stove.Use(tempCd, c.Id) }()
					c.CondVar.L.Lock()
					c.CounterAvailable -= 1
					c.CondVar.Signal()
					c.CondVar.L.Unlock()

				default:
					time.Sleep(time.Duration(int64(Menu.Foods[cd.FoodId-1].PreparationTime) * TIME_UNIT * int64(time.Millisecond)))
					c.CondVar.L.Lock()
					c.CounterAvailable -= 1
					c.CondVar.Signal()
					tempCd.CookId = c.Id
					tempCd.wg.Done()
					c.CondVar.L.Unlock()

				}
			}()
		// Thread Controller On Cooking items by Cook's Proficiency.
		case cda := <-c.Queue:
			tempCd := cda
			go func() {
				c.CondVar.L.Lock()
				for c.CounterAvailable >= c.Proficiency {
					c.CondVar.Wait()
				}
				c.CounterAvailable += 1
				c.CookChan <- tempCd
				c.CondVar.L.Unlock()
			}()

		default:
			// PickUpTime
			go func() {
				c.PickUpOrder(orderList, cooks)
			}()
			time.Sleep(1 * time.Millisecond)

		}

	}
}

func SendOrder(ord *Payload) {
	postBody, _ := json.Marshal(*ord)
	responseBody := bytes.NewBuffer(postBody)
	resp, err := http.Post("http://localhost:8080/distribution", "application/json", responseBody)
	if err != nil {
		log.Fatalf("An Error Occured %v", err)
	}
	defer resp.Body.Close()
	//Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	sb := string(body)
	log.Printf(sb)

}

func popFront(slice []int) []int {
	return slice[1:]
}
