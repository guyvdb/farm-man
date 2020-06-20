package iot

// this file defines the interface to the iot server and the active motes

import (
	"net"
	//"sync"
)

/* ------------------------------------------------------------------------
 *
 * --------------------------------------------------------------------- */
type Server interface {
	ControlPanel
	Serve(listener *net.TCPListener)
	RegisterMote(connection Connection)
	ConnectionComplete(connection Connection)
	Stop()
}

/* ------------------------------------------------------------------------
 *
 * --------------------------------------------------------------------- */
type ControlPanel interface {
	MoteConnected(moteid uint32) bool
	SetRelay(moteid uint32, pin uint8, value uint8) bool
}

/* ------------------------------------------------------------------------
 *
 * --------------------------------------------------------------------- */
type Connection interface {
	GetId() int64
	GetMoteId() uint32
	Read()
	Process()
	Stop()
	SetRelay(pin uint8, value uint8)
}
