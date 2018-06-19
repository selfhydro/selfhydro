package main

import (
	"github.com/stianeikeland/go-rpio"
	"time"
)

type HCSR04 struct {

	echoPin RaspberryPiPin
	pingPin RaspberryPiPin
}

func NewHCSR04Sensor(pingPin int,echoPin int) *HCSR04 {

	hcsr04 := new(HCSR04)
	hcsr04.pingPin = NewRaspberryPiPin(pingPin)
	hcsr04.echoPin = NewRaspberryPiPin(echoPin)

	return hcsr04
}

func (hcsr04 *HCSR04) MeasureDistance() float32 {
	hcsr04.echoPin.SetMode(rpio.Input)
	hcsr04.pingPin.SetMode(rpio.Output)
	hcsr04.pingPin.WriteState(rpio.Low)
	time.Sleep(time.Microsecond)
	hcsr04.pingPin.WriteState(rpio.High)
	time.Sleep(time.Microsecond*15)
	hcsr04.pingPin.WriteState(rpio.Low)




	return 0
}
