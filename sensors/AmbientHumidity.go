package sensors

import (
	"encoding/json"
	"log"

	mqttPaho "github.com/eclipse/paho.mqtt.golang"
	"github.com/selfhydro/selfhydro/mqtt"
)

type humidityMessage struct {
	sensorMessage
	Humidity float64 `json:"humidity"`
}

type AmbientHumidity struct {
	Sensor
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

func (e AmbientHumidity) GetLatestBatteryVoltage() float64 {
	return e.batteryVoltage
}

func (e AmbientHumidity) GetSensorID() int {
	return e.id
}
