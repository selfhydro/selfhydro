package main

import (
	"time"

	rpio "github.com/stianeikeland/go-rpio"
)

type Actuator interface {
	TurnOn()
	TurnOff()
	GetState() bool
	Setup() error
}

type WaterPump struct {
	pumpOn bool
	pin    RaspberryPiPin
}

var MAX_SECONDS_ON float64 = 90 * time.Second.Seconds()

func NewWaterPump(pin int) *WaterPump {
	return &WaterPump{
		pin: NewRaspberryPiPin(pin),
	}
}

func (wp *WaterPump) Setup() error {
	wp.pin.SetMode(rpio.Output)
	return nil
}

func (wp *WaterPump) TurnOn() {
	wp.pin.WriteState(rpio.High)
	wp.pumpOn = true
	turnedOn := time.Now()
	go func() {
		for wp.pumpOn {
			if time.Since(turnedOn).Seconds() > MAX_SECONDS_ON {
				wp.TurnOff()
			}
		}
	}()
}

func (wp *WaterPump) TurnOff() {
	wp.pin.WriteState(rpio.Low)
	wp.pumpOn = false
}

func (wp WaterPump) GetState() bool {
	return wp.pumpOn
}
