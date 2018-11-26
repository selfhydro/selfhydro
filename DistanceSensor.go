package main

import (
	"log"

	i2c "github.com/d2r2/go-i2c"
	vl53l0x "github.com/d2r2/go-vl53l0x"
)

type DistanceSensor interface {
	MeasureDistance() (float32, error)
}

type Vl53l0xSensor struct {
	distance float32
}

func NewDistanceSensor() DistanceSensor {
	sensor := new(Vl53l0xSensor)
	return sensor
}

func (ds *Vl53l0xSensor) MeasureDistance() (float32, error) {
	i2c, err := i2c.NewI2C(0x29, 1)
	if err != nil {
		log.Print(err)
		return 0, err
	}
	defer i2c.Close()

	sensor := vl53l0x.NewVl53l0x()
	err = sensor.Reset(i2c)
	if err != nil {
		log.Print(err)
		return 0, err
	}
	err = sensor.Init(i2c)
	if err != nil {
		log.Print(err)
		return 0, err
	}
	rng, err := sensor.ReadRangeSingleMillimeters(i2c)
	if err != nil {
		log.Fatal(err)
		return 0, err
	}
	log.Printf("Measured range = %v mm", rng)
	return ds.distance, nil
}
