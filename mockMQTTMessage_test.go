package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

type mockMessage struct {
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
	return float32ToByte(float32(2.24))
}

func float32ToByte(f float32) []byte {
	var buf bytes.Buffer
	err := binary.Write(&buf, binary.LittleEndian, f)
	if err != nil {
		fmt.Println("binary.Write failed:", err)
	}
	return buf.Bytes()
}

func (msg *mockMessage) Ack() {

}
