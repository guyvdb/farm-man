package iot

import (
	"fmt"
	"log"
	"net"
	"sync"
	"time"
)

/* ------------------------------------------------------------------------
 *
 * --------------------------------------------------------------------- */
type IoTServer struct {
	counter     int64
	ch          chan bool
	waitGroup   *sync.WaitGroup
	connections map[int64]*Connection
	motes       map[uint32]int64
}

/* ------------------------------------------------------------------------
 *
 * --------------------------------------------------------------------- */
func NewServer() *IoTServer {
	s := &IoTServer{
		counter:     0,
		ch:          make(chan bool),
		waitGroup:   &sync.WaitGroup{},
		connections: make(map[int64]*Connection),
		motes:       make(map[uint32]int64),
	}
	s.waitGroup.Add(1)
	return s
}

/* ------------------------------------------------------------------------
 *
 * --------------------------------------------------------------------- */
func (s *IoTServer) nextConnectionId() int64 {
	s.counter++
	return s.counter
}

/* ------------------------------------------------------------------------
 *
 * --------------------------------------------------------------------- */
func (s *IoTServer) Serve(listener *net.TCPListener) {
	defer s.waitGroup.Done()
	for {
		select {
		case <-s.ch:
			log.Println("Stopping listening on", listener.Addr())
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
		log.Println(conn.RemoteAddr(), "Connected")
		s.waitGroup.Add(1)

		c := NewConnection(s, conn, s.nextConnectionId())
		s.connections[c.Id] = c
		go c.Read()
		go c.Process()
	}
}

/* ------------------------------------------------------------------------
 *
 * --------------------------------------------------------------------- */
func (s *IoTServer) Stop() {
	for _, conn := range s.connections {
		conn.Stop()
	}

	close(s.ch)
	s.waitGroup.Done() // for the listener
	s.waitGroup.Wait()
	log.Println("Server shutdown complete.")
}

/* ------------------------------------------------------------------------
 *
 * --------------------------------------------------------------------- */
func (s *IoTServer) RegisterMote(connection *Connection) {
	s.motes[connection.MoteId] = connection.Id
	fmt.Printf("Mote: %d registered\n", connection.MoteId)
}

/* ------------------------------------------------------------------------
 *
 * --------------------------------------------------------------------- */
func (s *IoTServer) ConnectionComplete(connection *Connection) {
	s.waitGroup.Done()
	_, ok := s.connections[connection.Id]
	if ok {
		delete(s.connections, connection.Id)
	}
}
