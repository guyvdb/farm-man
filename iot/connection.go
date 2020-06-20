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
type IOTConnection struct {
	Id                  int64
	MoteId              uint32
	server              Server
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
func NewConnection(server Server, conn net.Conn, connid int64) Connection {
	return &IOTConnection{
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
func (c *IOTConnection) GetId() int64 {
	return c.Id
}

/* ------------------------------------------------------------------------
 *
 * --------------------------------------------------------------------- */
func (c *IOTConnection) GetMoteId() uint32 {
	return c.MoteId
}

/* ------------------------------------------------------------------------
 *
 * --------------------------------------------------------------------- */
func (c *IOTConnection) nextId() uint16 {
	c.counter++
	if c.counter == 0 {
		c.counter = 1
	}
	return c.counter
}

/* ------------------------------------------------------------------------
 *
 * --------------------------------------------------------------------- */
func (c *IOTConnection) Write(frame *Frame) {
	var werr error
	frame.Id = c.nextId()

	b := frame.NetworkBytes()

	_, werr = c.conn.Write(b)
	if werr != nil {
		log.Printf("[IOT] [%d:%d] Write Error: %v\n", c.Id, c.MoteId, werr)
	} else {
		if c.debug {
			log.Printf("[IOT] [%d:%d] TX %s\n", c.Id, c.MoteId, frame)
		}
		c.timeoutCount = 0
	}
}

/* ------------------------------------------------------------------------
 *
 * --------------------------------------------------------------------- */
func (c *IOTConnection) Read() {
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
func (c *IOTConnection) Process() {
	var err error
	defer c.conn.Close()

	for {
		f := <-c.frames
		if c.debug {
			log.Printf("[IOT] [%d:%d] RX %s\n", c.Id, c.MoteId, f)
		}
		switch f.Cmd {
		case IDENT:
			c.MoteId, err = CmdDecodeIDENT(f)
			if err != nil {
				log.Printf("[IOT] [%d:%d] [ERROR] decoding mote id from IDENT. %v\n", c.Id, c.MoteId, err)
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
func (c *IOTConnection) Stop() {
	c.running = false
}

/* ------------------------------------------------------------------------
 *
 * --------------------------------------------------------------------- */
func (c *IOTConnection) SetRelay(pin uint8, value uint8) {
	c.Write(CmdRELAYSET(pin, value))
}
