package simulator


import (
	"fmt"
)


type Simulator struct {
}


func New() *Simulator {

	return &Simulator{
	}
}


func (s *Simulator) Hello() {
	fmt.Printf("Hello from simulator\n")
}



func (s *Simluator) Start() {
}


func (s *Simulator) Shutdown() {
}


