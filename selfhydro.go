package main

import (
	"errors"
	"log"
	"strconv"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type selfhydro struct {
	currentTemp           float32
	waterLevel            *WaterLevel
	waterPump             Actuator
	waterPumpLastOnTime   time.Time
	lowWaterLevelReadings int
	airPump               Actuator
	airPumpOnDuration     time.Duration
	airPumpFrequency      time.Duration
	localMQTT             MQTTComms
	setup                 bool
}

const (
	WATER_PUMP_PIN = 18
)

var WaterMinLevel float32 = 80
var WATER_MAX_LEVEL float32 = 35
var AIR_PUMP_ON_DURATION = time.Minute * 30
var AIR_PUMP_FREQUENCY = time.Minute * 60
var MinWaterPumpOffPeriod = time.Hour * 5
var MIN_LOW_WATER_READINGS = 3

const (
	WATER_LEVEL_TOPIC = "/sensors/water_level"
)

func (sh *selfhydro) Setup(waterPump, airPump Actuator) error {
	sh.waterLevel = &WaterLevel{}
	sh.waterPump = waterPump
	sh.waterPump.Setup()
	sh.airPump = airPump
	sh.airPump.Setup()
	sh.localMQTT = NewLocalMQTT()
	sh.airPumpFrequency = AIR_PUMP_FREQUENCY
	sh.airPumpOnDuration = AIR_PUMP_ON_DURATION
	sh.setup = true
	return nil
}

func (sh *selfhydro) Start() error {
	if !sh.setup {
		return errors.New("must setup selfhydro before starting (use Setup())")
	}
	sh.localMQTT.ConnectDevice()
	sh.SubscribeToWaterLevel()
	sh.RunWaterPump()
	sh.runAirPump()
	return nil
}

func (sh *selfhydro) StopSelfhydro() {

}

func (sh selfhydro) SubscribeToWaterLevel() error {
	if err := sh.localMQTT.SubscribeToTopic(WATER_LEVEL_TOPIC, sh.waterLevelHandler); err != nil {
		log.Print(err.Error())
		return err
	}
	return nil
}

func (sh *selfhydro) RunWaterPump() {
	go func() {
		for {
			sh.checkWaterLevel()
		}
	}()
}

func (sh selfhydro) runAirPump() {
	go func() {
		for {
			sh.runAirPumpCycle()
			time.Sleep(sh.airPumpFrequency)
		}
	}()
}

func (sh selfhydro) runAirPumpCycle() {
	log.Print("turning on air pumps")
	sh.airPump.TurnOn()
	time.AfterFunc(sh.airPumpOnDuration, func() {
		log.Print("turning off air pumps")
		sh.airPump.TurnOff()
	})
}

func (sh *selfhydro) checkWaterLevel() {
	var turnOn = false
	var lastOnTimeTooRecently = time.Now().Sub(sh.waterPumpLastOnTime) <= MinWaterPumpOffPeriod && sh.waterPumpLastOnTime != time.Time{}
	if lastOnTimeTooRecently {
		log.Print("Water pump turned on too recently not turning on")
	}
	if sh.waterLevel.GetWaterLevel() > WaterMinLevel && !lastOnTimeTooRecently {
		sh.lowWaterLevelReadings++
	}
	if sh.lowWaterLevelReadings > MIN_LOW_WATER_READINGS {
		turnOn = true
	}
	if turnOn {
		log.Print("turning on water pump")
		log.Printf("water level %f", sh.waterLevel.GetWaterLevel())
		sh.waterPump.TurnOn()
		sh.waterPumpLastOnTime = time.Now()
		sh.lowWaterLevelReadings = 0
	} else if sh.waterLevel.GetWaterLevel() < WATER_MAX_LEVEL && sh.waterPump.GetState() {
		log.Print("turning off water pump")
		log.Printf("water level %f", sh.waterLevel.GetWaterLevel())
		sh.waterPump.TurnOff()
	}
}

func (sh *selfhydro) waterLevelHandler(client mqtt.Client, message mqtt.Message) {
	waterLevel := string(message.Payload()[:])
	waterLevelFloat, err := strconv.ParseFloat(waterLevel, 32)
	if err != nil {
		log.Print("error converting payload to float")
		return
	}
	sh.waterLevel.waterLevel = float32(waterLevelFloat)
}
