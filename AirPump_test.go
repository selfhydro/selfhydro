package main

import (
	"testing"

	rpio "github.com/stianeikeland/go-rpio"
	"gotest.tools/assert"
)

func Test_ShouldTurnOnAirPump(t *testing.T) {
	wp := &AirPump{
		pin: &mockRaspberryPiPinImpl{},
	}
	wp.TurnOn()
	assert.Equal(t, rpio.High, wp.pin.(*mockRaspberryPiPinImpl).ReadState())
	assert.Equal(t, true, wp.GetState())
}

func Test_ShouldTurnOffAirPump(t *testing.T) {
	wp := &AirPump{
		pin: &mockRaspberryPiPinImpl{},
	}
	wp.TurnOff()
	assert.Equal(t, rpio.Low, wp.pin.(*mockRaspberryPiPinImpl).ReadState())
	assert.Equal(t, false, wp.GetState())
}

func Test_ShouldSetupAirPump(t *testing.T) {
	ap := AirPump{
		pin: &mockRaspberryPiPinImpl{},
	}
	ap.Setup()
	assert.Equal(t, rpio.Output, ap.pin.(*mockRaspberryPiPinImpl).ModeOfPin)
}
