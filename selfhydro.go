package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type WaterLevelMessage struct {
	WaterLevel float32 `json:"waterLevel"`
	Time       string  `json:"time"`
}

type selfhydro struct {
	currentTemp           float32
	waterLevel            WaterLevelMeasurer
	waterPump             Actuator
	waterPumpLastOnTime   time.Time
	lowWaterLevelReadings int
	airPump               Actuator
	airPumpOnDuration     time.Duration
	airPumpFrequency      time.Duration
	localMQTT             MQTTComms
	externalMQTT          MQTTComms
	setup                 bool
}

const (
	WATER_PUMP_PIN = 18
)

var WaterMinLevel float32 = 95
var WATER_MAX_LEVEL float32 = 65
var AIR_PUMP_ON_DURATION = time.Minute * 30
var AIR_PUMP_FREQUENCY = time.Minute * 60
var MinWaterPumpOffPeriod = time.Hour * 24
var MIN_LOW_WATER_READINGS = 3
var waitTimeTillReconnectAgain = time.Second * 5

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
	sh.externalMQTT = &mqttComms{}
	sh.airPumpFrequency = AIR_PUMP_FREQUENCY
	sh.airPumpOnDuration = AIR_PUMP_ON_DURATION
	sh.setup = true
	return nil
}

func (sh *selfhydro) setupExternalMQTTComms() {
	var triedToConect int
	for {
		if err := sh.externalMQTT.ConnectDevice(); err != nil {
			triedToConect++
			if triedToConect < 5 {
				fmt.Print(triedToConect)
				time.Sleep(waitTimeTillReconnectAgain)
			} else {
				log.Print("cant connect to external mqtt")
				break
			}
		} else {
			break
		}
	}
}

func (sh *selfhydro) Start() error {
	if !sh.setup {
		return errors.New("must setup selfhydro before starting (use Setup())")
	}
	sh.localMQTT.ConnectDevice()
	sh.setupExternalMQTTComms()
	sh.SubscribeToWaterLevel()
	sh.RunWaterPump()
	sh.runAirPump()
	return nil
}

func (sh *selfhydro) StopSystem() {

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
	var currentWaterLevel = sh.waterLevel.GetWaterLevelFeed()
	if currentWaterLevel > WaterMinLevel && !lastOnTimeTooRecently {
		sh.lowWaterLevelReadings++
		log.Printf("received low water reading: %f", currentWaterLevel)
	} else {
		sh.lowWaterLevelReadings = 0
	}
	if sh.lowWaterLevelReadings > MIN_LOW_WATER_READINGS {
		turnOn = true
	}
	if turnOn {
		log.Print("turning on water pump")
		log.Printf("water level %f", currentWaterLevel)
		sh.waterPump.TurnOn()
		sh.waterPumpLastOnTime = time.Now()
		sh.lowWaterLevelReadings = 0
	} else if currentWaterLevel < WATER_MAX_LEVEL && sh.waterPump.GetState() {
		log.Print("turning off water pump")
		log.Printf("water level %f", currentWaterLevel)
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
	sh.waterLevel.SetWaterLevel(float32(waterLevelFloat))
}

func (sh *selfhydro) publishWaterLevel() {
	message, err := sh.createWaterLevelMessage()
	if err != nil {
		log.Printf("Error creating sensor message: %s", err)
	}
	fmt.Print(message)
	sh.externalMQTT.publishMessage("/devices/"+sh.externalMQTT.GetDeviceID()+"/events", message)
}

func (sh *selfhydro) createWaterLevelMessage() (string, error) {
	waterLevel, time := sh.waterLevel.GetWaterLevel()
	m := WaterLevelMessage{waterLevel, time.Format("20060102150405")}
	jsonMsg, err := json.Marshal(m)
	return string(jsonMsg), err
}
