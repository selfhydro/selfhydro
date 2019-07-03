// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"
import mqtt "github.com/selfhydro/selfhydro/mqtt"

// MQTTTopic is an autogenerated mock type for the MQTTTopic type
type MQTTTopic struct {
	mock.Mock
}

// GetLatestBatteryVoltage provides a mock function with given fields:
func (_m *MQTTTopic) GetLatestBatteryVoltage() float64 {
	ret := _m.Called()

	var r0 float64
	if rf, ok := ret.Get(0).(func() float64); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(float64)
	}

	return r0
}

// GetLatestData provides a mock function with given fields:
func (_m *MQTTTopic) GetLatestData() float64 {
	ret := _m.Called()

	var r0 float64
	if rf, ok := ret.Get(0).(func() float64); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(float64)
	}

	return r0
}

// Subscribe provides a mock function with given fields: _a0
func (_m *MQTTTopic) Subscribe(_a0 mqtt.MQTTComms) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func(mqtt.MQTTComms) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
