package iot

import (
//	"fmt"
)

const (
	ACK         uint8 = 0
	NACK              = 1
	IDENTREQ          = 4
	IDENT             = 5
	TIMEREQ           = 6
	TIMESET           = 7
	TIMEZONESET       = 8
	LOG               = 9
	RELAYSET          = 10
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
	RELAYSET:    "RELAYSET",
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

/* ------------------------------------------------------------------------
 *
 * --------------------------------------------------------------------- */
func CmdRELAYSET(pin uint8, value uint8) *Frame {
	args := NewFramePayload()
	args.AddUint8(pin)
	args.AddUint8(value)

	frame := NewFrame(args.Bytes())
	frame.Cmd = RELAYSET
	return frame
}
