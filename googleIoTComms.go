package main

import (
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/dgrijalva/jwt-go"
	"time"
	"fmt"
	"os"
	"log"
	"io/ioutil"
	"encoding/json"
	"math/rand"
)

type SensorMessage struct {
	UnitOneWaterTemp float64 `json:"unitOneWaterTemp"`
	UnitTwoWaterTemp float64 `json:"unitTwoWaterTemp"`
	PiCPUTemp        float64 `json:"piCPUTemp"`
	Time             string  `json:"time"`
}

const (
	location   = "asia-east1"
	projectId  = "selfhydro-197504"
	registryId = "raspberry-pis"
	deviceId   = "original-hydro"
)

type MQTTComms struct {
	client MQTT.Client
}

const (
	EVENTSTOPIC = "/devices/" + deviceId + "/events"
	JWTEXPIRYINHOURS = 6
)

var f MQTT.MessageHandler = func(client MQTT.Client, msg MQTT.Message) {
	fmt.Printf("TOPIC: %s\n", msg.Topic())
	fmt.Printf("MSG: %s\n", msg.Payload())
}

func (mqtt *MQTTComms) connectDevice() {
	mqtt.authenticateDevice()
	timerTillRefresh := time.NewTimer(JWTEXPIRYINHOURS * time.Hour)
	go func() {
		<-timerTillRefresh.C
		fmt.Println("Refreshing JWT Token and reconneting")
		mqtt.client.Disconnect(200)
		mqtt.authenticateDevice()
	}()
}

func (mqtt *MQTTComms) authenticateDevice() {

	tokenString, _ := createJWTToken(projectId)

	opts := MQTT.NewClientOptions().AddBroker("ssl://mqtt.googleapis.com:8883")

	clientId := "projects/" + projectId + "/locations/" + location + "/registries/" + registryId + "/devices/" + deviceId

	opts.SetClientID(clientId)
	opts.SetDefaultPublishHandler(f)
	opts.SetPassword(tokenString)
	opts.SetProtocolVersion(4)
	opts.SetUsername("unused")

	mqtt.client = MQTT.NewClient(opts)
	if token := mqtt.client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

}
func (mqtt *MQTTComms) subscribeToTopic(topic string) {
	if token := mqtt.client.Subscribe(topic, 0, nil); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}
}
func (mqtt *MQTTComms) unsubscribeFromTopic(topic string) {
	if token := mqtt.client.Unsubscribe(topic); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}
	mqtt.client.Disconnect(250)
}
func (mqtt *MQTTComms) publishMessage(topic string, message string) {
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

	fmt.Println(tokenString, err)
	return tokenString, err
}

func CreateSensorMessage(tempUnitOne float64, tempUnitTwo float64, piCPUTemp float64) (string, error) {
	m := SensorMessage{tempUnitOne, tempUnitTwo, piCPUTemp, time.Now().Format("20060102150405")}
	jsonMsg, err := json.Marshal(m)
	return string(jsonMsg), err
}
