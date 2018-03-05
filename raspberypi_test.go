package main

import (
	"testing"
	"github.com/stianeikeland/go-rpio"
)

func setupMock() *RaspberryPi {
	mockPi := new(RaspberryPi)
	mockPi.AirPumpPin = new(mockRaspberryPiPinImpl)
	mockPi.GrowLedPin = new(mockRaspberryPiPinImpl)
	return mockPi
}

func TestHydroCycle(t *testing.T) {
	mockPi := setupMock()
	t.Run("Testing Grow LEDS", func(t *testing.T) {
		mockPi.startLightCycle()
		if mockPi.GrowLedPin.ReadState() != rpio.High {
			t.Errorf("Error: GrowLED not turned on")
		}
	})
	
	t.Run("Test Air Pump cycle", func(t *testing.T) {
		mockPi.startAirPumpCycle()
		if mockPi.AirPumpPin.ReadState() != rpio.High {
			t.Errorf("Error: Airpump was not turned on")
		}
	})
}

