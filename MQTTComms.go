package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	MQTT "github.com/eclipse/paho.mqtt.golang"
)

type MQTTComms interface {
	ConnectDevice() error
	publishMessage(topic string, message string)
	GetDeviceID() string
	SubscribeToTopic(string, MQTT.MessageHandler)
	UnsubscribeFromTopic(topic string)
}

type SensorMessage struct {
	WaterTemp        float32 `json:"waterTemp"`
	AmbientTemp      float32 `json:"ambientTemp"`
	RelativeHumidity float32 `json:"relativeHumidity"`
	PiCPUTemp        float64 `json:"piCPUTemp"`
	WaterLevel       float32 `json:"waterLevel"`
	Time             string  `json:"time"`
}

type MQTTDetail struct {
	Location   string `json:"location"`
	ProjectID  string `json:"projectID"`
	RegistryID string `json:"registryID"`
	DeviceID   string `json:"deviceID"`
}

type mqttComms struct {
	client      MQTT.Client
	mqttDetails MQTTDetail
}

const (
	//EVENTSTOPIC      = "/devices/" + %s + "/events"
	JWTEXPIRYINHOURS = 6
)

var subscriptionTopic string
var subscribtionHandler MQTT.MessageHandler

var f MQTT.MessageHandler = func(client MQTT.Client, msg MQTT.Message) {
	log.Printf("TOPIC: %s\n", msg.Topic())
	log.Printf("MSG: %s\n", msg.Payload())
}

var subscribeHandler MQTT.MessageHandler = func(client MQTT.Client, message MQTT.Message) {
	fmt.Printf("MSG: %s\n", message.Payload())
}

func (mqtt *mqttComms) ConnectDevice() error {
	mqtt.loadMQTTConfig()
	if err := mqtt.authenticateDevice(); err != nil {
		return err
	}
	timerTillRefresh := time.NewTimer(JWTEXPIRYINHOURS * time.Hour)
	go func() {
		for {
			<-timerTillRefresh.C
			fmt.Println("Refreshing JWT Token and reconneting")
			mqtt.client.Disconnect(200)
			mqtt.authenticateDevice()
			mqtt.resubscribeToTopics()
			timerTillRefresh = time.NewTimer(JWTEXPIRYINHOURS * time.Hour)
		}
	}()
	return nil
}

func (mqtt *mqttComms) GetDeviceID() string {
	return mqtt.mqttDetails.DeviceID
}

func (mqtt *mqttComms) resubscribeToTopics() {
	mqtt.SubscribeToTopic(subscriptionTopic, subscribtionHandler)
}

func (mqtt *mqttComms) SubscribeToTopic(topic string, callback MQTT.MessageHandler) {
	log.Println("suibscribing to topic ", topic)
	subscriptionTopic = topic
	subscribtionHandler = callback
	if token := mqtt.client.Subscribe(topic, 1, callback); token.Wait() && token.Error() != nil {
		log.Println("error subscribing to topic ", topic)
		log.Println(token.Error())
	}
}

func (mqtt *mqttComms) loadMQTTConfig() {
	file, err := ioutil.ReadFile("/selfhydro/config/googleCloudIoTConfig.json")
	if err != nil {
		log.Printf("Could not find config file for Google Core IoT connection")
		log.Print(err)
	}

	err = json.Unmarshal(file, &mqtt.mqttDetails)
	if err != nil {
		panic(err)
	}
}

func (mqtt *mqttComms) authenticateDevice() error {

	tokenString, _ := createJWTToken(mqtt.mqttDetails.ProjectID)

	opts := MQTT.NewClientOptions().AddBroker("ssl://mqtt.googleapis.com:8883")

	clientId := "projects/" + mqtt.mqttDetails.ProjectID + "/locations/" + mqtt.mqttDetails.Location + "/registries/" + mqtt.mqttDetails.RegistryID + "/devices/" + mqtt.mqttDetails.DeviceID
	fmt.Print(clientId)
	opts.SetClientID(clientId)
	opts.SetDefaultPublishHandler(f)
	opts.SetPassword(tokenString)
	opts.SetProtocolVersion(4)
	opts.SetUsername("unused")

	mqtt.client = MQTT.NewClient(opts)
	if token := mqtt.client.Connect(); token.Wait() && token.Error() != nil {
		if token.Error().Error() == "" {

		} else {

			log.Print(token.Error())
			return token.Error()
		}
	}

	return nil
}

func (mqtt *mqttComms) UnsubscribeFromTopic(topic string) {
	if token := mqtt.client.Unsubscribe(topic); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}
	mqtt.client.Disconnect(250)
}
func (mqtt *mqttComms) publishMessage(topic string, message string) {
	if mqtt.client.IsConnected() {

		log.Printf("Sending: %v", message)
		token := mqtt.client.Publish(topic, 0, false, message)
		response := token.Wait()
		log.Printf("Response: %v", response)
	} else {
		log.Printf("Disconnected from google cloud")
	}
}

func createJWTToken(projectId string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(time.Hour * JWTEXPIRYINHOURS).Unix(),
		"aud": projectId,
	})

	file, err := os.Open("/selfhydro/rsa_private.pem") // For read access.
	if err != nil {
		log.Fatal(err)
	}

	key, _ := ioutil.ReadFile(file.Name())

	rsaPrivateKey, _ := jwt.ParseRSAPrivateKeyFromPEM(key)

	tokenString, err := token.SignedString(rsaPrivateKey)
	return tokenString, err
}

func CreateSensorMessage(waterTemp float32, ambientTemp float32, relativeHumidity float32, piCPUTemp float64, waterLevel float32) (string, error) {
	m := SensorMessage{waterTemp, ambientTemp, relativeHumidity, piCPUTemp, waterLevel, time.Now().Format("20060102150405")}
	jsonMsg, err := json.Marshal(m)
	return string(jsonMsg), err
}
