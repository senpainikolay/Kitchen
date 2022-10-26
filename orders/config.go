package orders

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type Conf struct {
	Port              string `json:"port"`
	DiningHallAddress string `json:"dining_hall_address"`
	NR_OF_STOVES      int    `json:"nr_of_stoves"`
	NR_OF_OVENS       int    `json:"nr_of_ovens"`
}

func GetConf() *Conf {
	jsonFile, err := os.Open("configurations/Conf.json")
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	var conf Conf
	json.Unmarshal(byteValue, &conf)
	return &conf

}
