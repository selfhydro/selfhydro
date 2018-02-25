package main

import (
	"testing"
)

func setup() {

}

func TestTemperature(t *testing.T) {
	dht22 := new(DHT22)
	temp, _ := dht22.Temperature()
	if temp != 24 {
		t.Error("Error: Not calculating temp correctlty")
	}

}
