package main

import "github.com/stianeikeland/go-rpio"

type RaspberryPiPin interface {
	ReadState() rpio.State
	WriteState(state rpio.State)
	SetMode(mode rpio.Mode)
	Frequency(freq int)
	DutyCycle(dutyLen, cycleLen uint32)
	Toggle()
}

type raspberryPiPinImpl struct {
	rpioPin rpio.Pin
}

func NewRaspberryPiPin(pin int) RaspberryPiPin {
	return &raspberryPiPinImpl{
		rpioPin: rpio.Pin(pin),
	}
}

func (r *raspberryPiPinImpl) ReadState() rpio.State {
	return r.rpioPin.Read()
}

func (r *raspberryPiPinImpl) WriteState(state rpio.State) {
	r.rpioPin.Write(state)
}

func (r *raspberryPiPinImpl) SetMode(mode rpio.Mode) {
	r.rpioPin.Mode(mode)
}

func (r *raspberryPiPinImpl) Frequency(freq int) {
	r.rpioPin.Freq(freq)
}

func (r *raspberryPiPinImpl) DutyCycle(dutyLen, cycleLen uint32) {
	r.rpioPin.DutyCycle(dutyLen, cycleLen)
}

func (r *raspberryPiPinImpl) Toggle() {
	r.rpioPin.Toggle()
}
