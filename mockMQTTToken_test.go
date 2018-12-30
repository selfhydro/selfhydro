package main

import (
	"errors"
	"time"
)

type mockMQTTToken struct {
	hasConnectionError bool
}

func (m *mockMQTTToken) Wait() bool {
	return true
}

func (m *mockMQTTToken) WaitTimeout(time time.Duration) bool {
	return true
}

func (m *mockMQTTToken) Error() error {
	if m.hasConnectionError {
		return errors.New("could not connect")
	}
	return nil
}
