package main

import (
	"github.com/stianeikeland/go-rpio"
	"os"
	"time"
	"os/signal"
	"log"
	"github.com/d2r2/go-dht"
	"fmt"
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

	turnOnTime, _ := time.Parse("15:04:05", "04:45:00")
	turnOffTime, _ := time.Parse("15:04:05", "23:45:00")

	go func() {
		for {
			temperature, humidity, retried, err :=
				dht.ReadDHTxxWithRetry(dht.DHT22, 17, true, 10)
			if err != nil {
				log.Fatal(err)
			}
			log.Printf("Ambient Temperature = %v*C, Humidity = %v%% (retried %d times)\n",
				temperature, humidity, retried)
			controller.getWaterTemp()
			time.Sleep(time.Minute)
		}


	}()

	go func() {
		for {
			if !controller.GrowLedState && betweenTime(turnOnTime, turnOffTime) {
				log.Printf("Turning on GROW LEDS")
				controller.turnOnGrowLed()
			} else if controller.GrowLedState && betweenTime(turnOffTime, turnOnTime.Add(time.Hour*24)) {
				log.Printf("Turning off GROW LEDS")
				controller.turnOffGrowLed()
			}
			time.Sleep(time.Minute * 1)
		}

	}()

	go func() {
		for {
			controller.turnOnWaterPump()
			time.Sleep(time.Second * 5)
			controller.turnOffWaterPump()
			time.Sleep(time.Minute * 120)
		}
	}()

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

	pi.GrowLedPin.Mode(rpio.Output)
	pi.WaterPumpPin.Mode(rpio.Output)

	pi.WaterTempSensor.id = "28-0316838ca7ff"

	return pi
}

func betweenTime(startTime time.Time, endTime time.Time) bool {
	currentTimeString := time.Now().Format("15:04:05")
	currentTime, _ := time.Parse("15:04:05", currentTimeString)
	if currentTime.After(startTime) && currentTime.Before(endTime) {
		return true
	}

	return false
}
