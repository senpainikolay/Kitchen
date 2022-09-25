package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/senpainikolay/Kitchen/orders"
)

//global
var orderList = orders.OrderList{}
var cooks = orders.GetCooks()

func PostHomePage(c *gin.Context) {

	body := c.Request.Body
	value, err := ioutil.ReadAll(body)
	if err != nil {
		fmt.Println(err.Error())
	}
	var order orders.Order
	json.Unmarshal(value, &order)
	orderList.Append(&order)
	fmt.Printf("Recieved:  %+v", order)
}
func main() {
	rand.Seed(time.Now().UnixNano())

	for i := 0; i < len(cooks.Cook); i++ {
		cooks.Cook[i].CookChan = make(chan *orders.CookingDetails)
	}

	r := gin.Default()
	r.POST("/order", PostHomePage)
	// fmt.Printf("Recieved:  %+v", cooks.Cook)

	fmt.Printf("KEKEKEKE %+v\n", orders.Menu.Foods[1])
	go func() {
		time.Sleep(3 * time.Second)
		for i := 0; i < len(cooks.Cook); i++ {
			idx := i
			go cooks.Cook[idx].Work(&orderList, cooks)

		}

	}()
	r.Run(":8081")

}
