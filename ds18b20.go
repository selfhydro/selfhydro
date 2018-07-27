package main

import (
	"errors"
	"io/ioutil"
	"strings"
	"strconv"
	"log"
	"path/filepath"
)

type ds18b20 struct {
	id string
}

var dataDirectory = "sys/bus/w1/devices"

var ErrReadSensor = errors.New("failed to read sensor temperature")

//Assumes that there is only one 1-wire device connected
func (ds *ds18b20) GetID() {
	files, err := ioutil.ReadDir(dataDirectory)
	if err != nil {
		log.Printf("error reading directory: %v", err)
	}

	for _, file := range files {
		if !strings.Contains(file.Name(), "w1") {
			ds.id = file.Name()
			log.Print(ds.id)
		}
	}
}

func (ds ds18b20) ReadTemperature() float64 {

	temp, err := ds.getTemp(ds.id)

	if err != nil {
		log.Printf("Error: Cant read temp of water tank")
		return 0.0
	}
	log.Printf("Water temperature: %.2f°C\n", temp)

	return temp
}

func (ds18b20) getTemp(sensor string) (float64, error) {
	sensorDirectory := filepath.Join(dataDirectory, sensor, "w1_slave")
	data, err := ioutil.ReadFile(sensorDirectory)
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
