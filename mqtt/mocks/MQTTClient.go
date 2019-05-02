package mocks

import (
	MQTT "github.com/eclipse/paho.mqtt.golang"
)

type MockMQTTClient struct {
	PublishCalled      bool
	Connected          bool
	HasErrorConnecting bool
}

func (c MockMQTTClient) AddRoute(topic string, callback MQTT.MessageHandler) {

}

func (c MockMQTTClient) IsConnected() bool {
	return c.Connected
}

func (c MockMQTTClient) IsConnectionOpen() bool {
	return true
}

func (c *MockMQTTClient) Connect() MQTT.Token {
	token := MockMQTTToken{
		hasConnectionError: c.HasErrorConnecting,
	}
	if c.HasErrorConnecting {
		c.Connected = false
	} else {
		c.Connected = true
	}
	return &token
}

func (c MockMQTTClient) Disconnect(quiesce uint) {

}

func (c MockMQTTClient) Publish(topic string, qos byte, retained bool, payload interface{}) MQTT.Token {
	c.PublishCalled = true
	token := new(MQTT.Token)
	return *token
}

func (c MockMQTTClient) Subscribe(topic string, qos byte, callback MQTT.MessageHandler) MQTT.Token {
	token := MockMQTTToken{
		hasConnectionError: c.HasErrorConnecting,
	}
	if c.HasErrorConnecting {
		c.Connected = false
	} else {
		c.Connected = true
	}
	return &token
}

func (c MockMQTTClient) SubscribeMultiple(filters map[string]byte, callback MQTT.MessageHandler) MQTT.Token {
	token := new(MQTT.Token)

	return *token
}

func (c MockMQTTClient) Unsubscribe(topics ...string) MQTT.Token {
	token := new(MQTT.Token)

	return *token
}

func (c MockMQTTClient) OptionsReader() MQTT.ClientOptionsReader {
	r := MQTT.ClientOptionsReader{}
	return r
}
