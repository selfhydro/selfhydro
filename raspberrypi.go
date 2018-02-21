package main

import (
	"github.com/stianeikeland/go-rpio"
	"log"
	"time"
	"github.com/d2r2/go-dht"
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
	WaterTempSensor ds18b20
	WaterLevelSensor rpio.Pin
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

func (pi RaspberryPi) getWaterTemp(){
	pi.WaterTempSensor.ReadTemperature()
}

func (pi RaspberryPi) StartWaterCycle() {
	go func() {
		for {
			if pi.WaterLevelSensor.Read() != rpio.High {
				log.Printf("ALERT: Water level is low")
			}
			pi.turnOnWaterPump()
			time.Sleep(time.Second * 5)
			pi.turnOffWaterPump()
			time.Sleep(time.Minute * 150)
		}
	}()
}
func (pi RaspberryPi) StartLightCycle() {
	turnOnTime, _ := time.Parse("15:04:05", "04:45:00")
	turnOffTime, _ := time.Parse("15:04:05", "23:45:00")
	go func() {
		for {
			if !pi.GrowLedState && betweenTime(turnOnTime, turnOffTime) {
				log.Printf("Turning on GROW LEDS")
				pi.turnOnGrowLed()
			} else if pi.GrowLedState && betweenTime(turnOffTime, turnOnTime.Add(time.Hour*24)) {
				log.Printf("Turning off GROW LEDS")
				pi.turnOffGrowLed()
			}
			time.Sleep(time.Minute * 1)
		}

	}()
}
func (pi RaspberryPi) StartSensorCycle() {

	go func() {
		for {
			temperature, humidity, retried, err :=
				dht.ReadDHTxxWithRetry(dht.DHT22, 17, true, 10)
			if err != nil {
				log.Printf("Error: Error with reading dht: %s", err.Error())
			}
			log.Printf("Ambient Temperature = %v*C, Humidity = %v%% (retried %d times)\n",
				temperature, humidity, retried)
			pi.getWaterTemp()
			time.Sleep(time.Hour)
		}

	}()
}


func betweenTime(startTime time.Time, endTime time.Time) bool {
	currentTimeString := time.Now().Format("15:04:05")
	currentTime, _ := time.Parse("15:04:05", currentTimeString)
	if currentTime.After(startTime) && currentTime.Before(endTime) {
		return true
	}

	return false
}

