package main

import (
	"github.com/stianeikeland/go-rpio"
	"log"
	"time"
	"github.com/d2r2/go-dht"
)

type Controller interface {
	StopSystem()
	StartHydroponics()

	startSensorCycle()
	startLightCycle()
	startWaterCycle()
	startAirPumpCycle()

	setPinHigh(pin RaspberryPiPin)
	setPinLow(pin RaspberryPiPin)
}

type RaspberryPi struct {
	GrowLedPin       RaspberryPiPin
	GrowLedState     bool
	WaterPumpPin     RaspberryPiPin
	WaterPumpState   bool
	WaterTempSensor  ds18b20
	AirPump          RaspberryPiPin
}


func NewRaspberryPi() *RaspberryPi {
	pi := new(RaspberryPi)

	rpio.Open()
	//defer rpio.Close()

	pi.GrowLedPin = NewRaspberryPiPin(19)
	pi.GrowLedState = false

	pi.WaterPumpPin = NewRaspberryPiPin(20)
	pi.WaterPumpState = false

	pi.GrowLedPin.SetMode(rpio.Output)
	pi.WaterPumpPin.SetMode(rpio.Output)

	pi.WaterTempSensor.id = "28-0316838ca7ff"

	pi.AirPump = NewRaspberryPiPin(21)
	pi.AirPump.SetMode(rpio.Output)

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
	pi.WaterPumpPin.WriteState(rpio.Low)
	pi.AirPump.WriteState(rpio.Low)
}

func (pi *RaspberryPi) turnOnGrowLed() {
	pi.GrowLedPin.WriteState(rpio.High)
	pi.GrowLedState = true
}

func (pi *RaspberryPi) turnOffGrowLed() {
	pi.GrowLedPin.WriteState(rpio.Low)
	pi.GrowLedState = false

}

func (pi *RaspberryPi) turnOffWaterPump() {
	pi.WaterPumpPin.WriteState(rpio.Low)
	pi.WaterPumpState = false
}

func (pi *RaspberryPi) turnOnWaterPump() {
	log.Printf("Turning on water Pump")
	pi.WaterPumpPin.WriteState(rpio.High)
	pi.WaterPumpState = true

}

func (pi RaspberryPi) getWaterTemp() {
	pi.WaterTempSensor.ReadTemperature()
}

func (pi RaspberryPi) startWaterCycle() {
	go func() {
		for {
			pi.turnOnWaterPump()
			time.Sleep(time.Second * 8)
			pi.turnOffWaterPump()
			time.Sleep(time.Hour * 6)
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
		//dht22 := NewDHT22(17)
		for {
			//temperature, tErr := dht22.Temperature()
			//humidity, hErr := dht22.Humidity()
			temperature, humidity, retried, err :=
				dht.ReadDHTxxWithRetry(dht.DHT22, 17, true, 10)
			if err != nil {
				log.Printf("Error: Error with reading dht: %v", err.Error())
			}

			log.Printf("Ambient Temperature = %v*C, Humidity = %v%% (retired: %v) \n ",
				temperature, humidity, retried)
			pi.getWaterTemp()
			time.Sleep(time.Hour)
		}

	}()
}

func (pi RaspberryPi) startAirPumpCycle() {
	go func() {
		for {
			log.Printf("Turning on air pump")
			pi.AirPump.WriteState(rpio.High)
			time.Sleep(time.Minute * 30)
			log.Printf("Turning off air pump")
			pi.AirPump.WriteState(rpio.Low)
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
