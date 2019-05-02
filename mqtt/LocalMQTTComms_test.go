package mqtt

import (
	"log"
	"testing"

	"github.com/bchalk101/selfhydro/mqtt/mocks"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"gotest.tools/assert"
)

func TestShouldCreateNewLocalMQTT(t *testing.T) {
	localMQTT := NewLocalMQTT()
	assert.Equal(t, localMQTT.mqttOptions.Servers[0].String(), "tcp://127.0.0.1:1883")
	assert.Equal(t, localMQTT.mqttOptions.ClientID, "selfhydro-controller")
}

func TestShouldConnectToLocalMQTTBroker(t *testing.T) {
	localMQTT := LocalMQTT{
		mqttClient: &mocks.MockMQTTClient{},
	}
	localMQTT.ConnectDevice()
	assert.Equal(t, localMQTT.mqttClient.IsConnected(), true)
}

func Test_ShouldLogErrorWhenCantConnectToMQTTBroker(t *testing.T) {
	localMQTT := LocalMQTT{
		mqttClient: &mocks.MockMQTTClient{
			HasErrorConnecting: true,
		},
	}
	error := localMQTT.ConnectDevice()
	assert.Equal(t, error.Error(), "could not connect")
}

func Test_ShouldSubscribeToAGivenTopic(t *testing.T) {
	localMQTT := LocalMQTT{
		mqttClient: &mocks.MockMQTTClient{},
	}
	error := localMQTT.SubscribeToTopic("/test/", mockMessageHandler)
	assert.Equal(t, localMQTT.subscribtionTopics[0], "/test/")
	assert.Equal(t, error, nil)
}

func Test_ShouldReturnErrorWhenCantSubscribeToTopic(t *testing.T) {
	localMQTT := LocalMQTT{
		mqttClient: &mocks.MockMQTTClient{
			HasErrorConnecting: true,
		},
	}
	error := localMQTT.SubscribeToTopic("/test/", mockMessageHandler)
	assert.Equal(t, len(localMQTT.subscribtionTopics), 0)
	assert.Equal(t, error.Error(), "could not connect")
}

var mockMessageHandler MQTT.MessageHandler = func(client MQTT.Client, msg MQTT.Message) {
	log.Printf("TOPIC: %s\n", msg.Topic())
	log.Printf("MSG: %s\n", msg.Payload())
}
