package main

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"gotest.tools/assert"
)

func TestShouldGetAmbientTemp(t *testing.T) {
	sh := selfhydro{}
	ambientTemp := sh.GetAmbientTemp()
	assert.Equal(t, float32(10), ambientTemp)
}

func Test_ShouldGetWaterLevelFromSensor(t *testing.T) {
	mockMQTT := &MockMQTTComms{}
	sh := selfhydro{
		localMQTT: mockMQTT,
	}
	mockMQTT.On("SubscribeToTopic", string("/sensors/water_level"), mock.AnythingOfType("mqtt.MessageHandler")).Return(nil)
	err := sh.SubscribeToWaterLevel()
	assert.Equal(t, err, nil)
}

func Test_ShouldLogErrorWhenCantSubscribeToWaterLevel(t *testing.T) {
	mockMQTT := &MockMQTTComms{}
	sh := selfhydro{
		localMQTT: mockMQTT,
	}
	mockMQTT.On("SubscribeToTopic", string("/sensors/water_level"), mock.AnythingOfType("mqtt.MessageHandler")).Return(errors.New("cant subscribe"))
	err := sh.SubscribeToWaterLevel()
	assert.Equal(t, err.Error(), "cant subscribe")
}

func Test_ShouldUpdateWaterLevelWhenReceivedFromTopic(t *testing.T) {
	mockMQTT := &MockMQTTComms{}
	mockMessage := &mockMessage{}
	sh := &selfhydro{
		localMQTT:  mockMQTT,
		waterLevel: &WaterLevel{},
	}
	mockMQTT.On("SubscribeToTopic", string("/sensors/water_level"), mock.AnythingOfType("mqtt.MessageHandler")).Return(nil).Run(func(args mock.Arguments) {
		waterLevelHandler(nil, mockMessage)
	})
	err := sh.SubscribeToWaterLevel()
	time.Sleep(time.Second)
	assert.Equal(t, sh.waterLevel.GetWaterLevel(), float32(2.24))
	assert.Equal(t, err, nil)
}
