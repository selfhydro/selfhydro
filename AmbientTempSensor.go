package main

import (
	"github.com/learnaddict/mcp9808"
	"log"
	"errors"
)

type AmbientTempSensor interface {
	GetTemp() float32
}

type mcp9808Sensor struct {
	address []uint8
	temp    float32
}

func NewMCP9808Sensor() (AmbientTempSensor, error) {

	sensor := mcp9808Sensor{}
	sensor.address = mcp9808.Find()
	if len(sensor.address) == 0 {
		log.Print("Can not find mcp9808 sensor")
		return nil, errors.New("no MCP9808 sensor found")
	}
	return &sensor, nil
}

func (sensor *mcp9808Sensor) GetTemp() float32 {
	for _, a := range sensor.address {
		temp, err := mcp9808.Read(a)
		if err != nil {
			log.Print(err)
			continue
		}
		sensor.temp = temp
	}

	return sensor.temp
}
