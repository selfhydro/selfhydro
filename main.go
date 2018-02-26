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
	controller := NewRaspberryPi()

	error := rpio.Open()
	if error != nil {
		os.Exit(1)
	}
	defer rpio.Close()

	handleExit(controller)

	controller.StartHydroponics()

	for {
		time.Sleep(time.Second)
	}

}

func handleExit(controller *RaspberryPi) {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs)
	go func() {
		s := <-sigs
		log.Println("Exiting program...")
		log.Println("RECEIVED SIGNAL: ", s)

		controller.StopSystem()

		os.Exit(0)
	}()
}
