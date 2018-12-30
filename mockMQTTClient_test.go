package main

import (
	MQTT "github.com/eclipse/paho.mqtt.golang"
)

type mockClient struct {
	publishCalled      bool
	connected          bool
	hasErrorConnecting bool
}

func (c mockClient) AddRoute(topic string, callback MQTT.MessageHandler) {

}

func (c mockClient) IsConnected() bool {
	return c.connected
}

func (c mockClient) IsConnectionOpen() bool {
	return true
}

func (c *mockClient) Connect() MQTT.Token {
	token := mockMQTTToken{
		hasConnectionError: c.hasErrorConnecting,
	}
	if c.hasErrorConnecting {
		c.connected = false
	} else {
		c.connected = true
	}
	return &token
}

func (c mockClient) Disconnect(quiesce uint) {

}

func (c mockClient) Publish(topic string, qos byte, retained bool, payload interface{}) MQTT.Token {
	c.publishCalled = true
	token := new(MQTT.Token)
	return *token
}

func (c mockClient) Subscribe(topic string, qos byte, callback MQTT.MessageHandler) MQTT.Token {
	token := mockMQTTToken{
		hasConnectionError: c.hasErrorConnecting,
	}
	if c.hasErrorConnecting {
		c.connected = false
	} else {
		c.connected = true
	}
	return &token
}

func (c mockClient) SubscribeMultiple(filters map[string]byte, callback MQTT.MessageHandler) MQTT.Token {
	token := new(MQTT.Token)

	return *token
}

func (c mockClient) Unsubscribe(topics ...string) MQTT.Token {
	token := new(MQTT.Token)

	return *token
}

func (c mockClient) OptionsReader() MQTT.ClientOptionsReader {
	r := MQTT.ClientOptionsReader{}
	return r
}
