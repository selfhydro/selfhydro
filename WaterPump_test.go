package main

import (
	"testing"
	"time"

	rpio "github.com/stianeikeland/go-rpio"
	"gotest.tools/assert"
)

func Test_ShouldSetupWaterPump(t *testing.T) {
	waterPump := WaterPump{
		pin: &mockRaspberryPiPinImpl{},
	}
	waterPump.Setup()
	assert.Equal(t, rpio.Output, waterPump.pin.(*mockRaspberryPiPinImpl).ModeOfPin)
}

func Test_ShouldTurnOnWaterPump(t *testing.T) {
	waterPump := WaterPump{
		pin: &mockRaspberryPiPinImpl{},
	}
	waterPump.TurnOn()
	assert.Equal(t, rpio.High, waterPump.pin.ReadState())
}

func Test_ShouldTurnOffWaterPump(t *testing.T) {
	waterPump := WaterPump{
		pin: &mockRaspberryPiPinImpl{},
	}
	waterPump.TurnOn()
	waterPump.TurnOff()
	assert.Equal(t, rpio.Low, waterPump.pin.ReadState())
}

func Test_ShouldTurnOffWaterPumpIfItHasBeenOnForMaxTime(t *testing.T) {
	waterPump := WaterPump{
		pin: &mockRaspberryPiPinImpl{},
	}
	MAX_SECONDS_ON = float64(.1)
	waterPump.TurnOn()
	time.Sleep(time.Second)
	assert.Equal(t, rpio.Low, waterPump.pin.ReadState())
}
