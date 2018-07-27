package main

import (
	"testing"
	"github.com/stianeikeland/go-rpio"
	"time"
)

func TestReadDistance(t *testing.T) {
	t.Run("Should get distance of object", func(t *testing.T) {
		hc := new(HCSR04)
		hc.echoPin = new(mockRaspberryPiPinImpl)
		hc.pingPin = new(mockRaspberryPiPinImpl)
		hc.pingPin.(*mockRaspberryPiPinImpl).stateOfPin = rpio.Low

		go func() {
			for i := 0; i < 10000 && hc.pingPin.ReadState() == rpio.Low; i++ {

			}
			for ; hc.pingPin.ReadState() == rpio.High; {

			}
			time.Sleep(time.Microsecond * 2)
			hc.echoPin.(*mockRaspberryPiPinImpl).stateOfPin = rpio.High
			time.Sleep(time.Microsecond * 58)
			hc.echoPin.(*mockRaspberryPiPinImpl).stateOfPin = rpio.Low
		}()
		distance := hc.MeasureDistance()

		if distance > 2 || distance < 1 {
			t.Errorf("Distance not measured as expected, was %f but expected 1", distance)
		}
	})
}
