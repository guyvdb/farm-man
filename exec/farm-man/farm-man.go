package main

import (
	"fmt"
	"github.com/guyvdb/farm-man/cli"
	"github.com/guyvdb/farm-man/config"
	"github.com/guyvdb/farm-man/iot"
	//"github.com/guyvdb/farm-man/platform/adapter"
	"github.com/guyvdb/farm-man/platform/adapter/mongo"
	"github.com/guyvdb/farm-man/platform/domain"
	"github.com/guyvdb/farm-man/platform/model/infrastructure"
	"log"
	"math/rand"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

/* ------------------------------------------------------------------------
 *
 * --------------------------------------------------------------------- */
func startIotServer(address string) iot.Server {
	ladd, err := net.ResolveTCPAddr("tcp", address)
	if err != nil {
		log.Fatalln(err)
	}

	listener, err := net.ListenTCP("tcp", ladd)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("[IOT] Listening on ", listener.Addr())

	moteServer := iot.NewServer()
	go moteServer.Serve(listener)
	return moteServer
}

/* ------------------------------------------------------------------------
 *
 * --------------------------------------------------------------------- */
func startCliServer(address string, control iot.ControlPanel) *cli.Server {
	ladd, err := net.ResolveTCPAddr("tcp", address)
	if err != nil {
		log.Fatalln(err)
	}

	listener, err := net.ListenTCP("tcp", ladd)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("[CLI] Listening on ", listener.Addr())

	cliServer := cli.NewServer(control)
	go cliServer.Serve(listener)

	return cliServer
}

/* ------------------------------------------------------------------------
 *
 * --------------------------------------------------------------------- */
func main() {
	fmt.Printf("Farm Manager\n")
	fmt.Printf("--------------\n")

	cfg, err := config.NewConfig("config")
	if err != nil {
		panic(err)
	}

	adapter, err := mongo.NewMongoAdapter(cfg.DatabaseURI, cfg.DatabaseName)
	if err != nil {
		panic(err)
	}

	service := domain.NewDomainService(adapter)

	fmt.Printf("service: %v\n", service)

	// func (s *Service) CreateBounds(btype infrastructure.BoundsType, name string, parent *infrastructure.Bounds) (*infrastructure.Bounds, error) {
	farm, err := service.CreateBounds(infrastructure.BOUNDSTYPE_FARM, "Farm", nil)
	if err != nil {
		panic(err)
	}

	t1, err := service.CreateBounds(infrastructure.BOUNDSTYPE_GREENHOUSE, "Tunnel 1", farm)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Create: farm = %v, tunnel = %v\n", farm, t1)

	// ==============
	return

	// prime the random number generator
	rand.Seed(time.Now().UnixNano())

	// Start the variouse server processes
	moteServer := startIotServer(cfg.IOTAddress)
	cliServer := startCliServer(cfg.CLIAddress, moteServer)

	// Block for SIGTERM
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	log.Println("Signal: ", <-ch)

	// Stop the various servers
	cliServer.Stop()
	moteServer.Stop()

	// done

}
