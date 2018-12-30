package main

import (
	"testing"

	"gotest.tools/assert"
)

func TestShouldReturnWaterLevel(t *testing.T) {
	waterLevel := WaterLevel{waterLevel: float32(2.3)}
	actualWaterLevel := waterLevel.GetWaterLevel()
	assert.Equal(t, actualWaterLevel, float32(2.3))
}
