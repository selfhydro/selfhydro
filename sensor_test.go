package main

import (
	"testing"
	"github.com/stianeikeland/go-rpio"
)

func TestGetState(t *testing.T) {
	sensor := new(sensor)
	mockPin := new(mockRaspberryPiPinImpl)
	sensor.pin = mockPin

	mockPin.stateOfPin = rpio.High

	if sensor.getState() != rpio.High {
		t.Error("Should be returning high")
	}


}
