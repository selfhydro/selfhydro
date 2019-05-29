package mqtt

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
	MQTT_DEFAULT_BROKER = ":1883"
	CLIENT_ID           = "selfhydro-controller"
)

func NewLocalMQTT(clientId string, brokerAddress string) *LocalMQTT {
	var broker string
	if brokerAddress == "" {
		broker = MQTT_DEFAULT_BROKER
	} else {
		broker = brokerAddress
	}
	opts := MQTT.NewClientOptions().AddBroker(broker)
	opts.SetClientID(clientId)
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

func (lmqtt *LocalMQTT) PublishMessage(topic string, message string) {
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
