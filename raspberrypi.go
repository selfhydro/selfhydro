package main

import (
	"github.com/stianeikeland/go-rpio"
	"log"
	"time"
	"github.com/d2r2/go-dht"
	"fmt"
	"os"
)

type Controller interface {
	StopSystem()
	StartHydroponics()

	startSensorCycle()
	startLightCycle()
	startAirPumpCycle()
}

type RaspberryPi struct {
	GrowLedPin      RaspberryPiPin
	GrowLedState    bool
	WaterPumpState  bool
	WaterTempSensor ds18b20
	AirPumpPin      RaspberryPiPin
	MQTTClient *MQTTComms
}

func NewRaspberryPi() *RaspberryPi {
	pi := new(RaspberryPi)

	error := rpio.Open()
	if error != nil {
		log.Fatalf("Could not open rpio pins %v", error.Error())
		os.Exit(1)
	}
	//defer rpio.Close()

	pi.GrowLedPin = NewRaspberryPiPin(19)
	pi.GrowLedState = false
	pi.GrowLedPin.SetMode(rpio.Output)

	pi.WaterTempSensor.id = "28-0316838ca7ff"

	pi.AirPumpPin = NewRaspberryPiPin(21)
	pi.AirPumpPin.SetMode(rpio.Output)

	pi.MQTTClient = new(MQTTComms)
	pi.MQTTClient.authenticateDevice()

	return pi
}

func (pi *RaspberryPi) StartHydroponics() {
	pi.startSensorCycle()
	pi.startLightCycle()
	pi.startAirPumpCycle()

}

func (pi *RaspberryPi) StopSystem() {
	pi.turnOffGrowLed()
	pi.AirPumpPin.WriteState(rpio.Low)
	rpio.Close()
}

func (pi *RaspberryPi) turnOnGrowLed() {
	pi.GrowLedPin.WriteState(rpio.High)
	pi.GrowLedState = true
}

func (pi *RaspberryPi) turnOffGrowLed() {
	pi.GrowLedPin.WriteState(rpio.Low)
	pi.GrowLedState = false

}

func (pi RaspberryPi) getWaterTemp() {
	temp := pi.WaterTempSensor.ReadTemperature()
	message, _ := CreateSensorMessage(float32(temp), 0.0, 0.0, true)
	pi.MQTTClient.publishMessage(HYDRO_EVENTS_TOPIC, message)
}

func (pi RaspberryPi) startLightCycle() {
	turnOnTime, _ := time.Parse("15:04:05", "04:45:00")
	turnOffTime, _ := time.Parse("15:04:05", "23:45:00")
	go func() {
		for {
			pi.changeLEDState(turnOnTime, turnOffTime)
			time.Sleep(time.Second * 4)
		}

	}()
}
func (pi RaspberryPi) changeLEDState(turnOnTime time.Time, turnOffTime time.Time) {
	if pi.GrowLedPin.ReadState() != rpio.High && betweenTime(turnOnTime, turnOffTime) {
		log.Printf("Turning on GROW LEDS")
		pi.turnOnGrowLed()
	} else if pi.GrowLedPin.ReadState() == rpio.High && betweenTime(turnOffTime, turnOnTime.Add(time.Hour*24)) {
		log.Printf("Turning off GROW LEDS")
		pi.turnOffGrowLed()
	}
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
			sensorReading := fmt.Sprintf("Ambient Temperature = %v*C, Humidity = %v%% (retired: %v)",
				temperature, humidity, retried)

			log.Printf(sensorReading)
			pi.getWaterTemp()
			time.Sleep(time.Hour * 2)
		}

	}()
}

func (pi RaspberryPi) startAirPumpCycle() {
	go func() {
		for {
			pi.airPumpCycle(time.Minute*30, time.Hour*2)
		}
	}()
}
func (pi RaspberryPi) airPumpCycle(airPumpOnDuration time.Duration, airPumpOffDuration time.Duration) {
	log.Printf("Turning on air pump")
	pi.AirPumpPin.WriteState(rpio.High)
	time.Sleep(airPumpOnDuration)
	log.Printf("Turning off air pump")
	pi.AirPumpPin.WriteState(rpio.Low)
	time.Sleep(airPumpOffDuration)
}

func betweenTime(startTime time.Time, endTime time.Time) bool {
	currentTimeString := time.Now().Format("15:04:05")
	currentTime, _ := time.Parse("15:04:05", currentTimeString)
	if currentTime.After(startTime) && currentTime.Before(endTime) {
		return true
	}
	return false
}
