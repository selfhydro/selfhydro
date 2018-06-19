package main

type mockAmbientTemp struct {

}

func (sensor mockAmbientTemp) GetTemp() float32 {
	return 0
}
