package main

import (
	"testing"
	"github.com/stianeikeland/go-rpio"
)

var ledState bool

func setupMock() *RaspberryPi {
	mockPi := new(RaspberryPi)
	mockPi.AirPumpPin = new(mockRaspberryPiPinImpl)
	mockPi.GrowLedPin = new(mockRaspberryPiPinImpl)
	return mockPi
}

func TestTurnOnGrowLed(t *testing.T) {
	ledState = false
	mockPi := setupMock()
	mockPi.startLightCycle()
	if mockPi.GrowLedPin.ReadState() != rpio.High {
		t.Errorf("Error: GrowLED not turned on")
	}
}




