package iot

import (
	//	"fmt"
	"log"
	"net"
	"sync"
	"time"
)

/* ------------------------------------------------------------------------
 *
 * --------------------------------------------------------------------- */
type IOTServer struct {
	counter     int64
	ch          chan bool
	waitGroup   *sync.WaitGroup
	connections map[int64]Connection
	motes       map[uint32]int64
}

/* ------------------------------------------------------------------------
 *
 * --------------------------------------------------------------------- */
func NewServer() Server {
	s := &IOTServer{
		counter:     0,
		ch:          make(chan bool),
		waitGroup:   &sync.WaitGroup{},
		connections: make(map[int64]Connection),
		motes:       make(map[uint32]int64),
	}
	s.waitGroup.Add(1)
	return s
}

/* ------------------------------------------------------------------------
 *
 * --------------------------------------------------------------------- */
func (s *IOTServer) nextConnectionId() int64 {
	s.counter++
	return s.counter
}

/* ------------------------------------------------------------------------
 *
 * --------------------------------------------------------------------- */
func (s *IOTServer) Serve(listener *net.TCPListener) {
	defer s.waitGroup.Done()
	for {
		select {
		case <-s.ch:
			log.Println("[IOT] Stopping listening on", listener.Addr())
			listener.Close()
			return
		default:
		}
		listener.SetDeadline(time.Now().Add(1e9))
		conn, err := listener.AcceptTCP()
		if nil != err {
			if opErr, ok := err.(*net.OpError); ok && opErr.Timeout() {
				continue
			}
			log.Println(err)
		}
		//		log.Println(conn.RemoteAddr(), "Connected")
		log.Printf("[IOT] %s Connected.\n", conn.RemoteAddr())
		s.waitGroup.Add(1)

		c := NewConnection(s, conn, s.nextConnectionId())
		s.connections[c.GetId()] = c
		go c.Read()
		go c.Process()
	}
}

/* ------------------------------------------------------------------------
 *
 * --------------------------------------------------------------------- */
func (s *IOTServer) Stop() {
	for _, conn := range s.connections {
		conn.Stop()
	}

	close(s.ch)
	s.waitGroup.Done() // for the listener
	s.waitGroup.Wait()
	log.Println("[IOT] Server shutdown complete.")
}

/* ------------------------------------------------------------------------
 *
 * --------------------------------------------------------------------- */
func (s *IOTServer) RegisterMote(connection Connection) {
	s.motes[connection.GetMoteId()] = connection.GetId()
	log.Printf("[IOT] Mote %d registered\n", connection.GetMoteId())
}

/* ------------------------------------------------------------------------
 *
 * --------------------------------------------------------------------- */
func (s *IOTServer) ConnectionComplete(connection Connection) {
	s.waitGroup.Done()
	_, ok := s.connections[connection.GetId()]
	if ok {
		delete(s.connections, connection.GetId())
	}
}

/* ------------------------------------------------------------------------
 *
 * --------------------------------------------------------------------- */
func (s *IOTServer) MoteConnected(moteid uint32) bool {
	_, ok := s.motes[moteid]
	return ok
}

/* ------------------------------------------------------------------------
 *
 * --------------------------------------------------------------------- */
func (s *IOTServer) SetRelay(moteid uint32, pin uint8, value uint8) bool {
	mid, ok := s.motes[moteid]
	if !ok {
		return false
	}

	conn, cok := s.connections[mid]
	if !cok {
		return false
	}

	conn.SetRelay(pin, value)
	return true
}
