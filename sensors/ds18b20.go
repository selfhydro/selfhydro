package sensors

import (
	"errors"
	"io/ioutil"
	"log"
	"path/filepath"
	"strconv"
	"strings"
)

type DS18b20 struct {
	Id string
}

var dataDirectory = "/sys/bus/w1/devices/"
var getReadDir = ioutil.ReadDir
var sensorError = errors.New("failed to read sensor temperature")

func NewDS18B20(id string) Sensor {
	return &DS18b20{Id: id}
}

func (ds *DS18b20) SetupDevice() error {
	files, err := getReadDir(dataDirectory)
	if err != nil {
		log.Printf("error reading directory: %v", err)
		return errors.New("error finding directory for sensor")
	}
	for _, file := range files {
		if !strings.Contains(file.Name(), "w1") {
			ds.Id = file.Name()
		}
	}
	return nil
}

func (ds DS18b20) GetState() (float32, error) {
	temp, err := ds.getTemp(ds.Id)
	if err != nil {
		log.Printf("Error: Cant read temp of water tank")
		log.Print(err.Error())
		return 0.0, err
	}
	log.Printf("Water temperature: %.2f°C\n", temp)
	return temp, nil
}

func (DS18b20) getTemp(sensor string) (float32, error) {
	sensorDirectory := filepath.Join(dataDirectory, sensor, "w1_slave")
	data, err := ioutil.ReadFile(sensorDirectory)
	if err != nil {
		return 0.0, sensorError
	}
	raw := string(data)
	i := strings.LastIndex(raw, "t=")
	if i == -1 {
		return 0.0, sensorError
	}
	c, err := strconv.ParseFloat(raw[i+2:len(raw)-1], 64)
	if err != nil {
		return 0.0, sensorError
	}
	return float32(c / 1000.0), nil
}
