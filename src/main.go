package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"github.com/senpainikolay/Kitchen/orders"
)

//global
var orderList = orders.OrderList{}
var olController = orders.OrderListPickUpController{}
var cooks = orders.GetCooks()
var cookProfiencySum = GetSumProfiencyOfCooks()
var Menu = orders.GetFoods()
var conf = orders.GetConf()
var oven, stove = orders.GetApparatus(Menu, conf.NR_OF_STOVES, conf.NR_OF_OVENS)

const (
	TIME_UNIT = 100
)

func PostDingHallOrders(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	var order orders.Order
	err := json.NewDecoder(r.Body).Decode(&order)
	if err != nil {
		log.Fatalln("There was an error decoding the request body into the struct")
	}

	orderList.Append(&order)

	olController.Mutex.Lock()
	olController.CounterOrdersPickedUp += 1
	olController.FoodCounter += len(order.Items)
	olController.Mutex.Unlock()

	go func() { orders.DistributeFoods(&orderList, cooks, Menu, conf.DiningHallAddress, &olController) }()
	// log.Printf("Order id %v recieved at Kitchen!", order.OrderId)

}

func main() {

	rand.Seed(time.Now().UnixMilli())

	for i := 0; i < len(cooks.Cook); i++ {
		cooks.Cook[i].CookChan = make(chan *orders.CookingDetails, cooks.Cook[i].Proficiency)
		cooks.Cook[i].CondVar = *sync.NewCond(&sync.Mutex{})
		cooks.Cook[i].Queue = make(chan *orders.CookingDetails, cooks.Cook[i].Proficiency)
		cooks.Cook[i].CounterAvailable = 0
	}

	r := mux.NewRouter()
	r.HandleFunc("/order", PostDingHallOrders).Methods("POST")
	r.HandleFunc("/estimationCalculation", GetEWT).Methods("GET")
	r.HandleFunc("/getOrderStatus", GetOrdersStatus).Methods("GET")
	r.HandleFunc("/getPreparedItems", GetPreparedItems).Methods("GET")

	for i := 0; i < len(cooks.Cook); i++ {
		idx := i
		go cooks.Cook[idx].Work(&orderList, cooks, oven, stove, Menu, conf.DiningHallAddress, &olController)
	}

	go oven.Work(&olController, &orders.CAOrdController{Counter: 0, C: *sync.NewCond(&sync.Mutex{})})
	go stove.Work(&olController, &orders.CAOrdController{Counter: 0, C: *sync.NewCond(&sync.Mutex{})})

	http.ListenAndServe(":"+conf.Port, r)

}

func GetEWT(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	var D, E int
	D = conf.NR_OF_OVENS + conf.NR_OF_STOVES
	olController.Mutex.Lock()
	E = olController.FoodCounter
	olController.Mutex.Unlock()

	BDE := orders.EWTCalculation{B: cookProfiencySum, D: D, E: E}

	resp, _ := json.Marshal(&BDE)
	fmt.Fprint(w, string(resp))

}

func GetOrdersStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	olController.Mutex.Lock()
	temp := olController.CounterOrdersPickedUp
	olController.Mutex.Unlock()
	// Buffer control
	if temp >= 5 {
		fmt.Fprint(w, 1)
	} else {
		fmt.Fprint(w, 0)

	}

}

func GetPreparedItems(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	olController.Mutex.Lock()
	temp := olController.PreparedItems
	olController.Mutex.Unlock()

	fmt.Fprint(w, temp)

}

func GetSumProfiencyOfCooks() int {
	var sum int
	for _, cook := range cooks.Cook {
		sum += cook.Proficiency

	}
	return sum

}
