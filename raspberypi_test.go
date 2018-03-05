package main

import (
	"testing"
	"github.com/stianeikeland/go-rpio"
	"time"
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
		startTimeString := time.Now().Add(-time.Minute).Format("15:04:05")
		startTime, _ := time.Parse("15:04:05", startTimeString)

		offTimeString := time.Now().Add(time.Minute).Format("15:04:05")
		offTime, _ := time.Parse("15:04:05", offTimeString)

		mockPi.changeLEDState(startTime, offTime)
		if mockPi.GrowLedPin.ReadState() != rpio.High {
			t.Errorf("Error: GrowLED not turned on")
		}
	})
	
	t.Run("Test Air Pump cycle", func(t *testing.T) {
		mockPi.airPumpCycle(time.Second, time.Second)
		if mockPi.AirPumpPin.ReadState() != rpio.Low {
			t.Errorf("Error: Airpump was not turned on")
		}
	})
}

