package main

type selfhydro struct{
  currentTemp float32
}

func (sh selfhydro) GetAmbientTemp() float32 {

  return 10
}
