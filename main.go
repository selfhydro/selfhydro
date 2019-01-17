package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	rpio "github.com/stianeikeland/go-rpio"
)

type Time struct {
	Hh int // Hours.
	Mm int // Minutes.
	Ss int // Seconds.
}

func main() {
	f, err := os.OpenFile("/selfhydro/selfhydro-logs", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()
	log.SetOutput(f)
	log.Println("Starting up SelfHydro")

	error := rpio.Open()
	if error != nil {
		log.Fatalf("Could not open rpio pins %v", error.Error())
		os.Exit(1)
	}
	defer rpio.Close()

	controller := NewRaspberryPi()
	sh := selfhydro{}
	waterPump := NewWaterPump(18)
	airPump := NewAirPump(21)
	err = sh.Setup(waterPump, airPump)
	if err != nil {
		log.Fatal(err)
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs)
	go func() {
		s := <-sigs
		if s != syscall.SIGPIPE {
			log.Println("Exiting program...")
			log.Println("RECEIVED SIGNAL: ", s)

			controller.StopSystem()

			os.Exit(0)
		}

	}()
	sh.Start()
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
