package main

import (
	rpio "github.com/stianeikeland/go-rpio"
)

type AirPump struct {
	pin    RaspberryPiPin
	pumpOn bool
}

func NewAirPump(pin int) *AirPump {
	return &AirPump{
		pin: NewRaspberryPiPin(pin),
	}
}

func (ap *AirPump) TurnOn() {
	ap.pin.WriteState(rpio.High)
	ap.pumpOn = true
}

func (ap *AirPump) TurnOff() {
	ap.pin.WriteState(rpio.Low)
	ap.pumpOn = false
}

func (ap *AirPump) GetState() bool {
	return ap.pumpOn
}

func (ap *AirPump) Setup() error {
	ap.pin.SetMode(rpio.Output)
	return nil
}
