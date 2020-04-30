package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	mqttPaho "github.com/eclipse/paho.mqtt.golang"
	mqtt "github.com/selfhydro/selfhydro/mqtt"
	"github.com/selfhydro/selfhydro/sensors"
)

type StateMessage struct {
	AmbientTemperature          float64 `json:"ambientTemperature"`
	AmbientHumidity             float64 `json:"ambientHumidity"`
	WaterTemperature            float64 `json:"waterTemperature"`
	WaterElectricalConductivity float64 `json:"waterElectricalConductivity"`
	Time                        string  `json:"time"`
}

type selfhydro struct {
	currentTemp                 float32
	ambientTemperature          sensors.MQTTTopic
	ambientHumidity             sensors.MQTTTopic
	waterTemperature            sensors.MQTTTopic
	waterElectricalConductivity sensors.MQTTTopic
	waterLevel                  WaterLevelMeasurer
	waterPumpLastOnTime         time.Time
	lowWaterLevelReadings       int
	airPumpOnDuration           time.Duration
	airPumpFrequency            time.Duration
	localMQTT                   mqtt.MQTTComms
	externalMQTT                mqtt.MQTTComms
	setup                       bool
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

func (sh *selfhydro) Setup() error {
	sh.waterLevel = &WaterLevel{}
	sh.ambientTemperature = &sensors.AmbientTemperature{}
	sh.ambientHumidity = &sensors.AmbientHumidity{}
	sh.waterTemperature = &sensors.WaterTemperature{}
	sh.waterElectricalConductivity = &sensors.WaterElectricalConductivity{}
	sh.localMQTT = mqtt.NewLocalMQTT(mqtt.CLIENT_ID, "")
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
	sh.waterElectricalConductivity.Subscribe(sh.localMQTT)
	sh.setupExternalMQTTComms()
	sh.SubscribeToWaterLevel()
	sh.runStatePublisherCycle()
	return nil
}

func (sh selfhydro) SubscribeToWaterLevel() error {
	if err := sh.localMQTT.SubscribeToTopic(WATER_LEVEL_TOPIC, sh.waterLevelHandler); err != nil {
		log.Print(err.Error())
		return err
	}
	return nil
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
		time.Sleep(time.Minute * 30)
		for {
			sh.publishState()
			time.Sleep(time.Hour)
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
	waterElectricalConductivity := sh.waterElectricalConductivity.GetLatestData()
	time := time.Now()
	m := StateMessage{temperature, humidity, waterTemperature, waterElectricalConductivity, time.Format("20060102150405")}
	jsonMsg, err := json.Marshal(m)
	return string(jsonMsg), err
}
