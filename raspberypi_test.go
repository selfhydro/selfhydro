package main

import (
	"testing"
	"github.com/stianeikeland/go-rpio"
	"os"
)

var ledState bool

func TestMain(m *testing.M){
	OldPinHigh := PinHigh
	defer func () {PinHigh = OldPinHigh}()
	PinHigh = func (pin rpio.Pin) {
		ledState = true
	}
	code := m.Run()
	os.Exit(code)

}

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



