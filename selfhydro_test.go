package main

import (
	"testing"

	"gotest.tools/assert"
)

func TestShouldGetAmbientTemp(t *testing.T) {
	sh := selfhydro{}
	ambientTemp := sh.GetAmbientTemp()
	assert.Equal(t, float32(10), ambientTemp)
}

func TestShouldGetWaterLevel(t *testing.T) {
	sh := selfhydro{}
	sh.waterLevel.waterLevel = float32(2.24)
	waterLevel := sh.GetWaterLevel()
	assert.Equal(t, float32(2.24), waterLevel)
}
