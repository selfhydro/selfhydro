package main

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"testing"
)

func TestGetID(t *testing.T) {
	tempSensor := new(ds18b20)
	dataDirectory = "./testdata/sensor_test/"
	tempSensor.GetID()
	if tempSensor.id != "testSensor" {
		t.Error("Error: Did not find correct ID for sensor")
	}
}

func TestShouldLogErrorIfCantGetId(t *testing.T) {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer log.SetOutput(os.Stdout)
	dataDirectory = ""
	getReadDir = mockReadDir
	tempSensor := ds18b20{}
	tempSensor.GetID()
	out := buf.String()
	fmt.Print(out)
	if strings.Contains(out, "error reading directory:") {
		t.Error("test failed: does not error if file directory not set correctly")
	}
}

func mockReadDir(dir string) ([]os.FileInfo, error) {
	if dir == "" {
		return nil, errors.New("error")
	}
	return nil, nil
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
