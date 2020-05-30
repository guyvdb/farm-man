package main

import (
	"fmt"
	iot "github.com/guyvdb/farm-man/iot/server"
)


func main() {
	fmt.Printf("Farm-Man-Server\n")

	s := iot.NewServer()
	s.Run()


}
