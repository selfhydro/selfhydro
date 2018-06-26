package main

import "github.com/stianeikeland/go-rpio"

type mockRaspberryPiPinImpl struct {
	stateOfPin rpio.State
	readStateCalled bool
	writeStateCalled bool
}


func (r *mockRaspberryPiPinImpl) ReadState() rpio.State {
	return r.stateOfPin
}

func (r *mockRaspberryPiPinImpl) WriteState(state rpio.State) {
	r.stateOfPin = state
	r.writeStateCalled = true
}

func (r *mockRaspberryPiPinImpl) SetMode(mode rpio.Mode) {

}

func (r *mockRaspberryPiPinImpl) Frequency(freq int){
}

func (r *mockRaspberryPiPinImpl) DutyCycle(dutyLen, cycleLen uint32) {
}

func (r *mockRaspberryPiPinImpl) Toggle() {
}

