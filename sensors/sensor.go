package sensors

type Sensor interface {
	SetupDevice() (error)
	GetState() (float32, error)
}
