package main

import (
	"fmt"
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
	SetupCloseHandler()
	sh := selfhydro{}
	sh.Setup()
	err = sh.Start()
	if err != nil {
		log.Fatal(err)
	}
	select {}
}

func SetupCloseHandler() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println("\r- Ctrl+C pressed in Terminal")
		os.Exit(0)
	}()
}

func panicHandler(output string) {
	log.Printf("The child panicked:\n\n%s\n", output)
	os.Exit(1)
}
