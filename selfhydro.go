package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

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
	localMQTT                   mqtt.MQTTComms
	externalMQTT                mqtt.MQTTComms
	setup                       bool
}

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
	err := sh.localMQTT.ConnectDevice()
	if err != nil {
		return err
	}
	err = sh.ambientTemperature.Subscribe(sh.localMQTT)
	if err != nil {
		return err
	}
	err = sh.ambientHumidity.Subscribe(sh.localMQTT)
	if err != nil {
		return err
	}
	err = sh.waterTemperature.Subscribe(sh.localMQTT)
	if err != nil {
		return err
	}
	sh.setupExternalMQTTComms()
	log.Println("all setup and subscribed...... going to start publishing")
	sh.runStatePublisherCycle()
	return nil
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
