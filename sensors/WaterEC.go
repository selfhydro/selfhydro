package sensors

import (
	"encoding/json"
	"log"

	mqttPaho "github.com/eclipse/paho.mqtt.golang"
	"github.com/selfhydro/selfhydro/mqtt"
)

type waterECMessage struct {
	ElectricalConductivity float64 `json:"temperature"`
}

type WaterEC struct {
	electricalConducivity float64
}

const WaterECTopic = "/state/water_ec"

func (e *WaterEC) Subscribe(mqtt mqtt.MQTTComms) error {
	if err := mqtt.SubscribeToTopic(WaterECTopic, e.ECHandler); err != nil {
		log.Print(err.Error())
		return err
	}
	return nil
}

func (e *WaterEC) ECHandler(client mqttPaho.Client, message mqttPaho.Message) {
	eM := &waterECMessage{}
	json.Unmarshal(message.Payload(), eM)
	e.electricalConducivity = eM.ElectricalConductivity
}

func (e WaterEC) GetLatestData() float64 {
	return e.electricalConducivity
}
