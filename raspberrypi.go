package main

import (
	"github.com/stianeikeland/go-rpio"
	"log"
)

type RaspberryPiGPIO interface {
	SetPinHigh(pin rpio.Pin)
	SetPinLow(pin rpio.Pin)
}

type RaspberryPi struct {
	GrowLedPin rpio.Pin
	GrowLedState bool
	WaterPumpPin rpio.Pin
	WaterPumpState bool
}

var PinHigh = rpio.Pin.High
var PinLow = rpio.Pin.Low

func (pi RaspberryPi) SetPinHigh(pin rpio.Pin){
	PinHigh(pin)
}

func (pi RaspberryPi) SetPinLow(pin rpio.Pin){
	PinLow(pin)
}

func (pi *RaspberryPi) turnOnGrowLed(){
	pi.SetPinHigh(pi.GrowLedPin)
	pi.GrowLedState = true
}

func (pi *RaspberryPi) turnOffGrowLed(){
	pi.SetPinLow(pi.GrowLedPin)
	pi.GrowLedState = false

}

func (pi *RaspberryPi) turnOffWaterPump(){
	pi.SetPinLow(pi.WaterPumpPin)
	pi.WaterPumpState = false
}

func (pi *RaspberryPi) turnOnWaterPump(){
	log.Printf("Turning on water Pump")
	pi.SetPinHigh(pi.WaterPumpPin)
	pi.WaterPumpState = true

}


