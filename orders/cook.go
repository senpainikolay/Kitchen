package orders

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type Cook struct {
	Id          int    `json:"id"`
	Rank        int    `json:"rank"`
	Proficiency int    `json:"proficiency"`
	Name        string `json:"name"`
	CatchPhrase string `json:"catch_phrase"`
}

type Cooks struct {
	Cook []Cook `json:"cooks"`
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
