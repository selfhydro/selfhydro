package main

import (
	"errors"
	"log"
	"strconv"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type selfhydro struct {
	currentTemp       float32
	waterLevel        *WaterLevel
	waterPump         Actuator
	airPump           Actuator
	airPumpOnDuration time.Duration
	airPumpFrequency  time.Duration
	localMQTT         MQTTComms
	setup             bool
}

const (
	WATER_PUMP_PIN = 18
)

var WATER_MIN_LEVEL float32 = 80
var WATER_MAX_LEVEL float32 = 35

const (
	WATER_LEVEL_TOPIC = "/sensors/water_level"
)

func (sh *selfhydro) Setup(waterPump, airPump Actuator) error {
	sh.waterLevel = &WaterLevel{}
	sh.waterPump = waterPump
	sh.waterPump.Setup()
	sh.airPumpOnDuration = time.Minute * 40
	sh.airPump = airPump
	sh.airPump.Setup()
	sh.localMQTT = NewLocalMQTT()
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
	sh.airPump.TurnOn()
	time.AfterFunc(sh.airPumpOnDuration, func() {
		sh.airPump.TurnOff()
	})
}

func (sh selfhydro) checkWaterLevel() {
	if sh.waterLevel.GetWaterLevel() > WATER_MIN_LEVEL && !sh.waterPump.GetState() {
		log.Print("turning on water pump")
		log.Printf("water level %f", sh.waterLevel.GetWaterLevel())
		sh.waterPump.TurnOn()
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
