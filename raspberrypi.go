package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
	"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/stianeikeland/go-rpio"
	Sensors "github.com/bchalk101/selfhydro/sensors"
)

type Controller interface {
	StopSystem()
	StartHydroponics()
	startSensorCycle()
	startLightCycle()
	startAirPumpCycle()
}

const (
	LowWaterLevel   = "LOW_WATER"
	LowWaterDefault = 1
	Stop            = "STOP"
	StartSlow       = "START-SLOW"
	StartFast       = "START-FAST"
)

var SensorReadingFrequency = 3
var configLocation = "/selfhydro/config/configData.json"
var cpuTempFileLocation = "/sys/class/thermal/thermal_zone0/temp"

type configData struct {
	WaterTempSensorId string `json:"waterTempSensorId"`
	LedOnTime         string `json:"ledOnTime"`
	LedOffTime        string `json:"ledOffTime"`
}

type RaspberryPi struct {
	GrowLedPin           RaspberryPiPin
	WiFiConnectButton    RaspberryPiPin
	WiFiConnectButtonLED RaspberryPiPin
	WaterLevelSensor     UltrasonicSensor
	WaterTempSensor      ds18b20
	ambientTempSensor    Sensors.Sensor
	AirPumpPin           RaspberryPiPin
	MQTTClient           MQTTComms
	alertChannel         chan string
	ledChannel           chan string
	ledStartTime         time.Time
	ledOffTime           time.Time
}

func NewRaspberryPi() *RaspberryPi {
	pi := new(RaspberryPi)

	error := rpio.Open()
	if error != nil {
		log.Fatalf("Could not open rpio pins %v", error.Error())
		os.Exit(1)
	}

	pi.WaterLevelSensor = NewHCSR04Sensor(16, 17)
	pi.WiFiConnectButton = NewRaspberryPiPin(13)
	pi.WiFiConnectButton.SetMode(rpio.Input)
	pi.WiFiConnectButtonLED = NewRaspberryPiPin(14)
	pi.WiFiConnectButtonLED.SetMode(rpio.Output)

	pi.GrowLedPin = NewRaspberryPiPin(19)
	pi.GrowLedPin.SetMode(rpio.Output)

	pi.AirPumpPin = NewRaspberryPiPin(21)
	pi.AirPumpPin.SetMode(rpio.Output)

	pi.MQTTClient = new(mqttComms)

	pi.loadConfig()

	if err := pi.MQTTClient.ConnectDevice(); err != nil {
		pi.handleConnectionError()
	}
	pi.ambientTempSensor = Sensors.NewMCP9808()
	if err := pi.ambientTempSensor.SetupDevice(); err != nil {
		log.Printf("could not setup ambient temp sensor, %s", err)
	}

	pi.alertChannel = make(chan string, 5)

	return pi
}

func (pi *RaspberryPi) loadConfig() {
	data, err := ioutil.ReadFile(configLocation)
	if err != nil {
		log.Print("error loading config data for raspberry pi")
		log.Print(err.Error())
		return
	}

	var configData = new(configData)
	err = json.Unmarshal(data, &configData)
	if err != nil {
		log.Print("error parsing config data")
		log.Print(err.Error())
	}
	startTimeString := configData.LedOnTime
	pi.ledStartTime, err = time.Parse("15:04:05", startTimeString)
	if err != nil {
		log.Printf("error could not parse turn on time, %s", err.Error())
	}
	endTimeString := configData.LedOffTime
	pi.ledOffTime, err = time.Parse("15:04:05", endTimeString)
	if err != nil {
		log.Printf("error could not parse turn off time, %s", err.Error())
	}
	pi.WaterTempSensor.id = configData.WaterTempSensorId
}

func (pi *RaspberryPi) handleConnectionError() {
	log.Print("Could not connect device to IoT platform\n Are you connected to the internet?")
	pi.startWifiConnect()
	for {
		pi.WiFiConnectButtonLED.Toggle()
		time.Sleep(time.Second)
		if pi.WiFiConnectButton.ReadState() == rpio.High {
			break
		}
	}
}

func (pi *RaspberryPi) StartHydroponics() {
	go pi.startSensorCycle()
	pi.subscribeToTopics()
	pi.startLightCycle()
	pi.startAirPumpCycle()
	pi.startWifiConnectCycle()
	pi.monitorAlerts()
}

func (pi *RaspberryPi) subscribeToTopics() {
	deviceID := pi.MQTTClient.GetDeviceID()
	topic := fmt.Sprintf("/devices/%s/commands/#", deviceID)
	pi.MQTTClient.SubscribeToTopic(topic, pi.commandSubscriptionHandler)
}

func (pi *RaspberryPi) commandSubscriptionHandler(client MQTT.Client, message MQTT.Message) {

	log.Printf("TOPIC: %s\n", message.Topic())
	log.Printf("MSG: %s\n", message.Payload())

	switch string(message.Payload()[:]) {
	case "REREAD_TELEMETRY":
		log.Print("message received for REREAD_TELEMETRY")
		pi.readSensorData()
	default:
		log.Print("message unhandled")
		log.Printf("message received: %s", message.Payload())
	}
}

func (pi *RaspberryPi) monitorAlerts() {
	go func() {
		for {
			alert := <-pi.alertChannel
			switch alert {
			case LowWaterLevel:
				log.Print("Water Level is Low")
			default:
				log.Print("WARNING CHECK SYSTEM")
			}
		}
	}()
}

func (pi *RaspberryPi) StopSystem() {
	deviceID := pi.MQTTClient.GetDeviceID()
	topic := fmt.Sprintf("/devices/%s/commands/#", deviceID)

	pi.GrowLedPin.WriteState(rpio.Low)
	pi.AirPumpPin.WriteState(rpio.Low)
	pi.MQTTClient.UnsubscribeFromTopic(topic)
	rpio.Close()
}

func (pi *RaspberryPi) publishState(waterTemp float64, ambientTemp float32, relativeHumidity float32, CPUTemp float64, waterLevel float32) {
	message, err := CreateSensorMessage(waterTemp, ambientTemp, relativeHumidity, CPUTemp, waterLevel)
	if err != nil {
		log.Printf("Error creating sensor message: %s", err)
	}
	fmt.Print(message)
	pi.MQTTClient.publishMessage("/devices/"+pi.MQTTClient.GetDeviceID()+"/events", message)
}

func (pi RaspberryPi) startLightCycle() {
	log.Printf("starting led cycle with on time at %s and off time at %s", pi.ledStartTime, pi.ledStartTime)
	go func() {
		for {
			pi.changeLEDState(pi.ledStartTime, pi.ledOffTime)
			time.Sleep(time.Second * 4)
		}
	}()
}
func (pi RaspberryPi) changeLEDState(turnOnTime time.Time, turnOffTime time.Time) {
	if pi.GrowLedPin.ReadState() != rpio.High && betweenTime(turnOnTime, turnOffTime) {
		log.Printf("Turning on GROW LEDS")
		pi.GrowLedPin.WriteState(rpio.High)
	} else if pi.GrowLedPin.ReadState() == rpio.High && betweenTime(turnOffTime, turnOnTime.Add(time.Hour*24)) {
		log.Printf("Turning off GROW LEDS")
		pi.GrowLedPin.WriteState(rpio.Low)
	}
}
func (pi RaspberryPi) startSensorCycle() {
	for {
		pi.readSensorData()
		time.Sleep(time.Hour * 3)
	}
}

func (pi RaspberryPi) readSensorData() {
	tankOneTemp := pi.WaterTempSensor.ReadTemperature()
	CPUTemp := pi.getCPUTemp()
	waterLevel := pi.checkWaterLevels()
	ambientTemp, _ := pi.ambientTempSensor.GetState()
	pi.publishState(tankOneTemp, ambientTemp, 0, CPUTemp, waterLevel)
}

func (pi RaspberryPi) checkWaterLevels() (level float32) {
	waterLevel := pi.WaterLevelSensor.MeasureDistance()
	fmt.Print(waterLevel)
	if waterLevel <= LowWaterDefault {
		pi.alertChannel <- LowWaterLevel
	}
	log.Printf("Water level is %f", waterLevel)
	return waterLevel
}

func (pi RaspberryPi) getCPUTemp() float64 {

	var temp float64
	data, err := ioutil.ReadFile(cpuTempFileLocation)
	if err != nil {
		log.Printf("Error: Can't read Raspberry Pi CPU Temp")
		return 0.0
	}
	tempData := strings.TrimSuffix(string(data), "\n")

	temp, err = strconv.ParseFloat(string(tempData), 64)
	if err != nil {
		log.Println("could not read cpu temp")
		log.Println("system will now shut down so as not to destroy the controller")
		panic(err)
	}
	log.Printf("CPU Temp: %v", temp/1000)
	return temp / 1000

}

func (pi RaspberryPi) startAirPumpCycle() {
	go func() {
		for {
			pi.airPumpCycle(time.Minute*30, time.Hour*2)
		}
	}()
}
func (pi RaspberryPi) airPumpCycle(airPumpOnDuration time.Duration, airPumpOffDuration time.Duration) {
	log.Printf("Turning on air pump")
	pi.AirPumpPin.WriteState(rpio.High)
	time.Sleep(airPumpOnDuration)
	log.Printf("Turning off air pump")
	pi.AirPumpPin.WriteState(rpio.Low)
	time.Sleep(airPumpOffDuration)
}
func (pi *RaspberryPi) startWifiConnectCycle() {
	go func() {
		for {
			pi.checkIfWifiButtonIsPressed()
		}
	}()
}
func (pi *RaspberryPi) checkIfWifiButtonIsPressed() {
	if pi.WiFiConnectButton.ReadState() == rpio.High {
		startTime := time.Now()
		for {
			if pi.WiFiConnectButton.ReadState() == rpio.Low {
				break
			}
		}
		if time.Since(startTime) >= time.Second*2 {

			pi.startWifiConnect()
		}
	}
}
func (pi RaspberryPi) startWifiConnect() {
	binary, lookErr := exec.LookPath("wifi-connect")
	if lookErr != nil {
		log.Printf("Error: Could not find wifi-connect")
	}
	args := []string{"wifi-connect", "-s=Selfhydro Connect"}
	env := os.Environ()
	execErr := syscall.Exec(binary, args, env)
	if execErr != nil {
		log.Printf("Error: Could not start wifi-connect")
	}
}

func (pi *RaspberryPi) flashLED() {
	go func() {
		for {
			pi.WiFiConnectButtonLED.Toggle()
			time.Sleep(time.Millisecond * 50)

		}

	}()
}

func betweenTime(startTime time.Time, endTime time.Time) bool {
	currentTimeString := time.Now().Format("15:04:05")
	currentTime, _ := time.Parse("15:04:05", currentTimeString)
	if currentTime.After(startTime) && currentTime.Before(endTime) {
		return true
	}
	return false
}
