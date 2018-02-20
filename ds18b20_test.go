package main

import "testing"

func TestReadTemp(t *testing.T) {
	waterTempSensor := new(ds18b20)
	temp := waterTempSensor.ReadTemp()
	if temp == nil {
		t.Errorf("Error: Not able to read temp")
	}
}
