package main

type selfhydro struct {
	currentTemp float32
	waterLevel  WaterLevel
}

func (sh selfhydro) GetAmbientTemp() float32 {

	return 10
}

func (sh selfhydro) GetWaterLevel() float32 {

	return sh.waterLevel.waterLevel
}
