package iot

import (
	"log"
	"net"
	"time"
)

const (
	MAX_TIMEOUTS = 10
)

var timeoutDuration time.Duration = 100 * time.Millisecond
var readTimeoutDuration time.Duration = 1 * time.Millisecond
var idleTimeoutDuration time.Duration = 60 * time.Second

/* ------------------------------------------------------------------------
 *
 * --------------------------------------------------------------------- */
type Connection struct {
	Id                  int64
	MoteId              uint32
	server              *IoTServer
	counter             uint16
	conn                net.Conn
	timeoutCount        int
	readTimeoutDuration int
	running             bool
	debug               bool
	frames              chan *Frame
}

/* ------------------------------------------------------------------------
 *
 * --------------------------------------------------------------------- */
func NewConnection(server *IoTServer, conn net.Conn, connid int64) *Connection {
	return &Connection{
		Id:           connid,
		server:       server,
		MoteId:       0,
		counter:      0,
		conn:         conn,
		timeoutCount: 0,
		running:      true,
		debug:        true,
		frames:       make(chan *Frame, 10),
	}
}

/* ------------------------------------------------------------------------
 *
 * --------------------------------------------------------------------- */
func (c *Connection) nextId() uint16 {
	c.counter++
	if c.counter == 0 {
		c.counter = 1
	}
	return c.counter
}

/* ------------------------------------------------------------------------
 *
 * --------------------------------------------------------------------- */
func (c *Connection) Write(frame *Frame) {
	var werr error
	frame.Id = c.nextId()

	b := frame.NetworkBytes()

	_, werr = c.conn.Write(b)
	if werr != nil {
		log.Printf("[%d:%d] Write Error: %v\n", c.Id, c.MoteId, werr)
	} else {
		if c.debug {
			log.Printf("[%d:%d] TX %s\n", c.Id, c.MoteId, frame)
		}
		c.timeoutCount = 0
	}
}

/* ------------------------------------------------------------------------
 *
 * --------------------------------------------------------------------- */
func (c *Connection) Read() {
	var length int
	var rerr error

	buf := make([]byte, 1024)

	for {
		c.conn.SetReadDeadline(time.Now().Add(readTimeoutDuration))
		length, rerr = c.conn.Read(buf)

		if rerr != nil {
			if err, ok := rerr.(net.Error); ok && err.Timeout() {
				// read timeout do nothing
			} else {
				// read error ... disconnect
				c.running = false
			}
		} else {
			// TODO need to handle frame fragmentation
			c.frames <- NewFrameFromNetworkBytes(buf[0:length])
		}

		if !c.running {
			return
		}
	}
}

/* ------------------------------------------------------------------------
 *
 * --------------------------------------------------------------------- */
func (c *Connection) Process() {
	var err error
	defer c.conn.Close()

	for {
		f := <-c.frames
		if c.debug {
			log.Printf("[%d:%d] RX %s\n", c.Id, c.MoteId, f)
		}
		switch f.Cmd {
		case IDENT:
			c.MoteId, err = CmdDecodeIDENT(f)
			if err != nil {
				log.Printf("[%d:%d] [ERROR] decoding mote id from IDENT. %v\n", c.Id, c.MoteId, err)
			} else {
				c.server.RegisterMote(c)
				c.Write(CmdCreateACK(f.Id))
			}
		case TIMEREQ:
			c.Write(CmdCreateTIMESET(f.Id, time.Now().Unix()))
		case LOG:
			c.Write(LogSave(c.MoteId, f))
		default:
			// do nothing
		}
		if !c.running {
			c.server.ConnectionComplete(c)
			return
		}
		time.Sleep(timeoutDuration)
	}
}

/* ------------------------------------------------------------------------
 *
 * --------------------------------------------------------------------- */
func (c *Connection) Stop() {
	c.running = false
}
