package main

import (
	"github.com/stianeikeland/go-rpio"
	"log"
	"time"
	"fmt"
	"os"
	"io/ioutil"
	"strconv"
	"strings"
)

type Controller interface {
	StopSystem()
	StartHydroponics()

	startSensorCycle()
	startLightCycle()
	startAirPumpCycle()
}

const (
	LowWaterLevel = "LOW_WATER"
)

type RaspberryPi struct {
	GrowLedPin              RaspberryPiPin
	TankOneWaterTempSensor  ds18b20
	TankTwoWaterTempSensor  ds18b20
	tankOneWaterLevelSensor Sensor
	AirPumpPin              RaspberryPiPin
	MQTTClient              MQTTComms
	alertChannel            chan string
}

func NewRaspberryPi() *RaspberryPi {
	pi := new(RaspberryPi)

	error := rpio.Open()
	if error != nil {
		log.Fatalf("Could not open rpio pins %v", error.Error())
		os.Exit(1)
	}

	pi.GrowLedPin = NewRaspberryPiPin(19)
	pi.GrowLedPin.SetMode(rpio.Output)

	pi.TankOneWaterTempSensor.id = "28-0316838ca7ff"
	pi.TankTwoWaterTempSensor.id = "28-0316838b3aff"

	pi.tankOneWaterLevelSensor = NewSensor(5)

	pi.AirPumpPin = NewRaspberryPiPin(21)
	pi.AirPumpPin.SetMode(rpio.Output)

	pi.MQTTClient = new(mqttComms)
	pi.MQTTClient.ConnectDevice()

	pi.alertChannel = make(chan string, 5)

	return pi
}

func (pi *RaspberryPi) StartHydroponics() {
	pi.startSensorCycle()
	pi.startLightCycle()
	pi.startAirPumpCycle()
	pi.monitorAlerts()
}

func (pi *RaspberryPi) monitorAlerts() {
	go func() {
		for {
			alert := <-pi.alertChannel
			switch alert {
			case LowWaterLevel:
				log.Print("Water Level is Low")
			}
		}
	}()
}

func (pi *RaspberryPi) StopSystem() {
	pi.GrowLedPin.WriteState(rpio.Low)
	pi.AirPumpPin.WriteState(rpio.Low)
	rpio.Close()
}

func (pi *RaspberryPi) publishState(tankOneTemp float64, tankTwoTemp float64, CPUTemp float64) {
	message, _ := CreateSensorMessage(tankOneTemp, tankTwoTemp, CPUTemp)
	pi.MQTTClient.publishMessage(EVENTSTOPIC, message)
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
		pi.GrowLedPin.WriteState(rpio.High)
	} else if pi.GrowLedPin.ReadState() == rpio.High && betweenTime(turnOffTime, turnOnTime.Add(time.Hour*24)) {
		log.Printf("Turning off GROW LEDS")
		pi.GrowLedPin.WriteState(rpio.Low)
	}
}
func (pi RaspberryPi) startSensorCycle() {

	go func() {
		for {
			fmt.Println("Sending sensor readings....")
			tankOneTemp := pi.TankOneWaterTempSensor.ReadTemperature()
			tankTwoTemp := pi.TankTwoWaterTempSensor.ReadTemperature()
			CPUTemp := pi.getCPUTemp()
			pi.checkWaterLevels()
			pi.publishState(tankOneTemp, tankTwoTemp, CPUTemp)
			time.Sleep(time.Hour * 4)
		}

	}()
}

func (pi RaspberryPi) checkWaterLevels() {
	tankOneState := pi.tankOneWaterLevelSensor.getState()
	if tankOneState == rpio.High {
		pi.alertChannel <- LowWaterLevel
	}
}

func (pi RaspberryPi) getCPUTemp() float64 {

	var temp float64
	data, err := ioutil.ReadFile("/sys/class/thermal/thermal_zone0/temp")
	if err != nil {
		log.Printf("Error: Can't read Raspberry Pi CPU Temp")
		return 0.0
	}
	tempData := strings.TrimSuffix(string(data), "\n")

	temp, err = strconv.ParseFloat(string(tempData), 64)
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
