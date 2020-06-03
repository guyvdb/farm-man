package server




const (
	ACK uint8 = 0         // can be used to respond in the affirmative to a number of commands 
	NACK = 1              // can be used to respond in the negative to a number of commands
	// RESERVED = 2
	// RESERVED = 3
	IDENT = 4             // sent by client to server - no response  
	TIMEREQ = 5           // sent by client to server - respond with timeresp
	TIMERESP = 6          // sent by server to client - in response to timereq 
	// RESERVED = 27 
	
)

func CmdCreateACK(refid uint16) *Frame {
	frame := NewFrame([]byte{})
	frame.Cmd = ACK
	frame.RefId = refid
	return frame
}

func CmdCreateNACK(refid uint16) *Frame {
	frame := NewFrame([]byte{})
	frame.Cmd = NACK
	frame.RefId = refid
	return frame 
}

// Create a frame to request and ident 
/*func CmdCreateIDENTREQ() *Frame {
	frame := NewFrame([]byte{})
	frame.Cmd = IDENTREQ
	return frame 
}*/

// Extract the ident out of a ident response 
func CmdDecodeIDENT(frame *Frame) (uint32, error) {
	return frame.GetUint32Arg()
}

// Create a time response for the time reques 
func CmdCreateTIMERESP(refid uint16, time int64) *Frame {
	args := NewFramePayload()
	args.AddUint32(uint32(time))

	frame := NewFrame(args.Bytes())
	frame.Cmd = TIMERESP
	frame.RefId = refid
	return frame 
}
