package sensors



type OperatingMethod int;

const (
	ANALOG OperatingMethod = iota
	DIGITAL
)






type Sensor struct {
	OperatingMethod OperatingMethod
	Pins []int 
	
}
