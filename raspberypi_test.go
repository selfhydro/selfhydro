package main

import (
	"testing"
)

var ledState bool

func setupMock() *RaspberryPi {
	mockPi := new(RaspberryPi)
	mockPi.WaterPumpPin = new(mockRaspberryPiPinImpl)
	mockPi.AirPump = new(mockRaspberryPiPinImpl)
	mockPi.GrowLedPin = new(mockRaspberryPiPinImpl)
	return mockPi
}

func TestTurnOnGrowLed(t *testing.T) {
	ledState = false
	mockPi := setupMock()
	mockPi.turnOnGrowLed()
	if !ledState {
		t.Errorf("Error: GrowLED not turned on")
	}
}

func TestTurnOffGrowLed(t *testing.T) {
	ledState = false
	mockPi := setupMock()
	mockPi.turnOnGrowLed()
	if !ledState {
		t.Errorf("Error: GrowLED not turned on")
	}
}
func TestTurnOnWaterPump(t *testing.T) {
	ledState = false
	mockPi := setupMock()
	mockPi.turnOnGrowLed()
	if !ledState {
		t.Errorf("Error: GrowLED not turned on")
	}
}
func TestTurnOffWaterPump(t *testing.T) {
	ledState = false
	mockPi := setupMock()
	mockPi.turnOnGrowLed()
	if !ledState {
		t.Errorf("Error: GrowLED not turned on")
	}
}



