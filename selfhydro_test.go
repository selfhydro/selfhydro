package main

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/bchalk101/selfhydro/mocks"
	"github.com/stretchr/testify/mock"
	"gotest.tools/assert"
)

func Test_ShouldSetupSelfhydro(t *testing.T) {
	mockMQTT := &mocks.MockMQTTComms{}
	mockWaterPump := &MockActuator{}
	mockAirPump := &MockActuator{}
	mockGrowLight := &MockActuator{}
	externalMockMQTT := &mocks.MockMQTTComms{}

	mockWaterPump.On("Setup").Return(nil)
	mockAirPump.On("Setup").Return(nil)
	mockGrowLight.On("Setup").Return(nil)
	sh := selfhydro{
		localMQTT:    mockMQTT,
		externalMQTT: externalMockMQTT,
	}
	sh.Setup(mockWaterPump, mockAirPump, mockGrowLight)
	assert.Equal(t, sh.airPumpOnDuration, time.Minute*30)
}

func Test_ShouldReturnErrorIfTryingToStartButNotSetup(t *testing.T) {
	mockMQTT := &mocks.MockMQTTComms{}
	mockMQTT.On("ConnectDevice").Return(nil)
	sh := selfhydro{
		localMQTT: mockMQTT,
	}
	err := sh.Start()
	assert.Error(t, err, "must setup selfhydro before starting (use Setup())")
}

func Test_ShouldStartSelfhydro(t *testing.T) {
	mockMQTT := &mocks.MockMQTTComms{}
	mockWaterPump := &MockActuator{}
	mockAirPump := &MockActuator{}
	mockGrowLight := &MockActuator{}
	mockExternalMQTT := &mocks.MockMQTTComms{}
	mockAmbientTemperature := &mocks.MQTTTopic{}
	mockMQTT.On("ConnectDevice").Return(nil)
	mockMQTT.On("SubscribeToTopic", string("/sensors/water_level"), mock.AnythingOfType("mqtt.MessageHandler")).Return(nil)
	mockAmbientTemperature.On("Subscribe", mock.Anything, mock.Anything).Return(nil)
	mockWaterPump.On("TurnOff").Return(nil)
	mockWaterPump.On("GetState").Return(false)
	mockAirPump.On("TurnOn").Return(nil)
	mockAirPump.On("TurnOff").Return(nil)
	mockGrowLight.On("GetState").Return(true)
	mockGrowLight.On("TurnOff").Return(nil)
	mockGrowLight.On("TurnOn").Return(nil)
	mockExternalMQTT.On("GetDeviceID").Return("sdss112")
	mockExternalMQTT.On("ConnectDevice").Return(nil)
	mockExternalMQTT.On("PublishMessage", mock.Anything, mock.Anything).Return(nil)
	mockAmbientTemperature.On("GetLatestData").Return(21.00)

	sh := selfhydro{
		localMQTT:          mockMQTT,
		setup:              true,
		waterPump:          mockWaterPump,
		waterLevel:         &WaterLevel{},
		airPump:            mockAirPump,
		externalMQTT:       mockExternalMQTT,
		growLight:          mockGrowLight,
		ambientTemperature: mockAmbientTemperature,
	}
	err := sh.Start()
	time.Sleep(time.Millisecond)
	assert.Equal(t, err, nil)
	mockMQTT.AssertNumberOfCalls(t, "ConnectDevice", 1)
	mockMQTT.AssertNumberOfCalls(t, "SubscribeToTopic", 1)
	mockAirPump.AssertCalled(t, "TurnOn")
}

func Test_ShouldGetWaterLevelFromSensor(t *testing.T) {
	mockMQTT := &mocks.MockMQTTComms{}
	sh := selfhydro{
		localMQTT: mockMQTT,
	}
	mockMQTT.On("SubscribeToTopic", string("/sensors/water_level"), mock.AnythingOfType("mqtt.MessageHandler")).Return(nil)
	err := sh.SubscribeToWaterLevel()
	assert.Equal(t, err, nil)
}

func Test_ShouldLogErrorWhenCantSubscribeToWaterLevel(t *testing.T) {
	mockMQTT := &mocks.MockMQTTComms{}
	sh := selfhydro{
		localMQTT: mockMQTT,
	}
	mockMQTT.On("SubscribeToTopic", string("/sensors/water_level"), mock.AnythingOfType("mqtt.MessageHandler")).Return(errors.New("cant subscribe"))
	err := sh.SubscribeToWaterLevel()
	assert.Equal(t, err.Error(), "cant subscribe")
}

func Test_ShouldUpdateWaterLevelWhenReceivedFromTopic(t *testing.T) {
	mockMQTT := &mocks.MockMQTTComms{}
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

func Test_ShouldTurnOnWaterPumpIfWaterIsVeryLow(t *testing.T) {
	mockMQTT := &mocks.MockMQTTComms{}
	mockWaterPump := &MockActuator{}
	mockWaterLevel := &mocks.MockWaterLevelMeasurer{}
	mockWaterPump.On("GetState").Return(false)
	sh := selfhydro{
		localMQTT:             mockMQTT,
		waterPump:             mockWaterPump,
		waterLevel:            mockWaterLevel,
		lowWaterLevelReadings: 4,
	}
	mockWaterPump.On("TurnOn").Return(nil)
	mockWaterLevel.On("GetWaterLevelFeed").Return(float32(100))
	mockWaterLevel.On("GetWaterLevel").Return(float32(100))
	sh.checkWaterLevel()
	mockWaterPump.AssertNumberOfCalls(t, "TurnOn", 1)
	assert.Equal(t, sh.lowWaterLevelReadings, 0)
}

func Test_ShouldTurnOffWaterPumpIfWaterGetsToGoodLevel(t *testing.T) {
	mockMQTT := &mocks.MockMQTTComms{}
	mockWaterPump := &MockActuator{}
	mockWaterLevel := &mocks.MockWaterLevelMeasurer{}
	sh := &selfhydro{
		localMQTT:  mockMQTT,
		waterPump:  mockWaterPump,
		waterLevel: mockWaterLevel,
	}
	mockWaterPump.On("GetState").Return(true)
	mockWaterPump.On("TurnOn").Return(nil)
	mockWaterPump.On("TurnOff").Return(nil)
	mockWaterLevel.On("GetWaterLevelFeed").Return(float32(25))
	mockWaterLevel.On("GetWaterLevel").Return(float32(25))
	sh.checkWaterLevel()
	mockWaterPump.AssertNumberOfCalls(t, "TurnOff", 1)
}

func Test_ShouldTurnOnAirPumps(t *testing.T) {
	mockAirPump := &MockActuator{}
	mockAirPump.On("TurnOn").Return(nil)
	mockAirPump.On("TurnOff").Return(nil)
	sh := &selfhydro{
		airPump: mockAirPump,
	}
	sh.runAirPumpCycle()
	time.Sleep(time.Microsecond)
	mockAirPump.AssertCalled(t, "TurnOn")
}

func Test_ShouldTurnOffAirPumpAfterSetDuration(t *testing.T) {
	mockAirPump := &MockActuator{}
	mockAirPump.On("TurnOn").Return(nil)
	mockAirPump.On("TurnOff").Return(nil)
	sh := &selfhydro{
		airPump:           mockAirPump,
		airPumpOnDuration: -time.Microsecond,
	}
	sh.runAirPumpCycle()
	time.Sleep(time.Millisecond)
	mockAirPump.AssertCalled(t, "TurnOn")
	mockAirPump.AssertCalled(t, "TurnOff")
}

func Test_ShouldNotTurnOnWaterPumpAgainWhenLastTurnOnTimeWasTooRecently(t *testing.T) {
	mockMQTT := &mocks.MockMQTTComms{}
	mockWaterPump := &MockActuator{}
	mockWaterLevel := &mocks.MockWaterLevelMeasurer{}
	mockWaterPump.On("GetState").Return(false)
	sh := selfhydro{
		localMQTT:           mockMQTT,
		waterPump:           mockWaterPump,
		waterLevel:          mockWaterLevel,
		waterPumpLastOnTime: time.Now().Add(time.Hour * -2),
	}
	mockWaterPump.On("TurnOn").Return(nil)
	mockWaterLevel.On("GetWaterLevelFeed").Return(float32(81))
	mockWaterLevel.On("GetWaterLevel").Return(float32(81))
	sh.checkWaterLevel()
	mockWaterPump.AssertNumberOfCalls(t, "TurnOn", 0)
}
func Test_OnlyTurnOnWaterPumpAfterEnoughTimeHasElasped(t *testing.T) {
	mockMQTT := &mocks.MockMQTTComms{}
	mockWaterPump := &MockActuator{}
	mockWaterLevel := &mocks.MockWaterLevelMeasurer{}
	mockWaterPump.On("GetState").Return(false)
	sh := selfhydro{
		localMQTT:             mockMQTT,
		waterPump:             mockWaterPump,
		waterLevel:            mockWaterLevel,
		lowWaterLevelReadings: 3,
		waterPumpLastOnTime:   time.Now().Add(time.Hour * -24),
	}
	mockWaterPump.On("TurnOn").Return(nil)
	mockWaterLevel.On("GetWaterLevelFeed").Return(float32(100))
	mockWaterLevel.On("GetWaterLevel").Return(float32(100))
	sh.checkWaterLevel()
	mockWaterPump.AssertNumberOfCalls(t, "TurnOn", 1)
}

func Test_ShouldWaitFor3ReadingsOverMinWateLevelToTurnOnWaterPump(t *testing.T) {
	mockMQTT := &mocks.MockMQTTComms{}
	mockWaterPump := &MockActuator{}
	mockWaterLevel := &mocks.MockWaterLevelMeasurer{}
	mockWaterPump.On("GetState").Return(false)
	sh := selfhydro{
		localMQTT:             mockMQTT,
		waterPump:             mockWaterPump,
		waterLevel:            mockWaterLevel,
		waterPumpLastOnTime:   time.Now().Add(time.Hour * -6),
		lowWaterLevelReadings: 1,
	}
	mockWaterPump.On("TurnOn").Return(nil)
	mockWaterLevel.On("GetWaterLevelFeed").Return(float32(81))
	mockWaterLevel.On("GetWaterLevel").Return(float32(81))
	sh.checkWaterLevel()
	mockWaterPump.AssertNumberOfCalls(t, "TurnOn", 0)
}

func Test_ShouldHandleErrorIfCantConnectToExternalMQTT(t *testing.T) {
	mockMQTT := &mocks.MockMQTTComms{}
	sh := selfhydro{
		externalMQTT: mockMQTT,
	}
	waitTimeTillReconnectAgain = time.Microsecond
	mockMQTT.On("ConnectDevice").Return(errors.New("cant connect"))
	sh.setupExternalMQTTComms()
	mockMQTT.AssertNumberOfCalls(t, "ConnectDevice", 5)
}

func Test_ShouldPublishState(t *testing.T) {
	mockMQTT := &mocks.MockMQTTComms{}
	mockWaterLevel := &mocks.MockWaterLevelMeasurer{}
	mockAmbientTemperature := &mocks.MQTTTopic{}
	mockMQTT.On("PublishMessage", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(nil)
	mockMQTT.On("GetDeviceID").Return("device")
	time := time.Now()
	mockWaterLevel.On("GetWaterLevel").Return(float32(20), time)
	mockAmbientTemperature.On("GetLatestData").Return(21.12)
	sh := selfhydro{
		externalMQTT:       mockMQTT,
		waterLevel:         mockWaterLevel,
		ambientTemperature: mockAmbientTemperature,
	}
	sh.publishState()
	expectedMessage := fmt.Sprintf(`{"temperature":21.12,"time":"%s"}`, time.Format("20060102150405"))
	mockMQTT.AssertCalled(t, "PublishMessage", "/devices/device/events", expectedMessage)
}

func Test_ShouldTurnOnGrowLight(t *testing.T) {
	mockGrowLight := &MockActuator{}
	startTimeString := time.Now().Add(-time.Minute).Format("15:04:05")
	startTime, _ := time.Parse("15:04:05", startTimeString)

	offTimeString := time.Now().Add(time.Minute).Format("15:04:05")
	offTime, _ := time.Parse("15:04:05", offTimeString)

	sh := selfhydro{
		growLight: mockGrowLight,
	}
	mockGrowLight.On("TurnOn").Return(nil)
	mockGrowLight.On("GetState").Return(false)
	sh.changeGrowLightState(startTime, offTime)
	mockGrowLight.AssertNumberOfCalls(t, "TurnOn", 1)
}

func Test_ShouldTurnOffGrowLights(t *testing.T) {
	mockGrowLight := &MockActuator{}
	startTimeString := time.Now().Add(time.Minute).Format("15:04:05")
	startTime, _ := time.Parse("15:04:05", startTimeString)

	offTimeString := time.Now().Add(-time.Minute).Format("15:04:05")
	offTime, _ := time.Parse("15:04:05", offTimeString)

	sh := selfhydro{
		growLight: mockGrowLight,
	}
	mockGrowLight.On("TurnOn").Return(nil)
	mockGrowLight.On("TurnOff").Return(nil)
	mockGrowLight.On("GetState").Return(true)
	sh.changeGrowLightState(startTime, offTime)
	mockGrowLight.AssertNumberOfCalls(t, "TurnOn", 0)
	mockGrowLight.AssertNumberOfCalls(t, "TurnOff", 1)
}
