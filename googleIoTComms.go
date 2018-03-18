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
	"bytes"
)

type SensorMessage struct {
	DeviceId int `json:"deviceId"`
	UnitOneWaterTemp float32 `json:"unitOneWaterTemp"`
	UnitTwoWaterTemp float32 `json:"unitTwoWaterTemp"`
	PiCPUTemp float32 `json:"piCPUTemp"`
	LEDState bool
	Time string `json:"time"`
}

const (
	location  = "asia-east1"
	projectId = "selfhydro-197504"
	registry  = "raspberry-pis"
	device    = "original-hydro"
)

type MQTTComms struct {
	client MQTT.Client
}

const (
	HYDRO_EVENTS_TOPIC = "/devices/"+device+"/events"
)

var f MQTT.MessageHandler = func(client MQTT.Client, msg MQTT.Message) {
	fmt.Printf("TOPIC: %s\n", msg.Topic())
	fmt.Printf("MSG: %s\n", msg.Payload())
}

func (mqtt *MQTTComms) authenticateDevice() {

	tokenString, _ := createJWTToken(projectId)

	opts := MQTT.NewClientOptions().AddBroker("ssl://mqtt.googleapis.com:8883")

	clientId := "projects/" + projectId + "/locations/" + location + "/registries/" + registry + "/devices/" + device

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
func (mqtt *MQTTComms) publishMessage(topic string, message []byte) {
	log.Printf("Sending: %v", bytes.NewBuffer(message))
	token := mqtt.client.Publish(topic, 0, false, bytes.NewBuffer(message))
	response := token.Wait()
	log.Printf("Response: %v",response)
}

func createJWTToken(projectId string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(time.Hour * 24).Unix(),
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

func CreateSensorMessage(tempUnitOne float32, tempUnitTwo float32, piCPUTemp float32, LEDState bool ) ([]byte, error) {
	m := SensorMessage{rand.Int(),tempUnitOne, tempUnitTwo, piCPUTemp, LEDState,time.Now().Format("20060102150405")}
	return json.Marshal(m)
}
