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
	testPi := new(RaspberryPi)
	testPi.MQTTClient = new(mockMQTTComms)

	testPi.WiFiConnectButton = new(mockRaspberryPiPinImpl)
	testPi.AirPumpPin = new(mockRaspberryPiPinImpl)
	testPi.GrowLedPin = new(mockRaspberryPiPinImpl)
	testPi.WaterLevelSensor = new(mockUltrasonicSensor)
	testPi.ambientTempSensor = new(mockAmbientTemp)
	testPi.alertChannel = make(chan string)
	return testPi
}

func TestHydroCycle(t *testing.T) {
	mockPi := setupMock()

	t.Run("Load config for device from file", func(t *testing.T) {
		configLocation = "./config/configData.json"
		mockPi.loadConfig()
		if mockPi.ledStartTime != "6:00" {
			t.Errorf("Did not load led start time from file, %s", mockPi.ledStartTime)
		}
		if mockPi.ledOffTime != "23:00" {
			t.Errorf("Did not load end time from file")
		}
	})

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

	//t.Run("Test Water Level sensor", func(t *testing.T) {
	//	mockPi.startSensorCycle()
	//	select {
	//	case x, ok := <- mockPi.alertChannel:
	//		if ok {
	//			fmt.Printf("Value %d was read.\n", x)
	//		} else {
	//			fmt.Println("Channel closed!")
	//			t.Error("Channel should have low level alert")
	//		}
	//	default:
	//			t.Error("Channel should have low level alert")
	//	}
	//})
	
	t.Run("Test that button activates wifi-connect ap", func(t *testing.T) {
		mockPi.WiFiConnectButton.(*mockRaspberryPiPinImpl).stateOfPin = rpio.High
		mockPi.startWifiConnectCycle()
		time.Sleep(time.Second*2)
		mockPi.WiFiConnectButton.(*mockRaspberryPiPinImpl).stateOfPin = rpio.Low
		time.Sleep(time.Second)
	})

	t.Run("Test when there are no alerts coming in", func(t *testing.T) {
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
		var buf bytes.Buffer
		log.SetOutput(&buf)
		defer log.SetOutput(os.Stdout)
		mockPi.monitorAlerts()
		mockPi.alertChannel <- "warning"
		time.Sleep(time.Millisecond)
		out := buf.String()

		if !strings.Contains(out, "WARNING CHECK SYSTEM")   {
			t.Error("Water Level alert not received")
		}
	})
}


