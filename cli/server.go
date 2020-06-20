package cli

// this server listens for commands from the cli client and transmits them to the appropriate mote

import (
	"log"
	"net"
	"sync"
	"time"

	"github.com/guyvdb/farm-man/iot"
)

type Server struct {
	ch        chan bool
	waitGroup *sync.WaitGroup
	control   iot.ControlPanel
}

/* ------------------------------------------------------------------------
 *
 * --------------------------------------------------------------------- */
func NewServer(control iot.ControlPanel) *Server {
	return &Server{
		ch:        make(chan bool),
		waitGroup: &sync.WaitGroup{},
		control:   control,
	}
}

/* ------------------------------------------------------------------------
 *
 * --------------------------------------------------------------------- */
func (s *Server) Serve(listener *net.TCPListener) {
	//defer s.waitGroup.Done()
	for {
		select {
		case <-s.ch:
			log.Println("[CLI] Stopping listening on", listener.Addr())
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
		log.Printf("[CLI] %s Connected.\n", conn.RemoteAddr())
		s.waitGroup.Add(1)

		// need to process new conn here
		c := NewConnection(conn, s.control)
		go c.Run()

	}
}

/* ------------------------------------------------------------------------
 *
 * --------------------------------------------------------------------- */
func (s *Server) Stop() {
	close(s.ch)
	//s.waitGroup.Done() // for the listener
	s.waitGroup.Wait()
	log.Println("[CLI] Server shutdown complete.")
}
