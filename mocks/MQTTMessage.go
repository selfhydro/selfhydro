package mocks

type MockMQTTMessage struct {
	ReceivedPayload []byte
}

func (msg *MockMQTTMessage) Duplicate() bool {
	return true
}

func (msg *MockMQTTMessage) Qos() byte {
	return byte(0x1)
}

func (msg *MockMQTTMessage) Retained() bool {
	return true
}

func (msg *MockMQTTMessage) Topic() string {
	return "/test/"
}

func (msg *MockMQTTMessage) MessageID() uint16 {
	return uint16(10)
}

func (msg *MockMQTTMessage) Payload() []byte {
	return msg.ReceivedPayload
}

func (msg *MockMQTTMessage) Ack() {

}
