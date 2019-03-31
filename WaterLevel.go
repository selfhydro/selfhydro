package main

import (
	"time"
)

type WaterLevelMeasurer interface {
	GetWaterLevel() (waterLevel float32, time time.Time)
	SetWaterLevel(float32)
	GetWaterLevelFeed() float32
}

type WaterLevel struct {
	waterLevel        float32
	waterLevelChannel chan float32
	time              time.Time
}

func (wl *WaterLevel) SetWaterLevel(level float32) {
	if cap(wl.waterLevelChannel) == 0 {
		wl.waterLevelChannel = make(chan float32, 1)
	}
	wl.waterLevelChannel <- level
	wl.waterLevel = level
	wl.time = time.Now()
}

func (wl *WaterLevel) GetWaterLevel() (waterLevel float32, time time.Time) {
	return wl.waterLevel, wl.time
}

func (wl WaterLevel) GetWaterLevelFeed() float32 {
	if wl.waterLevelChannel == nil {
		return 0
	}
	return <-wl.waterLevelChannel
}
