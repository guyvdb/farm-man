package service




/*

A node on a network registers a number of services to a registry. These services are either actuators or sensors services 

*/



type ServiceType int

const (
	SERVICE_SWITCH_GPIO
	SERVICE_TOGGLE_GPIO
	

)
