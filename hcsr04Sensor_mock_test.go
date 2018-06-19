package main

type mockUltrasonicSensor struct {

}

func (us *mockUltrasonicSensor) MeasureDistance() (cm float32) {
	return 0
}