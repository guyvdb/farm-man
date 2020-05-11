package sensors

/*

The DHT11 measures relative humidity. Relative humidity is the amount of water 
vapor in air vs. the saturation point of water vapor in air. At the saturation 
point, water vapor starts to condense and accumulate on surfaces forming dew.

The saturation point changes with air temperature. Cold air can hold less water 
vapor before it becomes saturated, and hot air can hold more water vapor before 
it becomes saturated.

The formula to calculate relative humidity is:

RH = Pw/Ps

RH is relative humidity
Pw is density of water vapor
Ps is density of water vapor at saturation 

*/



type HumiditySensor struct {
}
