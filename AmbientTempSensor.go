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
		log.Fatal(err)
	}
	defer i2c.Close()

	sensor := si7021.NewSi7021()
	if err != nil {
		log.Fatal(err)
	}
	rh, err := sensor.ReadRelativeHumidityMode1(i2c)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Relative humidity = %v%%\n", rh)
	tempSensor.humidity = rh

	temp, err := sensor.ReadTemperatureCelsiusMode1(i2c)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Temprature in celsius = %v*C\n", temp)
	tempSensor.temp = temp

	return tempSensor.temp, tempSensor.humidity
}
