package main

import (
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	iot "github.com/guyvdb/farm-man/iot/server"
)

const (
	//CONN = "192.168.0.99:6000"
	CONN = "192.168.8.100:6000"
)

func startIotServer() *iot.IoTServer {
	ladd, err := net.ResolveTCPAddr("tcp", CONN)
	if err != nil {
		log.Fatalln(err)
	}

	listener, err := net.ListenTCP("tcp", ladd)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("Listening on ", listener.Addr())

	moteServer := iot.NewServer()
	go moteServer.Serve(listener)

	return moteServer
}

func main() {
	fmt.Printf("Farm Manager\n")
	fmt.Printf("--------------\n")

	// prime the random number generator
	rand.Seed(time.Now().UnixNano())

	moteServer := startIotServer()

	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	log.Println("Signal: ", <-ch)

	moteServer.Stop()

}
