package main

import (
	"encoding/json"
	"log"

	"github.com/bchalk101/selfhydro/mqtt"
	mqttPaho "github.com/eclipse/paho.mqtt.golang"
)

type MQTTTopic interface {
	Subscribe(mqtt mqtt.MQTTComms) error
	GetLatestData() interface{}
}

type environmentMessage struct {
	Temperature float64 `json:"temperature"`
	Humidity    float64 `json:"humidity"`
}

type Environment struct {
	temperature float64
	humidity    float64
}

const AmbientTopic = "/sensors/ambient_temp_humidity"

func (e *Environment) Subscribe(mqtt mqtt.MQTTComms) error {
	if err := mqtt.SubscribeToTopic(AmbientTopic, e.TempAndHumidityHandler); err != nil {
		log.Print(err.Error())
		return err
	}
	return nil
}

func (e *Environment) TempAndHumidityHandler(client mqttPaho.Client, message mqttPaho.Message) {
	eM := &environmentMessage{}
	json.Unmarshal(message.Payload(), eM)
	e.humidity = eM.Humidity
	e.temperature = eM.Temperature
}

func (e Environment) GetLatestData() interface{} {
	return e
}

func (e Environment) GetTemp() float64 {
	return e.temperature
}

func (e Environment) GetHumidity() float64 {
	return e.humidity
}
