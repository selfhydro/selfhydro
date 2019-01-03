package main

import "time"

type WaterLevelMeasurer interface {
	GetWaterLevel() float32
	SetWaterLevel(float32)
}

type WaterLevel struct {
	waterLevel float32
	time       time.Time
}

func (wl *WaterLevel) GetWaterLevel() float32 {
	return wl.waterLevel
}
