package server

import (
	"fmt"
	"net"
	"time"
	"math/rand"
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


func (c *Connection) CreateRollingFrame() *Frame {
	f := NewFrame([]byte{0xA,0xB,0xC,0xD})
	f.Id = c.nextId( )
	f.TCount = 1
	f.Transmitted = int32(time.Now().Unix())
	f.RefId = 0
	f.Cmd = 33
	return f
}

func (c *Connection) WriteFrame(frame *Frame) {

	var length int
	var werr error

	fmt.Printf("=======================> WRITE\n")

	
	b := frame.NetworkBytes()
	
	length, werr = c.conn.Write(b)
	if werr != nil {
		fmt.Printf("Write Error: %v\n", werr)
	} else {
		fmt.Printf("Wrote %d bytes.\n", length)
		c.timeoutCount = 0
	}
}

// we are going to write the bytes in two writes with a delay between them
// this is only for testing fragment combination on the client 
func (c *Connection) WriteFrameFragmented(frame *Frame) {
	var length int
	var werr error


	fmt.Printf("=======================> FRAGMENTED WRITE\n")

	b := frame.NetworkBytes()


	// we want a range in 4..n-4 bytes
  // the frame is a minimum of 18 bytes long


	min := 4
	max := len(b) - 4
	split := rand.Intn(max - min + 1) + min

	ba := b[0:split]
	bb := b[split:]

	
	length, werr = c.conn.Write(ba)
	if werr != nil {
		fmt.Printf("Write Error: %v\n", werr)
	} else {
		fmt.Printf("Wrote %d bytes.\n", length)
	}

	time.Sleep(2000 * time.Millisecond)

	length, werr = c.conn.Write(bb)
	if werr != nil {
		fmt.Printf("Write Error: %v\n", werr)
	} else {
		fmt.Printf("Wrote %d bytes.\n", length)
		c.timeoutCount = 0
	}

	
}

func (c *Connection) shouldFragment() bool {

	//	return false 

	
	
	max := float64(0.8)
	val := rand.Float64()

	if val >= max {
		return true
	}

	return false 

}

// Handles incoming requests.
func (c *Connection) Run() {
	var length int
	var rerr error
	//var werr error

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
		//		f := c.CreateTestFrameWithControlChars()
		f := c.CreateRollingFrame()
		f.Print()


		if(c.shouldFragment()) {
			c.WriteFrameFragmented(f)
		} else {
			c.WriteFrame(f)
		}

		/*
		b := f.NetworkBytes()

		// write out the frame
		length, werr = c.conn.Write(b)
		if werr != nil {
			fmt.Printf("Write Error: %v\n", werr)
		} else {
			fmt.Printf("Wrote %d bytes.\n", length)
			c.timeoutCount = 0
		}


    */

		time.Sleep(10000 * time.Millisecond)

		if c.timeoutCount > MAX_TIMEOUTS {
			fmt.Printf("Max timeouts reached. Closing connection.\n")
			c.conn.Close()
			return
		}

	}
	c.conn.Close()
}
