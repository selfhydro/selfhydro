package sensors

import (
	"encoding/json"
	"log"

	"github.com/selfhydro/selfhydro/mqtt"
	mqttPaho "github.com/eclipse/paho.mqtt.golang"
)

type MQTTTopic interface {
	Subscribe(mqtt mqtt.MQTTComms) error
	GetLatestData() float64
}

type temperatureMessage struct {
	Temperature float64 `json:"temperature"`
}

type AmbientTemperature struct {
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
