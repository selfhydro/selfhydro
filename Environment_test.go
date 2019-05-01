package main

import (
	"testing"

	"github.com/bchalk101/selfhydro/mocks"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/stretchr/testify/mock"
	"gotest.tools/assert"
)

func Test_ShouldReturnCurrentTempAndHumidity(t *testing.T) {
	e := Environment{
		temperature: 23.2,
		humidity:    45.2,
	}
	currentReadings := e.GetLatestData()
	assert.Equal(t, e.temperature, currentReadings.(Environment).temperature)
	assert.Equal(t, e.humidity, currentReadings.(Environment).humidity)
}

func Test_ShouldSubscribeToEnvironmentTopic(t *testing.T) {
	mockMQTT := &mocks.MockMQTTComms{}
	mockMQTTClient := &mocks.MockMQTTClient{}
	mockMQTTMessage := &mocks.MockMQTTMessage{
		ReceivedPayload: []byte(`{"humidity":54.7338,"temperature":20.76101}`),
	}
	mockMQTT.On("SubscribeToTopic", string("/sensors/ambient_temp_humidity"), mock.Anything).Run(func(args mock.Arguments) {
		args[1].(mqtt.MessageHandler)(mockMQTTClient, mockMQTTMessage)
	}).Return(nil)
	e := &Environment{}
	e.Subscribe(mockMQTT)
	assert.Equal(t, e.temperature, 20.76101)
	assert.Equal(t, e.humidity, 54.7338)
}
