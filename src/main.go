package main

import (
	"fmt"
	"io/ioutil"

	"github.com/gin-gonic/gin"
)

func PostHomePage(c *gin.Context) {

	body := c.Request.Body
	value, err := ioutil.ReadAll(body)
	if err != nil {
		fmt.Println(err.Error())
	}

	c.JSON(200, gin.H{
		"message": string(value),
	})
}

func main() {

	r := gin.Default()
	r.POST("/order", PostHomePage)
	r.Run(":8081")

}
