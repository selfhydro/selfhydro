package main

import (
	"errors"
	"io/ioutil"
	"strings"
	"strconv"
	"log"
)

type ds18b20 struct {
	id string
}

var ErrReadSensor = errors.New("failed to read sensor temperature")

func (ds ds18b20) ReadTemperature() float64{

	temp, err := ds.getTemp(ds.id)

	if err != nil {
		log.Printf("Error: Cant read temp of water tank")
		return 0.0
	}
	log.Printf("Water temperature: %.2f°C\n", temp)

	return temp
}

func (ds18b20) getTemp(sensor string) (float64, error) {
	data, err := ioutil.ReadFile("/sys/bus/w1/devices/" + sensor + "/w1_slave")
	if err != nil {
		return 0.0, ErrReadSensor
	}

	raw := string(data)

	i := strings.LastIndex(raw, "t=")
	if i == -1 {
		return 0.0, ErrReadSensor
	}

	c, err := strconv.ParseFloat(raw[i+2:len(raw)-1], 64)
	if err != nil {
		return 0.0, ErrReadSensor
	}

	return c / 1000.0, nil
}
