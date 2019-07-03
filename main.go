package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/mitchellh/panicwrap"
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

	exitStatus, err := panicwrap.BasicWrap(panicHandler)
	if err != nil {
		panic(err)
	}
	if exitStatus >= 0 {
		os.Exit(exitStatus)
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

	sh := selfhydro{}
	waterPump := NewWaterPump(18)
	airPump := NewAirPump(21)
	growLight := NewGrowLight(19)
	err = sh.Setup(waterPump, airPump, growLight)
	if err != nil {
		log.Fatal(err)
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs)
	sh.Start()
	defer sh.StopSystem()
	handleExit(<-sigs)
}

func panicHandler(output string) {
	log.Printf("The child panicked:\n\n%s\n", output)
	os.Exit(1)
}

func handleExit(signal os.Signal) {
	if signal != syscall.SIGPIPE {
		log.Println("Exiting program...")
		log.Println("RECEIVED SIGNAL: ", signal)

		os.Exit(0)
	}
}
