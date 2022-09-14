package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/senpainikolay/Kitchen/orders"

	"github.com/gin-gonic/gin"
)

//global
var orderList = orders.OrderList{}
var menu = orders.GetFoods()

func PostHomePage(c *gin.Context) {

	body := c.Request.Body
	value, err := ioutil.ReadAll(body)
	if err != nil {
		fmt.Println(err.Error())
	}
	var order orders.Order
	json.Unmarshal(value, &order)
	orderList.Append(&order)

	c.JSON(200, gin.H{
		"message": fmt.Sprintf(" Order %d recived at Kitchen", order.OrderId),
		"order":   string(value),
	})
}

func prepare() {
	if !orderList.IsEmpty() {
		order := orderList.PickUp()
		postOrder := orders.Order{OrderId: order.OrderId, Priority: order.Priority, MaxWait: order.MaxWait}
		for i := range order.Items {
			temp := order.Items[i]
			go func() {
				time.Sleep(time.Duration(menu.Foods[temp-1].PreparationTime) * 50 * time.Millisecond)
				postOrder.Items = append(postOrder.Items, temp)
			}()
		}
		for {
			time.Sleep(50 * time.Millisecond)
			if len(order.Items) == len(postOrder.Items) {
				break
			}
		}

		postBody, _ := json.Marshal(postOrder)
		responseBody := bytes.NewBuffer(postBody)
		resp, err := http.Post("http://dininghall:8080/distribution", "application/json", responseBody)
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

}

func main() {

	r := gin.Default()
	r.POST("/order", PostHomePage)
	go r.Run(":8081")
	for {
		prepare()
		time.Sleep(4000 * time.Millisecond)
	}

}
