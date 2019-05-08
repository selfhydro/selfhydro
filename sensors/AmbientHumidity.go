package sensors

import (
	"encoding/json"
	"log"

	"github.com/selfhydro/selfhydro/mqtt"
	mqttPaho "github.com/eclipse/paho.mqtt.golang"
)

type humidityMessage struct {
	Humidity float64 `json:"humidity"`
}

type AmbientHumidity struct {
	humidity float64
}

const AmbientHumidityTopic = "/state/ambient_humidity"

func (e *AmbientHumidity) Subscribe(mqtt mqtt.MQTTComms) error {
	if err := mqtt.SubscribeToTopic(AmbientHumidityTopic, e.HumidityHandler); err != nil {
		log.Print(err.Error())
		return err
	}
	return nil
}

func (e *AmbientHumidity) HumidityHandler(client mqttPaho.Client, message mqttPaho.Message) {
	humidity := &humidityMessage{}
	json.Unmarshal(message.Payload(), humidity)
	e.humidity = humidity.Humidity
}

func (e AmbientHumidity) GetLatestData() float64 {
	return e.humidity
}
