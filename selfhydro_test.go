package main

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/selfhydro/selfhydro/mocks"
	mqttMocks "github.com/selfhydro/selfhydro/mqtt/mocks"
	"github.com/stretchr/testify/mock"
	"gotest.tools/assert"
)

func Test_ShouldSetupSelfhydro(t *testing.T) {
	mockMQTT := &mqttMocks.MockMQTTComms{}
	externalMockMQTT := &mqttMocks.MockMQTTComms{}

	sh := selfhydro{
		localMQTT:    mockMQTT,
		externalMQTT: externalMockMQTT,
	}
	sh.Setup()
}

func Test_ShouldReturnErrorIfTryingToStartButNotSetup(t *testing.T) {
	mockMQTT := &mqttMocks.MockMQTTComms{}
	mockMQTT.On("ConnectDevice").Return(nil)
	sh := selfhydro{
		localMQTT: mockMQTT,
	}
	err := sh.Start()
	assert.Error(t, err, "must setup selfhydro before starting (use Setup())")
}

func Test_ShouldStartSelfhydro(t *testing.T) {
	mockMQTT := &mqttMocks.MockMQTTComms{}
	mockExternalMQTT := &mqttMocks.MockMQTTComms{}
	mockAmbientTemperature := &mocks.MQTTTopic{}
	mockAmbientHumidity := &mocks.MQTTTopic{}
	mockWaterTemperature := &mocks.MQTTTopic{}
	mockWaterElectricalConductivity := &mocks.MQTTTopic{}
	mockMQTT.On("ConnectDevice").Return(nil)
	mockMQTT.On("SubscribeToTopic", string("/sensors/water_level"), mock.AnythingOfType("mqtt.MessageHandler")).Return(nil)
	mockAmbientTemperature.On("Subscribe", mock.Anything, mock.Anything).Return(nil)
	mockAmbientHumidity.On("Subscribe", mock.Anything, mock.Anything).Return(nil)
	mockWaterTemperature.On("Subscribe", mock.Anything, mock.Anything).Return(nil)
	mockWaterElectricalConductivity.On("Subscribe", mock.Anything, mock.Anything).Return(nil)
	mockExternalMQTT.On("GetDeviceID").Return("sdss112")
	mockExternalMQTT.On("ConnectDevice").Return(nil)
	mockExternalMQTT.On("PublishMessage", mock.Anything, mock.Anything).Return(nil)
	mockAmbientTemperature.On("GetLatestData").Return(21.00)
	mockAmbientHumidity.On("GetLatestData").Return(21.00)
	mockWaterTemperature.On("GetLatestData").Return(21.00)

	sh := selfhydro{
		localMQTT:                   mockMQTT,
		setup:                       true,
		waterLevel:                  &WaterLevel{},
		externalMQTT:                mockExternalMQTT,
		ambientTemperature:          mockAmbientTemperature,
		ambientHumidity:             mockAmbientHumidity,
		waterTemperature:            mockWaterTemperature,
		waterElectricalConductivity: mockWaterElectricalConductivity,
	}
	err := sh.Start()
	time.Sleep(time.Millisecond)
	assert.Equal(t, err, nil)
	mockMQTT.AssertNumberOfCalls(t, "ConnectDevice", 1)
	mockMQTT.AssertNumberOfCalls(t, "SubscribeToTopic", 1)
}

func Test_ShouldGetWaterLevelFromSensor(t *testing.T) {
	mockMQTT := &mqttMocks.MockMQTTComms{}
	sh := selfhydro{
		localMQTT: mockMQTT,
	}
	mockMQTT.On("SubscribeToTopic", string("/sensors/water_level"), mock.AnythingOfType("mqtt.MessageHandler")).Return(nil)
	err := sh.SubscribeToWaterLevel()
	assert.Equal(t, err, nil)
}

func Test_ShouldLogErrorWhenCantSubscribeToWaterLevel(t *testing.T) {
	mockMQTT := &mqttMocks.MockMQTTComms{}
	sh := selfhydro{
		localMQTT: mockMQTT,
	}
	mockMQTT.On("SubscribeToTopic", string("/sensors/water_level"), mock.AnythingOfType("mqtt.MessageHandler")).Return(errors.New("cant subscribe"))
	err := sh.SubscribeToWaterLevel()
	assert.Equal(t, err.Error(), "cant subscribe")
}

func Test_ShouldUpdateWaterLevelWhenReceivedFromTopic(t *testing.T) {
	mockMQTT := &mqttMocks.MockMQTTComms{}
	mockMessage := &mocks.MockMQTTMessage{
		ReceivedPayload: []byte("22.4"),
	}
	sh := &selfhydro{
		localMQTT:  mockMQTT,
		waterLevel: &WaterLevel{},
	}
	mockMQTT.On("SubscribeToTopic", string("/sensors/water_level"), mock.AnythingOfType("mqtt.MessageHandler")).Return(nil).Run(func(args mock.Arguments) {
		sh.waterLevelHandler(nil, mockMessage)
	})
	err := sh.SubscribeToWaterLevel()
	waterLevel, _ := sh.waterLevel.GetWaterLevel()
	assert.Equal(t, waterLevel, float32(22.4))
	assert.Equal(t, err, nil)
}

func Test_ShouldHandleErrorIfCantConnectToExternalMQTT(t *testing.T) {
	mockMQTT := &mqttMocks.MockMQTTComms{}
	sh := selfhydro{
		externalMQTT: mockMQTT,
	}
	waitTimeTillReconnectAgain = time.Microsecond
	mockMQTT.On("ConnectDevice").Return(errors.New("cant connect"))
	sh.setupExternalMQTTComms()
	mockMQTT.AssertNumberOfCalls(t, "ConnectDevice", 5)
}

func Test_ShouldPublishState(t *testing.T) {
	mockMQTT := &mqttMocks.MockMQTTComms{}
	mockWaterLevel := &mocks.MockWaterLevelMeasurer{}
	mockAmbientTemperature := &mocks.MQTTTopic{}
	mockAmbientHumidity := &mocks.MQTTTopic{}
	mockWaterTemperature := &mocks.MQTTTopic{}
	mockWaterElectricalConductivity := &mocks.MQTTTopic{}
	mockMQTT.On("PublishMessage", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(nil)
	mockMQTT.On("GetDeviceID").Return("device")
	time := time.Now()
	mockWaterLevel.On("GetWaterLevel").Return(float32(20), time)
	mockAmbientTemperature.On("GetLatestData").Return(21.12)
	mockAmbientHumidity.On("GetLatestData").Return(43.22)
	mockWaterTemperature.On("GetLatestData").Return(13.22)
	mockWaterElectricalConductivity.On("GetLatestData").Return(1.22)

	sh := selfhydro{
		externalMQTT:                mockMQTT,
		waterLevel:                  mockWaterLevel,
		ambientTemperature:          mockAmbientTemperature,
		ambientHumidity:             mockAmbientHumidity,
		waterTemperature:            mockWaterTemperature,
		waterElectricalConductivity: mockWaterElectricalConductivity,
	}
	sh.publishState()
	expectedMessage := fmt.Sprintf(`{"ambientTemperature":21.12,"ambientHumidity":43.22,"waterTemperature":13.22,"waterElectricalConductivity":1.22,"time":"%s"}`, time.Format("20060102150405"))
	mockMQTT.AssertCalled(t, "PublishMessage", "/devices/device/events", expectedMessage)
}

