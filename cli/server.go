package cli

// this server listens for commands from the cli client and transmits them to the appropriate mote

import (
	"net"
	"sync"
)

type Server struct {
	ch        chan bool
	waitGroup *sync.WaitGroup
}

/* ------------------------------------------------------------------------
 *
 * --------------------------------------------------------------------- */
func NewServer() *Server {
	return &Server{
		ch:        make(chan bool),
		waitGroup: &sync.WaitGroup{},
	}
}

/* ------------------------------------------------------------------------
 *
 * --------------------------------------------------------------------- */
func (s *Server) Serve(listener *net.TCPListener) {
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

		// need to process new conn here

	}
}

/* ------------------------------------------------------------------------
 *
 * --------------------------------------------------------------------- */
func (s *Server) Stop() {
	close(s.ch)
	s.waitGroup.Done() // for the listener
	s.waitGroup.Wait()
	log.Println("Server shutdown complete.")
}
