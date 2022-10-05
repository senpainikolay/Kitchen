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
var cooks = orders.GetCooks()

func PostDingHallOrders(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	var order orders.Order
	err := json.NewDecoder(r.Body).Decode(&order)
	if err != nil {
		log.Fatalln("There was an error decoding the request body into the struct")
	}

	orderList.Append(&order)

	fmt.Fprint(w, "Order recieved at Kitchen")
	log.Printf("Order id %v recieved at Kitchen!", order.OrderId)

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

	go func() {
		for i := 0; i < len(cooks.Cook); i++ {
			idx := i
			go cooks.Cook[idx].Work(&orderList, cooks)
		}

	}()

	http.ListenAndServe(":8081", r)

}
