package sensors

import (
	"encoding/json"
	"log"

	mqttPaho "github.com/eclipse/paho.mqtt.golang"
	"github.com/selfhydro/selfhydro/mqtt"
)

type waterTemperatureMessage struct {
	sensorMessage
	Temperature float64 `json:"temperature"`
}

type WaterTemperature struct {
	Sensor
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

func (e WaterTemperature) GetLatestBatteryVoltage() float64 {
	return e.batteryVoltage
}

func (e WaterTemperature) GetSensorID() int {
	return e.id
}
