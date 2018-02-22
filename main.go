package main

import (
	"github.com/stianeikeland/go-rpio"
	"os"
	"time"
	"os/signal"
	"log"
)

type Time struct {
	Hh int // Hours.
	Mm int // Minutes.
	Ss int // Seconds.
}

func main() {
	f, err := os.OpenFile("SelfHydroLogs", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	log.SetOutput(f)

	log.Println("Starting up SelfHydro")

	error := rpio.Open()
	if error != nil {
		os.Exit(1)
	}
	defer rpio.Close()

	controller := NewController()

	sigs := make(chan os.Signal, 1)

	signal.Notify(sigs, os.Interrupt)
	signal.Notify(sigs, os.Kill)

	go func() {
		s := <-sigs
		log.Println("Exiting program...")
		log.Println("RECEIVED SIGNAL: ", s)

		controller.turnOffGrowLed()
		controller.WaterPumpPin.Low()
		os.Exit(0)
	}()

	controller.StartSensorCycle()
	controller.StartLightCycle()
	controller.StartWaterCycle()
	controller.StartAirPumpCycle()

	for {
		time.Sleep(time.Second)
	}

}

func NewController() *RaspberryPi {
	pi := new(RaspberryPi)

	pi.GrowLedPin = rpio.Pin(19)
	pi.GrowLedState = false

	pi.WaterPumpPin = rpio.Pin(20)
	pi.WaterPumpState = false

	pi.WaterLevelSensor = rpio.Pin(4)
	pi.WaterLevelSensor.Input()

	pi.GrowLedPin.Mode(rpio.Output)
	pi.WaterPumpPin.Mode(rpio.Output)

	pi.WaterTempSensor.id = "28-0316838ca7ff"

	pi.AirPump = rpio.Pin(21)
	pi.AirPump.Mode(rpio.Output)

	return pi
}

