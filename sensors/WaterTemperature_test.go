package sensors

import (
	"testing"

	"github.com/selfhydro/selfhydro/mocks"
	mqttMocks "github.com/selfhydro/selfhydro/mqtt/mocks"
	mqtt "github.com/eclipse/paho.mqtt.golang"

	"github.com/stretchr/testify/mock"
	"gotest.tools/assert"
)

func Test_ShouldSubscribeToWaterTemperatureTopic(t *testing.T) {
	mockMQTT := &mqttMocks.MockMQTTComms{}
	mockMQTTClient := &mqttMocks.MockMQTTClient{}
	mockMQTTMessage := &mocks.MockMQTTMessage{
		ReceivedPayload: []byte(`{"temperature":10.76101}`),
	}
	mockMQTT.On("SubscribeToTopic", string("/state/water_temperature"), mock.Anything).Run(func(args mock.Arguments) {
		args[1].(mqtt.MessageHandler)(mockMQTTClient, mockMQTTMessage)
	}).Return(nil)
	e := &WaterTemperature{}
	e.Subscribe(mockMQTT)
	assert.Equal(t, e.temperature, 10.76101)
}
