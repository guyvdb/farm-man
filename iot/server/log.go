package server

const (
	LOGTYPE_POWERUP            uint8 = 0
	LOGTYPE_REBOOT                   = 1
	LOGTYPE_TIMESET                  = 2
	LOGTYPE_RELAY_STATE_CHANGE       = 3
	LOGTYPE_SENSOR                   = 4
)

var LOGTYPE map[uint8]string = map[uint8]string{
	LOGTYPE_POWERUP:            "POWERUP",
	LOGTYPE_REBOOT:             "REBOOT",
	LOGTYPE_TIMESET:            "TIMESET",
	LOGTYPE_RELAY_STATE_CHANGE: "RELAY_STATE_CHANGE",
	LOGTYPE_SENSOR:             "SENSOR",
}

func LogTypeToString(logtype uint8) string {
	v, ok := LOGTYPE[logtype]
	if ok {
		return v
	} else {
		return "LOGTYPE_UNKNOWN"
	}
}

/* ------------------------------------------------------------------------
 * We return an ACK or a NACK on receiving a log
 * --------------------------------------------------------------------- */
func LogSave(moteid uint32, frame *Frame) *Frame {
	return CmdCreateACK(frame.Id)
}
