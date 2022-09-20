package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/gin-gonic/gin"
	"github.com/senpainikolay/Kitchen/orders"
)

//global
var orderList = orders.OrderList{}
var menu = orders.GetFoods()
var cooks = orders.GetCooks()

func PostHomePage(c *gin.Context) {

	body := c.Request.Body
	value, err := ioutil.ReadAll(body)
	if err != nil {
		fmt.Println(err.Error())
	}
	var order orders.Order
	json.Unmarshal(value, &order)
	order.Wg.Add(len(order.Items))
	orderList.Append(&order)
	fmt.Printf("Recieved:  %+v", order)
}

func prepare() {

	order := orderList.PickUp()

	order.Wg.Wait()

}

func main() {

	r := gin.Default()
	r.POST("/order", PostHomePage)
	go r.Run(":8081")
	fmt.Printf("Recieved:  %+v", cooks.Cook[0])
	/*
		for {
			prepare()
			time.Sleep(4000 * time.Millisecond)
		}
	*/

}
