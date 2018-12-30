package main

import (
	"log"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

type LocalMQTT struct {
	mqttClient         MQTT.Client
	mqttOptions        *MQTT.ClientOptions
	subscribtionTopics []string
}

const (
	MQTT_LOCAL_BROKER = ":1883"
	CLIENT_ID         = "selfhydro-controller"
)

func NewLocalMQTT() *LocalMQTT {
	opts := MQTT.NewClientOptions().AddBroker(MQTT_LOCAL_BROKER)
	opts.SetClientID(CLIENT_ID)
	return &LocalMQTT{
		mqttOptions: opts,
		mqttClient:  MQTT.NewClient(opts),
	}
}

func (lmqtt *LocalMQTT) ConnectDevice() error {
	if token := lmqtt.mqttClient.Connect(); token.Wait() && token.Error() != nil {
		log.Print(token.Error())
		return token.Error()
	}
	return nil
}

func (lmqtt *LocalMQTT) publishMessage(topic string, message string) {
}

func (lmqtt *LocalMQTT) GetDeviceID() string {
	return ""
}

func (lmqtt *LocalMQTT) SubscribeToTopic(topic string, callback MQTT.MessageHandler) error {
	if token := lmqtt.mqttClient.Subscribe(topic, 1, callback); token.Wait() && token.Error() != nil {
		log.Println("error subscribing to topic ", topic)
		log.Println(token.Error())
		return token.Error()
	}
	lmqtt.subscribtionTopics = append(lmqtt.subscribtionTopics, topic)
	return nil
}

func (lmqtt *LocalMQTT) UnsubscribeFromTopic(topic string) {

}
