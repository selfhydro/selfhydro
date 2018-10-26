package main

type mockAmbientTemp struct {
}

func (sensor mockAmbientTemp) GetReadings() (float32, float32) {
	return 0, 0
}
