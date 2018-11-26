package main

import (
	"log"

	i2c "github.com/d2r2/go-i2c"
	si7021 "github.com/d2r2/go-si7021"
)

type AmbientTempSensor interface {
	GetReadings() (float32, float32)
}

type i2cTempSensor struct {
	address  uint8
	temp     float32
	humidity float32
	sensor   *si7021.Si7021
}

func NewTempSensor() (AmbientTempSensor, error) {

	sensor := i2cTempSensor{}
	sensor.address = 0x40
	return &sensor, nil
}

func (tempSensor *i2cTempSensor) GetReadings() (float32, float32) {
	i2c, err := i2c.NewI2C(0x40, 1)
	if err != nil {
		log.Print(err)
		return 0, 0
	}
	defer i2c.Close()

	sensor := si7021.NewSi7021()
	if err != nil {
		log.Print(err)
		return 0, 0
	}
	rh, temp, err := sensor.ReadRelativeHumidityAndTemperature(i2c)

	if err != nil {
		log.Print(err)
		return 0, 0
	}
	log.Printf("Temprature in celsius = %v*C\n", temp)
	tempSensor.temp = temp
	tempSensor.humidity = rh

	return tempSensor.temp, tempSensor.humidity
}
