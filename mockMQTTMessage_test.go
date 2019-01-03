package main

type mockMessage struct {
	payload []byte
}

func (msg *mockMessage) Duplicate() bool {
	return true
}

func (msg *mockMessage) Qos() byte {
	return byte(0x1)
}

func (msg *mockMessage) Retained() bool {
	return true
}

func (msg *mockMessage) Topic() string {
	return "/test/"
}

func (msg *mockMessage) MessageID() uint16 {
	return uint16(10)
}

func (msg *mockMessage) Payload() []byte {
	return msg.payload
}

func (msg *mockMessage) Ack() {

}
