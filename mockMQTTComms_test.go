package main

import MQTT "github.com/eclipse/paho.mqtt.golang"

type mockMQTTComms struct {
}

func (m *mockMQTTComms) ConnectDevice() error {
	return nil
}

func (m *mockMQTTComms) publishMessage(topic string, message string) {

}

func (m *mockMQTTComms) GetDeviceID() string {
	return ""
}

func (m *mockMQTTComms) SubscribeToTopic(topix string, callback MQTT.MessageHandler) {

}
func (m *mockMQTTComms) UnsubscribeFromTopic(topic string) {

}
