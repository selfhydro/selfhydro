package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	mqtt "github.com/bchalk101/selfhydro/mqtt"
	mqttPaho "github.com/eclipse/paho.mqtt.golang"
)

type StateMessage struct {
	AmbientTemperature float64 `json:"ambientTemperature"`
	AmbientHumidity    float64 `json:"ambientHumidity"`
	WaterTemperature   float64 `json:"waterTemperature"`
	Time               string  `json:"time"`
}

type selfhydro struct {
	currentTemp           float32
	ambientTemperature    MQTTTopic
	ambientHumidity       MQTTTopic
	waterTemperature      MQTTTopic
	waterLevel            WaterLevelMeasurer
	waterPump             Actuator
	growLight             Actuator
	waterPumpLastOnTime   time.Time
	lowWaterLevelReadings int
	airPump               Actuator
	airPumpOnDuration     time.Duration
	airPumpFrequency      time.Duration
	localMQTT             mqtt.MQTTComms
	externalMQTT          mqtt.MQTTComms
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

func (sh *selfhydro) Setup(waterPump, airPump, growLight Actuator) error {
	sh.waterLevel = &WaterLevel{}
	sh.ambientTemperature = &AmbientTemperature{}
	sh.ambientHumidity = &AmbientHumidity{}
	sh.waterTemperature = &WaterTemperature{}
	sh.waterPump = waterPump
	sh.waterPump.Setup()
	sh.airPump = airPump
	sh.airPump.Setup()
	sh.growLight = growLight
	sh.growLight.Setup()
	sh.localMQTT = mqtt.NewLocalMQTT()
	sh.externalMQTT = &mqtt.GCPMQTTComms{}
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
	sh.ambientTemperature.Subscribe(sh.localMQTT)
	sh.ambientHumidity.Subscribe(sh.localMQTT)
	sh.waterTemperature.Subscribe(sh.localMQTT)
	sh.setupExternalMQTTComms()
	sh.SubscribeToWaterLevel()
	sh.runStatePublisherCycle()
	sh.RunWaterPump()
	sh.runAirPump()
	sh.runGrowLights()
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

func (sh *selfhydro) runGrowLights() {
	turnOnTime, _ := time.Parse("15:04:05", "06:00:00")
	turnOffTime, _ := time.Parse("15:04:05", "18:30:00")
	go func() {
		for {
			sh.changeGrowLightState(turnOnTime, turnOffTime)
		}
	}()
}

func (sh selfhydro) changeGrowLightState(turnOnTime time.Time, turnOffTime time.Time) {
	if !sh.growLight.GetState() && betweenTime(turnOnTime, turnOffTime) {
		log.Printf("Turning on GROW LEDS")
		sh.growLight.TurnOn()
	} else if sh.growLight.GetState() && betweenTime(turnOffTime, turnOnTime.Add(time.Hour*24)) {
		log.Printf("Turning off GROW LEDS")
		sh.growLight.TurnOff()
	}

}
func betweenTime(startTime time.Time, endTime time.Time) bool {
	currentTimeString := time.Now().Format("15:04:05")
	currentTime, _ := time.Parse("15:04:05", currentTimeString)
	if currentTime.After(startTime) && currentTime.Before(endTime) {
		return true
	}
	return false
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

func (sh *selfhydro) waterLevelHandler(client mqttPaho.Client, message mqttPaho.Message) {
	waterLevel := string(message.Payload()[:])
	waterLevelFloat, err := strconv.ParseFloat(waterLevel, 32)
	if err != nil {
		log.Print("error converting payload to float")
		return
	}
	sh.waterLevel.SetWaterLevel(float32(waterLevelFloat))
}

func (sh *selfhydro) runStatePublisherCycle() {
	go func() {
		for {
			time.Sleep(time.Minute)
			sh.publishState()
			time.Sleep(time.Minute * 15)
		}
	}()
}

func (sh *selfhydro) publishState() {
	message, err := sh.createStateMessage()
	if err != nil {
		log.Printf("error creating sensor message: %s", err)
	}
	sh.externalMQTT.PublishMessage("/devices/"+sh.externalMQTT.GetDeviceID()+"/events", message)
}

func (sh *selfhydro) createStateMessage() (string, error) {
	temperature := sh.ambientTemperature.GetLatestData()
	humidity := sh.ambientHumidity.GetLatestData()
	waterTemperature := sh.waterTemperature.GetLatestData()
	time := time.Now()
	m := StateMessage{temperature, humidity, waterTemperature, time.Format("20060102150405")}
	jsonMsg, err := json.Marshal(m)
	return string(jsonMsg), err
}
