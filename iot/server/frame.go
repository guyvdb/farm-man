package server


import (
	"bytes"
)

const  (
	ESCAPE = 0x1B
	SFLAG = 0x2
	EFLAG = 0x3
)

var escapable []byte = []byte{
	SFLAG,
	EFLAG,
}


type Frame struct {
	Command byte
	Data []byte
}

func NewFrame(cmd byte, data []byte) *Frame {
	return &Frame{
		Command: cmd,
		Data: data,
	}
}




// Return the data in an escaped format 
func (f *Frame) Bytes() []byte {
	var buf bytes.Buffer


	// add the START FLAG
	buf.WriteByte(SFLAG)


	// write the data buffer, escaping if needed
	for _, b := range f.Data {
		for _, e := range escapable {
			if b == e {
				buf.WriteByte(ESCAPE)
				break
			}
		}
		buf.WriteByte(b)
	}

	// add the END FLAG
	buf.WriteByte(EFLAG)

	return buf.Bytes()
}
