package server

import (
//	"fmt"
)

const (
	ACK  uint8 = 0 // can be used to respond in the affirmative to a number of commands
	NACK       = 1 // can be used to respond in the negative to a number of commands
	// RESERVED = 2
	// RESERVED = 3
	IDENTREQ    = 4
	IDENT       = 5 // sent by client to server - no response
	TIMEREQ     = 6 // sent by client to server - respond with timeresp
	TIMESET     = 7 // sent by server to client - in response to timereq - client should set time
	TIMEZONESET = 8
	LOG         = 9
	// RESERVED = 27

)

var COMMANDS map[uint8]string = map[uint8]string{
	ACK:         "ACK",
	NACK:        "NACK",
	IDENTREQ:    "IDENTREQ",
	IDENT:       "IDENT",
	TIMEREQ:     "TIMEREQ",
	TIMESET:     "TIMESET",
	TIMEZONESET: "TIMEZONESET",
	LOG:         "LOG",
}

/* ------------------------------------------------------------------------
 *
 * --------------------------------------------------------------------- */
func CmdToString(cmd uint8) string {
	v, ok := COMMANDS[cmd]
	if ok {
		return v
	} else {
		return "UNKNOWN"
	}
}

/* ------------------------------------------------------------------------
 *
 * --------------------------------------------------------------------- */
func CmdCreateACK(refid uint16) *Frame {
	frame := NewFrame([]byte{})
	frame.Cmd = ACK
	frame.RefId = refid
	return frame
}

/* ------------------------------------------------------------------------
 *
 * --------------------------------------------------------------------- */
func CmdCreateNACK(refid uint16) *Frame {
	frame := NewFrame([]byte{})
	frame.Cmd = NACK
	frame.RefId = refid
	return frame
}

/* ------------------------------------------------------------------------
 *
 * --------------------------------------------------------------------- */
// Extract the ident out of a ident response
func CmdDecodeIDENT(frame *Frame) (uint32, error) {
	return frame.GetUint32Arg()
}

/* ------------------------------------------------------------------------
 *
 * --------------------------------------------------------------------- */
// Create a time response for the time request
// TIMESET args = |timestamp 4 byte|strlen 1 btye|len bytes timezone|
func CmdCreateTIMESET(refid uint16, time int64) *Frame {
	args := NewFramePayload()
	args.AddUint32(uint32(time))
	args.AddString("SAST-2") // TODO look up tz string from config

	frame := NewFrame(args.Bytes())
	frame.Cmd = TIMESET
	frame.RefId = refid
	return frame
}
