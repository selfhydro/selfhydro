package mocks

import (
	"errors"
	"time"
)

type MockMQTTToken struct {
	hasConnectionError bool
}

func (m *MockMQTTToken) Wait() bool {
	return true
}

func (m *MockMQTTToken) WaitTimeout(time time.Duration) bool {
	return true
}

func (m *MockMQTTToken) Error() error {
	if m.hasConnectionError {
		return errors.New("could not connect")
	}
	return nil
}
