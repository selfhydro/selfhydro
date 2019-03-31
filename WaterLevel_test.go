package main

import (
	"testing"

	"gotest.tools/assert"
)

func TestShouldReturnWaterLevel(t *testing.T) {
	waterLevel := WaterLevel{waterLevel: float32(2.3)}
	actualWaterLevel, _ := waterLevel.GetWaterLevel()
	assert.Equal(t, actualWaterLevel, float32(2.3))
}

func Test_ShouldSetWaterLevel(t *testing.T) {
	waterLevel := WaterLevel{}
	waterLevel.SetWaterLevel(float32(10))
	assert.Equal(t, <-waterLevel.waterLevelChannel, float32(10))
}

func Test_ShouldGetWaterLevelFromChannelFeed(t *testing.T) {
	waterLevel := WaterLevel{
		waterLevelChannel: make(chan float32, 1),
	}
	waterLevel.SetWaterLevel(float32(10))
	assert.Equal(t, waterLevel.GetWaterLevelFeed(), float32(10))
}
