package main

import (
	"fmt"
	iot "github.com/guyvdb/farm-man/iot/server"

	"math/rand"
	"time"
)


func main() {
	fmt.Printf("Farm-Man-Server\n")


	// prime the random number generator
	rand.Seed(time.Now().UnixNano())

	


	// run the server 
	s := iot.NewServer()
	s.Run()


}
