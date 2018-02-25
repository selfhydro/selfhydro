package main

import (
	"github.com/stianeikeland/go-rpio"
	"log"
	"time"
)

type Controller interface {
	setPinHigh(pin rpio.Pin)
	setPinLow(pin rpio.Pin)

	startSensorCycle()
	startLightCycle()
	startWaterCycle()
	startAirPumpCycle()
}

type RaspberryPi struct {
	GrowLedPin       rpio.Pin
	GrowLedState     bool
	WaterPumpPin     rpio.Pin
	WaterPumpState   bool
	WaterTempSensor  ds18b20
	WaterLevelSensor rpio.Pin
	AirPump          rpio.Pin
}

var PinHigh = rpio.Pin.High
var PinLow = rpio.Pin.Low

func NewRaspberryPi() *RaspberryPi {
	pi := new(RaspberryPi)

	pi.GrowLedPin = rpio.Pin(19)
	pi.GrowLedState = false

	pi.WaterPumpPin = rpio.Pin(20)
	pi.WaterPumpState = false

	pi.WaterLevelSensor = rpio.Pin(4)
	pi.WaterLevelSensor.Input()

	pi.GrowLedPin.Mode(rpio.Output)
	pi.WaterPumpPin.Mode(rpio.Output)

	pi.WaterTempSensor.id = "28-0316838ca7ff"

	pi.AirPump = rpio.Pin(21)
	pi.AirPump.Mode(rpio.Output)

	return pi
}

func (pi *RaspberryPi) StartHydroponics() {
	pi.startSensorCycle()
	pi.startLightCycle()
	pi.startWaterCycle()
	pi.startAirPumpCycle()

}

func (pi *RaspberryPi) StopSystem() {
	pi.turnOffGrowLed()
	pi.WaterPumpPin.Low()
	pi.AirPump.Low()
}

func (pi RaspberryPi) setPinHigh(pin rpio.Pin) {
	PinHigh(pin)
}

func (pi RaspberryPi) setPinLow(pin rpio.Pin) {
	PinLow(pin)
}

func (pi *RaspberryPi) turnOnGrowLed() {
	pi.setPinHigh(pi.GrowLedPin)
	pi.GrowLedState = true
}

func (pi *RaspberryPi) turnOffGrowLed() {
	pi.setPinLow(pi.GrowLedPin)
	pi.GrowLedState = false

}

func (pi *RaspberryPi) turnOffWaterPump() {
	pi.setPinLow(pi.WaterPumpPin)
	pi.WaterPumpState = false
}

func (pi *RaspberryPi) turnOnWaterPump() {
	log.Printf("Turning on water Pump")
	pi.setPinHigh(pi.WaterPumpPin)
	pi.WaterPumpState = true

}

func (pi RaspberryPi) getWaterTemp() {
	pi.WaterTempSensor.ReadTemperature()
}

func (pi RaspberryPi) startWaterCycle() {
	go func() {
		for {
			//if pi.WaterLevelSensor.Read() != rpio.High {
			//	log.Printf("ALERT: Water level is low")
			//}
			pi.turnOnWaterPump()
			time.Sleep(time.Second * 1)
			pi.turnOffWaterPump()
			time.Sleep(time.Hour * 8)
		}
	}()
}
func (pi RaspberryPi) startLightCycle() {
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
			time.Sleep(time.Second * 4)
		}

	}()
}
func (pi RaspberryPi) startSensorCycle() {

	go func() {
		dht22 := NewDHT22(17)
		for {
			temperature, tErr := dht22.Temperature()
			humidity, hErr := dht22.Humidity()
			//temperature, humidity, retried, err :=
			//	dht.ReadDHTxxWithRetry(dht.DHT22, 17, true, 10)
			if tErr != nil {
				log.Printf("Error: Error with reading dht: %s", tErr.Error())
			}
			if hErr != nil {
				log.Printf("Error: Error with reading dht: %s", tErr.Error())
			}
			log.Printf("Ambient Temperature = %v*C, Humidity = %v%% \n",
				temperature, humidity)
			pi.getWaterTemp()
			time.Sleep(time.Second * 5)
		}

	}()
}

func (pi RaspberryPi) startAirPumpCycle() {
	go func() {
		for {
			log.Printf("Turning on air pump")
			pi.AirPump.High()
			time.Sleep(time.Minute * 30)
			log.Printf("Turning off air pump")
			pi.AirPump.Low()
			time.Sleep(time.Hour * 3)
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
