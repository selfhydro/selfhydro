package main

import (
	"time"

	"github.com/stianeikeland/go-rpio"
)

type UltrasonicSensor interface {
	MeasureDistance() (cm float32)
}

type HCSR04 struct {
	echoPin RaspberryPiPin
	pingPin RaspberryPiPin
}

const HardStop = 1000000

func NewHCSR04Sensor(pingPin int, echoPin int) UltrasonicSensor {

	hcsr04 := new(HCSR04)
	hcsr04.pingPin = NewRaspberryPiPin(pingPin)
	hcsr04.echoPin = NewRaspberryPiPin(echoPin)

	return hcsr04
}

func (hcsr04 *HCSR04) MeasureDistance() (cm float32) {
	hcsr04.initPins()

	strobeZero := 0
	strobeOne := 0

	delayUs(200)
	hcsr04.pingPin.WriteState(rpio.High)
	delayUs(15)
	hcsr04.pingPin.WriteState(rpio.Low)

	for strobeZero = 0; strobeZero < HardStop && hcsr04.echoPin.ReadState() != rpio.High; strobeZero++ {
	}
	startTime := time.Now()
	for strobeOne = 0; strobeOne < HardStop && hcsr04.echoPin.ReadState() != rpio.Low; strobeOne++ {
		delayUs(1)
	}
	endTime := time.Now()

	return float32(endTime.UnixNano()-startTime.UnixNano()) / (58.0 * 1000)
}

func (hcsr04 *HCSR04) initPins() {
	hcsr04.echoPin.SetMode(rpio.Output)
	hcsr04.pingPin.SetMode(rpio.Output)
	hcsr04.echoPin.WriteState(rpio.Low)
	hcsr04.pingPin.WriteState(rpio.Low)
	time.Sleep(time.Microsecond)
	hcsr04.echoPin.SetMode(rpio.Input)
}

func delayUs(ms int) {
	time.Sleep(time.Duration(ms) * time.Microsecond)
}
