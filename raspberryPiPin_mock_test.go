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
