// Code generated by mockery v1.0.0. DO NOT EDIT.

package main

import mock "github.com/stretchr/testify/mock"
import mqtt "github.com/eclipse/paho.mqtt.golang"

// MockMQTTComms is an autogenerated mock type for the MQTTComms type
type MockMQTTComms struct {
	mock.Mock
}

// ConnectDevice provides a mock function with given fields:
func (_m *MockMQTTComms) ConnectDevice() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetDeviceID provides a mock function with given fields:
func (_m *MockMQTTComms) GetDeviceID() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// SubscribeToTopic provides a mock function with given fields: _a0, _a1
func (_m *MockMQTTComms) SubscribeToTopic(_a0 string, _a1 mqtt.MessageHandler) error {
	ret := _m.Called(_a0, _a1)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, mqtt.MessageHandler) error); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UnsubscribeFromTopic provides a mock function with given fields: topic
func (_m *MockMQTTComms) UnsubscribeFromTopic(topic string) {
	_m.Called(topic)
}

// publishMessage provides a mock function with given fields: topic, message
func (_m *MockMQTTComms) publishMessage(topic string, message string) {
	_m.Called(topic, message)
}