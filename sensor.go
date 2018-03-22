package main

import "github.com/stianeikeland/go-rpio"

type Sensor interface {
	getState() rpio.State
}

type sensor struct {
	pin RaspberryPiPin
}

func NewSensor(pin int) *sensor {
	sensor := new(sensor)
	sensor.pin = NewRaspberryPiPin(pin)
	sensor.pin.SetMode(rpio.Input)
	return sensor
}

func (s *sensor) getState() rpio.State {
	return s.pin.ReadState()
}
