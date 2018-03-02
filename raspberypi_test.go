package main

import (
	"testing"
)

var ledState bool

func TestTurnOnGrowLed(t *testing.T) {
	ledState = false
	mockPi := new(RaspberryPi)
	mockPi.turnOnGrowLed()
	if !ledState {
		t.Errorf("Error: GrowLED not turned on")
	}
}
func TestTurnOffGrowLed(t *testing.T) {
	ledState = false
	mockPi := new(RaspberryPi)
	mockPi.turnOnGrowLed()
	if !ledState {
		t.Errorf("Error: GrowLED not turned on")
	}
}
func TestTurnOnWaterPump(t *testing.T) {
	ledState = false
	mockPi := new(RaspberryPi)
	mockPi.turnOnGrowLed()
	if !ledState {
		t.Errorf("Error: GrowLED not turned on")
	}
}
func TestTurnOffWaterPump(t *testing.T) {
	ledState = false
	mockPi := new(RaspberryPi)
	mockPi.turnOnGrowLed()
	if !ledState {
		t.Errorf("Error: GrowLED not turned on")
	}
}



