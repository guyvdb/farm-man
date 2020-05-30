package server


import (
	"net"
	"time"
	"fmt"
	//"math/rand"
)

const (
	MAX_TIMEOUTS = 10
)

var timeoutDuration time.Duration = 100 * time.Millisecond

type Connection struct {
	conn net.Conn
	timeoutCount int 
}


func NewConnection(conn net.Conn) *Connection {
	return &Connection{
		conn: conn,
		timeoutCount: 0,
	}
}



func (c *Connection) PrintFrame(buf []byte) {

	fmt.Printf("\n--FRAME--\n")

	fmt.Printf("%s\n[", string(buf))
	for _, value := range buf {	
		fmt.Printf("%d ", value)
	}

	fmt.Printf("] %d bytes\n", len(buf))
	
	fmt.Printf("--END FRAME--\n\n")
}


// Handles incoming requests.
func (c *Connection) Run() {



	var length int
	var rerr error
	var werr error

	
	//frameno := 0
	buf := make([]byte, 1024)
	
	
	for {
		//Error reading: read tcp 192.168.8.100:3000->192.168.8.100:44672: i/o timeout


		c.conn.SetReadDeadline(time.Now().Add(timeoutDuration))
	
		
		// do some reading 
		length, rerr = c.conn.Read(buf)
		if rerr != nil {
			if err, ok := rerr.(net.Error); ok && err.Timeout() {
				fmt.Printf("Timeout Error: %v\n", err)
				c.timeoutCount++
			} else {
				fmt.Printf("Read Error: %v\n", rerr)
				break
			}
		} else {
			c.timeoutCount = 0
			request := string(buf[0:length-1])
			fmt.Printf("Read: %s\n", request)
		}


		// do some writing
		f := NewFrame('A',[]byte("DATA"))
		b := f.Bytes()
		c.PrintFrame(b)

		// write out the frame 
		length, werr = c.conn.Write(b)
		if werr != nil {
			fmt.Printf("Write Error: %v\n", werr)
		} else {
			fmt.Printf("Wrote %d bytes.\n", length)
			c.timeoutCount = 0 
		}

		time.Sleep(10000 * time.Millisecond)

		if c.timeoutCount > MAX_TIMEOUTS {
			fmt.Printf("Max timeouts reached. Closing connection.\n")
			c.conn.Close()
			return
		}

		
		
	}
  c.conn.Close()
}
