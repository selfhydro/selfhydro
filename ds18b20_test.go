package main

import (
	"testing"
	"fmt"
)


func TestGetID(t *testing.T) {
	tempSensor := new(ds18b20)
	dataDirectory = "testdata"
	tempSensor.GetID()
	fmt.Print(tempSensor.id)
	if tempSensor.id != "testSensor" {
		t.Error("Error: Did not find correct ID for sensor")
	}
}

func TestReadTemp(t *testing.T) {
	waterTempSensor := new(ds18b20)
	waterTempSensor.id = "testSensor"
	dataDirectory = "testdata"
	temp := waterTempSensor.ReadTemperature()
	if temp != 10.00 {
		t.Errorf("Error: Not able to read temp")
	}
}
