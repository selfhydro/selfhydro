package main

import (
	"testing"
	"github.com/stianeikeland/go-rpio"
	"time"
	"bytes"
	"log"
	"strings"
	"os"
)


func setupMock() *RaspberryPi {
	mockPi := new(RaspberryPi)
	mockPi.MQTTClient = new(mockMQTTComms)

	mockPi.WiFiConnectButton = new(mockRaspberryPiPinImpl)
	mockPi.AirPumpPin = new(mockRaspberryPiPinImpl)
	mockPi.GrowLedPin = new(mockRaspberryPiPinImpl)
	mockPi.tankOneWaterLevelSensor = new(mockSensor)
	mockPi.alertChannel = make(chan string)
	return mockPi
}

func TestHydroCycle(t *testing.T) {
	mockPi := setupMock()
	t.Run("Testing Grow LEDS", func(t *testing.T) {
		startTimeString := time.Now().Add(-time.Minute).Format("15:04:05")
		startTime, _ := time.Parse("15:04:05", startTimeString)

		offTimeString := time.Now().Add(time.Minute).Format("15:04:05")
		offTime, _ := time.Parse("15:04:05", offTimeString)

		mockPi.changeLEDState(startTime, offTime)
		if mockPi.GrowLedPin.ReadState() != rpio.High {
			t.Errorf("Error: GrowLED not turned on")
		}
	})
	
	t.Run("Test Air Pump cycle", func(t *testing.T) {
		mockPi.airPumpCycle(time.Second, time.Second)
		if mockPi.AirPumpPin.ReadState() != rpio.Low {
			t.Errorf("Error: Airpump was not turned on")
		}
	})

	t.Run("Test Water Level sensor", func(t *testing.T) {
		mockPi.tankOneWaterLevelSensor.(*mockSensor).sensorState = rpio.High
		mockPi.startSensorCycle()
		if <-mockPi.alertChannel != LowWaterLevel {
			t.Error("Channel should have low level alert")
		}

	})
	
	t.Run("Test that button activates wifi-connect ap", func(t *testing.T) {
		mockPi.WiFiConnectButton.(*mockRaspberryPiPinImpl).stateOfPin = rpio.High
		mockPi.startWifiConnectCycle()
		time.Sleep(time.Second*2)
		mockPi.WiFiConnectButton.(*mockRaspberryPiPinImpl).stateOfPin = rpio.Low
		time.Sleep(time.Second)


	})

	t.Run("Test when there are no alerts coming in", func(t *testing.T) {
		mockPi.tankOneWaterLevelSensor.(*mockSensor).sensorState = rpio.Low
		var buf bytes.Buffer
		log.SetOutput(&buf)
		defer log.SetOutput(os.Stdout)
		mockPi.monitorAlerts()
		mockPi.startSensorCycle()
		time.Sleep(time.Millisecond)
		out := buf.String()

		if strings.Contains(out, "Water Level is Low")   {
			t.Error("Water Level alert not received")
		}

	})

	t.Run("Alerts should be logged when ever they come in", func(t *testing.T){
		mockPi.tankOneWaterLevelSensor.(*mockSensor).sensorState = rpio.High
		var buf bytes.Buffer
		log.SetOutput(&buf)
		defer log.SetOutput(os.Stdout)
		mockPi.monitorAlerts()
		mockPi.startSensorCycle()
		time.Sleep(time.Millisecond)
		out := buf.String()

		if !strings.Contains(out, "Water Level is Low")   {
			t.Error("Water Level alert not received")
		}
	})
}


