package main

import (
	"encoding/binary"
	"log"
	"math"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type selfhydro struct {
	currentTemp float32
	waterLevel  *WaterLevel
	localMQTT   MQTTComms
}

const (
	WATER_LEVEL_TOPIC = "/sensors/water_level"
)

var waterLevelChannel chan float32

func (sh *selfhydro) Setup() error {
	sh.localMQTT.ConnectDevice()
	sh.waterLevel = &WaterLevel{}
	return nil
}

func (sh selfhydro) GetAmbientTemp() float32 {

	return 10
}

func (sh selfhydro) SubscribeToWaterLevel() error {
	go sh.updateWaterLevel()
	if err := sh.localMQTT.SubscribeToTopic(WATER_LEVEL_TOPIC, waterLevelHandler); err != nil {
		log.Print(err.Error())
		return err
	}
	return nil
}

func (sh *selfhydro) updateWaterLevel() {
	waterLevelChannel = make(chan float32, 1)
	for {
		waterLevel := <-waterLevelChannel
		sh.waterLevel.waterLevel = waterLevel
	}
}

var waterLevelHandler = func(client mqtt.Client, message mqtt.Message) {
	waterLevelChannel <- float32frombytes(message.Payload())
}

func float32frombytes(bytes []byte) float32 {
	bits := binary.LittleEndian.Uint32(bytes)
	float := math.Float32frombits(bits)
	return float
}
