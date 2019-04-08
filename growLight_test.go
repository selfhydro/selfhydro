package main

import (
	"testing"

	rpio "github.com/stianeikeland/go-rpio"
	"github.com/stretchr/testify/assert"
)

func Test_ShouldSetLightOnOnSetup(t *testing.T) {
	mockPin := &mockRaspberryPiPinImpl{}
	mockPin.WriteState(rpio.High)
	growLight := &GrowLight{
		pin: mockPin,
	}
	growLight.Setup()
	assert.True(t, growLight.GetState())
}

func Test_ShouldSetLightStateOffOnSetup(t *testing.T) {
	mockPin := &mockRaspberryPiPinImpl{}
	mockPin.WriteState(rpio.Low)
	growLight := &GrowLight{
		pin: mockPin,
	}
	growLight.Setup()
	assert.False(t, growLight.GetState())
}
