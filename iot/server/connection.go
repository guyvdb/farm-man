package server

import (
	"fmt"
	"net"
	"time"
	//"math/rand"
)

const (
	MAX_TIMEOUTS = 10
)

var timeoutDuration time.Duration = 100 * time.Millisecond

type Connection struct {
	counter      uint16
	conn         net.Conn
	timeoutCount int
}

func NewConnection(conn net.Conn) *Connection {
	return &Connection{
		counter:      0,
		conn:         conn,
		timeoutCount: 0,
	}
}

func (c *Connection) nextId() uint16 {
	c.counter++
	if c.counter == 0 {
		c.counter = 1
	}
	return c.counter
}

func (c *Connection) CreateTestFrame() *Frame {
	//f := NewFrame([]byte("DATA"))

	f := NewFrame([]byte{})
	f.Id = 343
	f.TCount = 12
	f.Transmitted = 1590917131
	f.RefId = 231
	f.Cmd = 33
	return f
}

func (c *Connection) CreateTestFrameWithControlChars() *Frame {

	f := NewFrame([]byte{0xA, SFLAG, 0xB, EFLAG, 0xC, ESCAPE, 0xD})
	f.Id = ESCAPE
	f.TCount = SFLAG
	f.Transmitted = 1590917131
	f.RefId = EFLAG
	f.Cmd = 33
	return f
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
			request := string(buf[0 : length-1])
			fmt.Printf("Read: %s\n", request)
		}

		// do some writing
		// 		f := c.CreateTestFrame()
		f := c.CreateTestFrameWithControlChars()
		f.Print()
		b := f.NetworkBytes()

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
