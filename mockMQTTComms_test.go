package main

type mockMQTTComms struct {
}

func (m *mockMQTTComms) ConnectDevice() error {
	return nil
}

func (m *mockMQTTComms) publishMessage(topic string, message string) {

}

func (m *mockMQTTComms) GetDeviceID() string {
	return ""
}
