package main

import (
	"github.com/stianeikeland/go-rpio"
	"os"
	"time"
	"os/signal"
	"log"
	//"github.com/morus12/dht22"
	"github.com/morus12/dht22"
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

	go func() {
		s := <-sigs
		log.Println("Exiting program...")
		log.Println("RECEIVED SIGNAL: ", s)

		controller.turnOffGrowLed()
		controller.WaterPumpPin.Low()
		os.Exit(0)
	}()

	go func() {
		sensor := dht22.New("GPIO_17")

		for {
			log.Println("Reading temp")
			temperature, errtemp := sensor.Temperature()
			humidity, errhumid := sensor.Humidity()

			if errtemp != nil {
				log.Println("Error: Cant Read Temp ", errtemp.Error())
			} else {
				log.Println("Temp: ", temperature)

			}

			if errhumid != nil {
				log.Println("Error: Cant Read Humidity ", errhumid.Error())
			} else {
				log.Println("Humidity: ", humidity)
			}

			time.Sleep(time.Minute)
		}

	}()


	turnOnTime, _ := time.Parse("15:04:05", "04:45:00")
	turnOffTime, _ := time.Parse("15:04:05", "23:45:00")

	go func() {
		for {
			log.Printf("GrowLED state: %v", controller.GrowLedState)
			if !controller.GrowLedState && betweenTime(turnOnTime, turnOffTime) {
				log.Printf("Turning on GROW LEDS")
				controller.turnOnGrowLed()
			} else if controller.GrowLedState && betweenTime(turnOffTime, turnOnTime.Add(time.Hour*24)) {
				log.Printf("Turning off GROW LEDS")
				controller.turnOffGrowLed()
			}
			time.Sleep(time.Minute*1)
		}


	}()


	go func() {
		for {
			controller.turnOnWaterPump()
			time.Sleep(time.Second*5)
			controller.turnOffWaterPump()
			time.Sleep(time.Minute*120)
		}

	} ()

	for {
		time.Sleep(time.Second)
	}

}

func NewController() *RaspberryPi {
	pi := &RaspberryPi{
		GrowLedPin:   rpio.Pin(19),
		GrowLedState: false,
		WaterPumpPin: rpio.Pin(20),
		WaterPumpState: false,
	}

	pi.GrowLedPin.Mode(rpio.Output)
	pi.WaterPumpPin.Mode(rpio.Output)
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
