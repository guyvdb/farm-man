package command


/*

  Commands are sent between the server and the nodes, possiable via a relay. The commands have
  an opcode and data fields.  

*/

type OpCode int

const (
	OP_NOOP OpCode = iota

	// OP_ACK is used to acknowledge the receipt of a command 
	OP_ACK 

	// The OP_ID_* commands are used to set identification of the nodes
	OP_ID_REQUEST

	
	


	
	// The OP_LOG_* commands are sent from the Node to the server. They include
	// a data point that is to be saved in the time series database.
	OP_LOG_AIR_TEMPRATURE
	OP_LOG_WATER_TEMPRATURE
	OP_LOG_HUMIDITY
  OP_LOG_PH
	OP_LOG_PPM

	OP_LOG_DIO_ON
	OP_LOG_DIO_OFF


	// The OP_DIGITAL_* commands are to controll digital IO pins on the Node. These
	// commands are sent from the server to the node where the node executes the
	// actuation 
	OP_DIGITAL_ON
	OP_DIGITAL_OFF
	OP_DIGITAL_TOGGLE_DURATION
	OP_DIGITAL_SET_SCHEDULE
	OP_DIGITAL_CLEAR_SCHEDULE 


	// The OP_READ_* command is used to force a read of an analog or digital pin. These
	// commands are sent from the server to the node where the node executes the read and
	// responds with the appropriate OP_LOG_* 
	OP_READ_DIGITAL
	OP_READ_ANALOG

	
	


	
	
)


type Command struct {
	From  Id
	To Id
	Via []Id
	OpCode OpCode
	Params []string 
}



func LogAirTemprature() Command {
}
