package main

import rpio "github.com/stianeikeland/go-rpio"

type GrowLight struct {
	lightsOn bool
	pin      RaspberryPiPin
}

func NewGrowLight(pin int) *GrowLight {
	return &GrowLight{
		pin: NewRaspberryPiPin(pin),
	}
}

func (gl *GrowLight) TurnOn() {
	gl.pin.WriteState(rpio.High)
}

func (gl *GrowLight) TurnOff() {
	gl.pin.WriteState(rpio.Low)
}

func (gl *GrowLight) GetState() bool {
	return gl.lightsOn
}

func (gl *GrowLight) Setup() error {
	gl.pin.SetMode(rpio.Output)
	return nil
}
