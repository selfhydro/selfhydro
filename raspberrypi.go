package main

import (
	"github.com/stianeikeland/go-rpio"
	"log"
	"time"
	"github.com/d2r2/go-dht"
	"fmt"
	"os"
	"io/ioutil"
	"strconv"
)

type Controller interface {
	StopSystem()
	StartHydroponics()

	startSensorCycle()
	startLightCycle()
	startAirPumpCycle()
}

type RaspberryPi struct {
	GrowLedPin             RaspberryPiPin
	GrowLedState           bool
	WaterPumpState         bool
	TankOneWaterTempSensor ds18b20
	TankTwoWaterTempSensor ds18b20
	AirPumpPin             RaspberryPiPin
	MQTTClient             *MQTTComms
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

	pi.TankOneWaterTempSensor.id = "28-0316838ca7ff"
	pi.TankTwoWaterTempSensor.id = ""

	pi.AirPumpPin = NewRaspberryPiPin(21)
	pi.AirPumpPin.SetMode(rpio.Output)

	pi.MQTTClient = new(MQTTComms)
	pi.MQTTClient.ConnectDevice()

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

func (pi *RaspberryPi) publishState(tankOneTemp float64, tankTwoTemp float64, CPUTemp float64) {
	message, _ := CreateSensorMessage(tankOneTemp, tankTwoTemp, CPUTemp)
	pi.MQTTClient.publishMessage(EVENTSTOPIC, message)
}

func (pi *RaspberryPi) turnOnGrowLed() {
	pi.GrowLedPin.WriteState(rpio.High)
	pi.GrowLedState = true
}

func (pi *RaspberryPi) turnOffGrowLed() {
	pi.GrowLedPin.WriteState(rpio.Low)
	pi.GrowLedState = false

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
		for {
			temperature, humidity, retried, err :=
				dht.ReadDHTxxWithRetry(dht.DHT22, 17, true, 10)
			if err != nil {
				log.Printf("Error: Error with reading dht: %v", err.Error())
			}
			sensorReading := fmt.Sprintf("Ambient Temperature = %v*C, Humidity = %v%% (retired: %v)",
				temperature, humidity, retried)

			log.Printf(sensorReading)
			tankOneTemp := pi.TankOneWaterTempSensor.ReadTemperature()
			//tankTwoTemp := pi.TankTwoWaterTempSensor.ReadTemperature()
			CPUTemp := pi.getCPUTemp()
			fmt.Println("Sending sensor readings....")
			pi.publishState(tankOneTemp, 0.0, CPUTemp)
			time.Sleep(time.Hour * 4)
		}

	}()
}

func (pi RaspberryPi) getCPUTemp() float64 {

	var temp float64
	data, err := ioutil.ReadFile("/sys/class/thermal/thermal_zone0/temp")
	if err != nil {
		log.Printf("Error: Can't read Raspberry Pi CPU Temp")
		return 0.0
	}
	temp, err = strconv.ParseFloat(string(data), 64)
	if err != nil {
		panic(err)
	}
	log.Printf("CPU Temp: %v", temp/1000)
	return temp / 1000

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
