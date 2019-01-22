package main

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"gotest.tools/assert"
)

func Test_ShouldSetupSelfhydro(t *testing.T) {
	mockMQTT := &MockMQTTComms{}
	mockWaterPump := &MockActuator{}
	mockAirPump := &MockActuator{}
	mockWaterPump.On("Setup").Return(nil)
	mockAirPump.On("Setup").Return(nil)

	sh := selfhydro{
		localMQTT: mockMQTT,
	}
	sh.Setup(mockWaterPump, mockAirPump)
	assert.Equal(t, sh.waterLevel.waterLevel, float32(0))
	assert.Equal(t, sh.airPumpOnDuration, time.Minute*30)
}

func Test_ShouldReturnErrorIfTryingToStartButNotSetup(t *testing.T) {
	mockMQTT := &MockMQTTComms{}
	mockMQTT.On("ConnectDevice").Return(nil)
	sh := selfhydro{
		localMQTT: mockMQTT,
	}
	err := sh.Start()
	assert.Error(t, err, "must setup selfhydro before starting (use Setup())")
}

func Test_ShouldStartSelfhydro(t *testing.T) {
	mockMQTT := &MockMQTTComms{}
	mockWaterPump := &MockActuator{}
	mockAirPump := &MockActuator{}
	mockMQTT.On("ConnectDevice").Return(nil)
	mockMQTT.On("SubscribeToTopic", string("/sensors/water_level"), mock.AnythingOfType("mqtt.MessageHandler")).Return(nil)
	mockWaterPump.On("TurnOff").Return(nil)
	mockWaterPump.On("GetState").Return(false)
	mockAirPump.On("TurnOn").Return(nil)
	mockAirPump.On("TurnOff").Return(nil)
	sh := selfhydro{
		localMQTT:  mockMQTT,
		setup:      true,
		waterPump:  mockWaterPump,
		waterLevel: &WaterLevel{},
		airPump:    mockAirPump,
	}
	err := sh.Start()
	time.Sleep(time.Millisecond)
	assert.Equal(t, err, nil)
	mockMQTT.AssertNumberOfCalls(t, "ConnectDevice", 1)
	mockMQTT.AssertNumberOfCalls(t, "SubscribeToTopic", 1)
	mockAirPump.AssertCalled(t, "TurnOn")
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
	mockMessage := &mockMessage{
		payload: []byte("22.4"),
	}
	sh := &selfhydro{
		localMQTT:  mockMQTT,
		waterLevel: &WaterLevel{},
	}
	mockMQTT.On("SubscribeToTopic", string("/sensors/water_level"), mock.AnythingOfType("mqtt.MessageHandler")).Return(nil).Run(func(args mock.Arguments) {
		sh.waterLevelHandler(nil, mockMessage)
	})
	err := sh.SubscribeToWaterLevel()
	assert.Equal(t, sh.waterLevel.GetWaterLevel(), float32(22.4))
	assert.Equal(t, err, nil)
}

func Test_ShouldTurnOnWaterPumpIfWaterIsVeryLow(t *testing.T) {
	mockMQTT := &MockMQTTComms{}
	mockWaterPump := &MockActuator{}
	waterLevel := new(WaterLevel)
	mockWaterPump.On("GetState").Return(false)
	sh := selfhydro{
		localMQTT:             mockMQTT,
		waterPump:             mockWaterPump,
		waterLevel:            waterLevel,
		lowWaterLevelReadings: 4,
	}
	mockWaterPump.On("TurnOn").Return(nil)
	sh.waterLevel.waterLevel = float32(81.0)
	sh.checkWaterLevel()
	mockWaterPump.AssertNumberOfCalls(t, "TurnOn", 1)
	assert.Equal(t, sh.lowWaterLevelReadings, 0)
}

func Test_ShouldTurnOffWaterPumpIfWaterGetsToGoodLevel(t *testing.T) {
	mockMQTT := &MockMQTTComms{}
	mockWaterPump := &MockActuator{}
	sh := &selfhydro{
		localMQTT:  mockMQTT,
		waterPump:  mockWaterPump,
		waterLevel: &WaterLevel{},
	}
	mockWaterPump.On("GetState").Return(true)
	mockWaterPump.On("TurnOn").Return(nil)
	mockWaterPump.On("TurnOff").Return(nil)
	sh.waterLevel.waterLevel = float32(25.0)

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
	mockMQTT := &MockMQTTComms{}
	mockWaterPump := &MockActuator{}
	waterLevel := new(WaterLevel)
	mockWaterPump.On("GetState").Return(false)
	sh := selfhydro{
		localMQTT:           mockMQTT,
		waterPump:           mockWaterPump,
		waterLevel:          waterLevel,
		waterPumpLastOnTime: time.Now().Add(time.Hour * -2),
	}
	mockWaterPump.On("TurnOn").Return(nil)
	sh.waterLevel.waterLevel = float32(81.0)
	sh.checkWaterLevel()
	mockWaterPump.AssertNumberOfCalls(t, "TurnOn", 0)
}
func Test_OnlyTurnOnWaterPumpAfterEnoughTimeHasElasped(t *testing.T) {
	mockMQTT := &MockMQTTComms{}
	mockWaterPump := &MockActuator{}
	waterLevel := new(WaterLevel)
	mockWaterPump.On("GetState").Return(false)
	sh := selfhydro{
		localMQTT:             mockMQTT,
		waterPump:             mockWaterPump,
		waterLevel:            waterLevel,
		lowWaterLevelReadings: 3,
		waterPumpLastOnTime:   time.Now().Add(time.Hour * -6),
	}
	mockWaterPump.On("TurnOn").Return(nil)
	sh.waterLevel.waterLevel = float32(81.0)
	sh.checkWaterLevel()
	mockWaterPump.AssertNumberOfCalls(t, "TurnOn", 1)
}

func Test_ShouldWaitFor3ReadingsOverMinWateLevelToTurnOnWaterPump(t *testing.T) {
	mockMQTT := &MockMQTTComms{}
	mockWaterPump := &MockActuator{}
	waterLevel := new(WaterLevel)
	mockWaterPump.On("GetState").Return(false)
	sh := selfhydro{
		localMQTT:             mockMQTT,
		waterPump:             mockWaterPump,
		waterLevel:            waterLevel,
		waterPumpLastOnTime:   time.Now().Add(time.Hour * -6),
		lowWaterLevelReadings: 1,
	}
	mockWaterPump.On("TurnOn").Return(nil)
	sh.waterLevel.waterLevel = float32(81.0)
	sh.checkWaterLevel()
	mockWaterPump.AssertNumberOfCalls(t, "TurnOn", 0)
}
