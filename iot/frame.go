package iot

import (
	"bytes"
	"errors"
	"fmt"
	"time"
)

const (
	FRAMEVERSION = 0x1
)

const (
	ESCAPE = 0x1B
	SFLAG  = 0x2
	EFLAG  = 0x3
)

var ControlCharacters []byte = []byte{
	SFLAG,
	EFLAG,
	ESCAPE,
}

var (
	ArgumentOverflow error = errors.New("ArgumentOverflow")
)

/* ------------------------------------------------------------------------
 *
 * --------------------------------------------------------------------- */
// You can use a FramePayload to generate the payload byte array
// of a Frame before you create the Frame
type FramePayload struct {
	buffer *bytes.Buffer
}

//
// Byte No: 0             1            2   3            4             5  6  7  8           9  10           11          12           13..n
// |--------------|---------------|------------|---------------|---------------------|---------------|------------|-----------|--------------------|--------------|
// | SFLAG 1 byte | version 1byte | id 2 bytes | tcount 1 byte | transmitted 4 bytes | refid 2 bytes | cmd 1 byte | len 1byte | payload <varbytes> | EFLAG 1 byte |
// |--------------|---------------|------------|---------------|---------------------|---------------|------------|-----------|--------------------|--------------|
//
// The SFLAG and EFLAG are not part of the frame.
// The frame is 16 bytes + payload size
// The transmitted bytes are 18 bytes + payload size
// The meaning of payload is individually defined by each command
//
// The frame is in network byte order when encoded. Big endian.
//
type Frame struct {
	Version     uint8
	Id          uint16
	TCount      uint8
	Transmitted int32
	RefId       uint16
	Cmd         uint8
	Len         uint8
	Payload     []byte
	argPtr      uint8
}

/* ------------------------------------------------------------------------
 *
 * --------------------------------------------------------------------- */
func NewFramePayload() *FramePayload {
	return &FramePayload{
		buffer: &bytes.Buffer{},
	}
}

/* ------------------------------------------------------------------------
 *
 * --------------------------------------------------------------------- */
func (payload *FramePayload) AddUint8(data uint8) {
	payload.buffer.WriteByte(byte(data))
}

/* ------------------------------------------------------------------------
 *
 * --------------------------------------------------------------------- */
func (payload *FramePayload) AddUint16(data uint16) {
	payload.buffer.Write(bytes_uint16_encode(data))
}

/* ------------------------------------------------------------------------
 *
 * --------------------------------------------------------------------- */
func (payload *FramePayload) AddUint32(data uint32) {
	payload.buffer.Write(bytes_uint32_encode(data))
}

/* ------------------------------------------------------------------------
 *
 * --------------------------------------------------------------------- */
func (payload *FramePayload) AddString(s string) {
	if len(s) > 254 {
		fmt.Printf("Cannot encode string. Too long.\n")
		return
	}

	length := uint8(len(s))

	payload.AddUint8(length)

	for _, v := range s {
		payload.AddUint8(uint8(v))
	}

	// encode the string length as uint8
	// write the string bytes as uint8
}

/* ------------------------------------------------------------------------
 *
 * --------------------------------------------------------------------- */
func (payload *FramePayload) Bytes() []byte {
	return payload.buffer.Bytes()
}

/* ------------------------------------------------------------------------
 *
 * --------------------------------------------------------------------- */
// create a new frame with a payload
func NewFrame(payload []byte) *Frame {
	return &Frame{
		Version: FRAMEVERSION,
		TCount:  1,
		//		Transmitted: int32(time.Now().Unix()),
		RefId:   0,
		Len:     uint8(len(payload)),
		Payload: payload,
		argPtr:  0,
	}
}

// Byte No: 0             1            2   3            4             5  6  7  8           9  10           11          12           13.. n
// |--------------|---------------|------------|---------------|---------------------|---------------|------------|-----------|--------------------|--------------|
// | SFLAG 1 byte | version 1byte | id 2 bytes | tcount 1 byte | transmitted 4 bytes | refid 2 bytes | cmd 1 byte | len 1byte | payload <varbytes> | EFLAG 1 byte |
// |--------------|---------------|------------|---------------|---------------------|---------------|------------|-----------|--------------------|--------------|
// Create a frame from network order escaped bytes
func NewFrameFromNetworkBytes(data []byte) *Frame {
	var f Frame
	buf := stripControlCharacters(data)

	f.Version = bytes_uint8_decode(buf[0:])
	f.Id = bytes_uint16_decode(buf[1:])
	f.TCount = bytes_uint8_decode(buf[3:])
	f.Transmitted = int32(bytes_uint32_decode(buf[4:]))
	f.RefId = bytes_uint16_decode(buf[8:])
	f.Cmd = bytes_uint8_decode(buf[10:])
	f.Len = bytes_uint8_decode(buf[11:])
	f.argPtr = 0

	if f.Len > 0 {
		f.Payload = make([]byte, int(f.Len))
		for idx, d := range buf[12:] {
			f.Payload[idx] = d
		}
	} else {
		f.Payload = []byte{}
	}

	return &f

}

/* ------------------------------------------------------------------------
 *
 * --------------------------------------------------------------------- */
// Check if a character is a control character and thus
// should be escaped
func isControlCharacter(data uint8) bool {
	for _, c := range ControlCharacters {
		if data == c {
			return true
		}
	}
	return false
}

/* ------------------------------------------------------------------------
 *
 * --------------------------------------------------------------------- */
// Given a set of bytes that has just been pulled off the
// wire, return a set of bytes that have had the network
// control characters stripped out
func stripControlCharacters(data []byte) []byte {
	var buf bytes.Buffer

	escaped := false

	for index, b := range data {
		// if this is the first or last byte it is the SFLAG or EFLAG
		if index != 0 && index != len(data)-1 {
			if escaped {
				// if the escaped flag is set then this byte is an escaped char
				buf.WriteByte(b)
				escaped = false
			} else {
				// if b is an escape character then enter escaped mode else write the byte
				if b == ESCAPE {
					escaped = true
				} else {
					buf.WriteByte(b)
				}
			}
		}
	}

	return buf.Bytes()
}

/* ------------------------------------------------------------------------
 *
 * --------------------------------------------------------------------- */
func (f *Frame) TransmittedNow() {
	f.Transmitted = int32(time.Now().Unix())
}

/* ------------------------------------------------------------------------
 *
 * --------------------------------------------------------------------- */
func (f *Frame) String() string {
	var buf bytes.Buffer

	for idx, c := range f.Payload {
		if idx+1 < int(f.Len) {
			buf.Write([]byte(fmt.Sprintf("%d ", c)))
		} else {
			buf.Write([]byte(fmt.Sprintf("%d", c)))
		}
	}

	// if it is a log command append the LOGTYPE
	if f.Cmd == LOG {

		// the logtype is the first byte of the params buffer
		if f.Len < 1 {
			return "ERROR: cmd is log but no log type specified"
		} else {
			logtype := LogTypeToString(f.Payload[0])

			return fmt.Sprintf("%s (%s) { v: %d, id: %d, tc: %d, time: %d, ref: %d, cmd: %d, len: %d, args: [%s] }",
				CmdToString(f.Cmd),
				logtype,
				f.Version,
				f.Id,
				f.TCount,
				f.Transmitted,
				f.RefId,
				f.Cmd,
				f.Len,
				string(buf.Bytes()),
			)
		}
	} else {
		return fmt.Sprintf("%s { v: %d, id: %d, tc: %d, time: %d, ref: %d, cmd: %d, len: %d, args: [%s] }",
			CmdToString(f.Cmd),
			f.Version,
			f.Id,
			f.TCount,
			f.Transmitted,
			f.RefId,
			f.Cmd,
			f.Len,
			string(buf.Bytes()),
		)
	}

}

/* ------------------------------------------------------------------------
 *
 * --------------------------------------------------------------------- */
// Print this frame to stdout
func (f *Frame) Print() {
	buf := f.NetworkBytes() // the network representation of this frame

	fmt.Printf("\n-- FRAME --\n")
	fmt.Printf("version: %d\n", f.Version)
	fmt.Printf("id: %d\n", f.Id)
	fmt.Printf("tcount: %d\n", f.TCount)
	fmt.Printf("transmitted: %d\n", f.Transmitted)
	fmt.Printf("refid: %d\n", f.RefId)
	fmt.Printf("cmd: %d\n", f.Cmd)
	fmt.Printf("len: %d\n", f.Len)

	if len(f.Payload) > 0 {
		fmt.Printf("Payload:\n  [")
		for i := 0; i < len(f.Payload); i++ {
			fmt.Printf("%d ", f.Payload[i])
		}
		fmt.Printf("]\n")
	} else {
		fmt.Printf("Payload:\n  [NULL]\n")
	}

	fmt.Printf("Network Bytes:\n  [")
	for _, b := range buf {
		fmt.Printf("%d ", b)
	}
	fmt.Printf("] Size=%d\n", len(buf))
	fmt.Printf("-- END FRAME --\n\n")
}

/* ------------------------------------------------------------------------
 *
 * --------------------------------------------------------------------- */
// Return an escaped network order representation of the frame with control characters
func (f *Frame) NetworkBytes() []byte {
	var buf bytes.Buffer

	// add start flag
	buf.WriteByte(SFLAG)

	// write data, escaping where needed
	for _, b := range f.FrameBytes() {
		if isControlCharacter(b) {
			buf.WriteByte(ESCAPE)
			buf.WriteByte(byte(b))
		} else {
			buf.WriteByte(byte(b))
		}
	}

	// add end flag
	buf.WriteByte(EFLAG)

	return buf.Bytes()
}

/* ------------------------------------------------------------------------
 *
 * --------------------------------------------------------------------- */
// Return a network order representation of the frame without any control characters
func (f *Frame) FrameBytes() []byte {
	var buf bytes.Buffer

	// Version
	buf.WriteByte(byte(f.Version))

	// Id
	buf.Write(bytes_uint16_encode(f.Id))

	// TCount
	buf.WriteByte(byte(f.TCount))

	// Transmitted
	buf.Write(bytes_uint32_encode(uint32(f.Transmitted)))

	// RefId
	buf.Write(bytes_uint16_encode(f.RefId))

	// Cmd
	buf.WriteByte(byte(f.Cmd))

	// Len
	buf.WriteByte(byte(f.Len))

	// Add the payload
	if len(f.Payload) > 0 {
		buf.Write(f.Payload)
	}

	return buf.Bytes()
}

/* ------------------------------------------------------------------------
 *
 * --------------------------------------------------------------------- */
// Check that we can still pull len bytes from the payload
func (f *Frame) checkArgs(len uint8) bool {
	if f.argPtr+len < f.Len {
		return true
	} else {
		return false
	}
}

/* ------------------------------------------------------------------------
 *
 * --------------------------------------------------------------------- */
func (f *Frame) GetUint8Arg() (uint8, error) {
	if f.checkArgs(1) {
		val := bytes_uint8_decode(f.Payload[f.argPtr:])
		f.argPtr += 1
		return val, nil
	} else {
		return 0, ArgumentOverflow
	}
}

/* ------------------------------------------------------------------------
 *
 * --------------------------------------------------------------------- */
func (f *Frame) GetUint16Arg() (uint16, error) {
	if f.checkArgs(2) {
		val := bytes_uint16_decode(f.Payload[f.argPtr:])
		f.argPtr += 2
		return val, nil
	} else {
		return 0, ArgumentOverflow
	}
}

/* ------------------------------------------------------------------------
 *
 * --------------------------------------------------------------------- */
func (f *Frame) GetUint32Arg() (uint32, error) {
	if f.checkArgs(1) {
		val := bytes_uint32_decode(f.Payload[f.argPtr:])
		f.argPtr += 4
		return val, nil
	} else {
		return 0, ArgumentOverflow
	}
}
