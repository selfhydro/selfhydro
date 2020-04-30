package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/mitchellh/panicwrap"
)

type Time struct {
	Hh int // Hours.
	Mm int // Minutes.
	Ss int // Seconds.
}

func main() {
	f, err := os.OpenFile("./selfhydro-logs", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
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

	sh := selfhydro{}
	sh.Setup()
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs)
	sh.Start()
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
