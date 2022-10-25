package main

import (
	"encoding/json"
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
var Menu = orders.GetFoods()
var oven, stove = orders.GetApparatus(Menu)
var conf = orders.GetConf()

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
	log.Println(olController.CounterOrdersPickedUp)
	olController.Mutex.Unlock()

	go func() { orders.DistributeFoods(&orderList, cooks, Menu, conf.DiningHallAddress, &olController) }()

	// fmt.Fprint(w, "Order recieved at Kitchen")
	// log.Printf("Order id %v recieved at Kitchen!", order.OrderId)

}

func main() {

	rand.Seed(time.Now().UnixMilli())

	for i := 0; i < len(cooks.Cook); i++ {
		cooks.Cook[i].CookChan = make(chan *orders.CookingDetails, cooks.Cook[i].Proficiency)
		cooks.Cook[i].CondVar = *sync.NewCond(&sync.Mutex{})
		cooks.Cook[i].Queue = make(chan *orders.CookingDetails, 10)
		cooks.Cook[i].CounterAvailable = 0
	}

	r := mux.NewRouter()
	r.HandleFunc("/order", PostDingHallOrders).Methods("POST")

	for i := 0; i < len(cooks.Cook); i++ {
		idx := i
		go cooks.Cook[idx].Work(&orderList, cooks, oven, stove, Menu, conf.DiningHallAddress, &olController)
	}

	http.ListenAndServe(":"+conf.Port, r)

}
