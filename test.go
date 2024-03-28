package main

import (
	"fmt"

	"github.com/vantu-fit/master-go-be/utils"
)

func Test() {
	config, err := utils.LoadConfig(".")

	if err!=nil {
		fmt.Println(err)
	}

	fmt.Println(config.BDSource)
}