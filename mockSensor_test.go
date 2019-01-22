package main

import "github.com/stianeikeland/go-rpio"

type mockSensor struct {
	sensorState rpio.State
}

func (ms mockSensor) getState() rpio.State {
	return ms.sensorState
}
