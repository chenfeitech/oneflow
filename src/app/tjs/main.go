package main

import (
	"fmt"
	"math/rand"
	"microservice"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	var reply interface{}
	err := microservice.Request("localhost", "Script.Exec", "return 3+5", &reply)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(reply)
}
