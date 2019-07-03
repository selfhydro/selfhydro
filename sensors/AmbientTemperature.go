package sensors

import (
	"encoding/json"
	"log"

	mqttPaho "github.com/eclipse/paho.mqtt.golang"
	"github.com/selfhydro/selfhydro/mqtt"
)

type MQTTTopic interface {
	Subscribe(mqtt mqtt.MQTTComms) error
	GetLatestData() float64
	GetLatestBatteryVoltage() float64
}

type temperatureMessage struct {
	sensorMessage
	Temperature float64 `json:"temperature"`
}

type Sensor struct {
	batteryVoltage float64
	id             int
}

type sensorMessage struct {
	BatteryVoltage float32 `json:"batteryVoltage"`
	ID             int     `json:"id"`
}

type AmbientTemperature struct {
	Sensor
	temperature float64
}

const AmbientTemperatureTopic = "/state/ambient_temperature"

func (e *AmbientTemperature) Subscribe(mqtt mqtt.MQTTComms) error {
	if err := mqtt.SubscribeToTopic(AmbientTemperatureTopic, e.TemperatureHandler); err != nil {
		log.Print(err.Error())
		return err
	}
	return nil
}

func (e *AmbientTemperature) TemperatureHandler(client mqttPaho.Client, message mqttPaho.Message) {
	eM := &temperatureMessage{}
	json.Unmarshal(message.Payload(), eM)
	e.temperature = eM.Temperature
}

func (e AmbientTemperature) GetLatestData() float64 {
	return e.temperature
}

func (e AmbientTemperature) GetLatestBatteryVoltage() float64 {
	return e.batteryVoltage
}
