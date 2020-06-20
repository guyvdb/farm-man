package cli

import (
	"fmt"
	"github.com/guyvdb/farm-man/iot"
	"net"
	"strings"
)

// A connection by a cli client
type Connection struct {
	conn    net.Conn
	running bool
	control iot.ControlPanel
}

func NewConnection(conn net.Conn, control iot.ControlPanel) *Connection {
	return &Connection{
		conn:    conn,
		running: true,
		control: control,
	}
}

func (c *Connection) Run() {
	var length int
	var rerr error
	//	var cmd string

	buf := make([]byte, 1024)

	for {
		//c.conn.SetReadDeadline(time.Now().Add(readTimeoutDuration))
		length, rerr = c.conn.Read(buf)

		if rerr != nil {
			if err, ok := rerr.(net.Error); ok && err.Timeout() {
				// read timeout do nothing
			} else {
				// read error ... disconnect
				c.running = false
			}
		} else {
			cmd := string(buf[0 : length-1])
			parts := strings.Split(cmd, " ")

			// mote 55 relay off 1
			if parts[2] == "relay" {
				if parts[3] == "off" {
					fmt.Printf("mote 55 turn off relay\n")
					c.control.SetRelay(55, 4, 0)
				} else if parts[3] == "on" {
					fmt.Printf("mote 55 turnon relay\n")
					c.control.SetRelay(55, 4, 1)
				}
			}

			fmt.Printf("[CLI] Received: %s\n", string(buf[0:length-1]))
			c.conn.Write(buf[0:length])
		}

		if !c.running {
			return
		}
	}
}
