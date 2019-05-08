package sensors

import (
	"encoding/json"
	"log"

	"github.com/selfhydro/selfhydro/mqtt"
	mqttPaho "github.com/eclipse/paho.mqtt.golang"
)

type waterTemperatureMessage struct {
	Temperature float64 `json:"temperature"`
}

type WaterTemperature struct {
	temperature float64
}

const WaterTemperatureTopic = "/state/water_temperature"

func (e *WaterTemperature) Subscribe(mqtt mqtt.MQTTComms) error {
	if err := mqtt.SubscribeToTopic(WaterTemperatureTopic, e.TemperatureHandler); err != nil {
		log.Print(err.Error())
		return err
	}
	return nil
}

func (e *WaterTemperature) TemperatureHandler(client mqttPaho.Client, message mqttPaho.Message) {
	eM := &waterTemperatureMessage{}
	json.Unmarshal(message.Payload(), eM)
	e.temperature = eM.Temperature
}

func (e WaterTemperature) GetLatestData() float64 {
	return e.temperature
}
