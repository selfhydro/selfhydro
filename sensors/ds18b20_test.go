package sensors

import (
	"errors"
	"log"
	"os"
	"testing"
  "gotest.tools/assert"
)

func TestGetID(t *testing.T) {
	tempSensor := ds18b20{}
	dataDirectory = "./testdata/sensor_test/"
	tempSensor.SetupDevice()
  assert.Equal(t, "testSensor", tempSensor.id)
}

func TestShouldLogErrorIfCantGetId(t *testing.T) {
	dataDirectory = ""
	getReadDir = mockReadDir
	tempSensor := ds18b20{}
  assert.ErrorContains(t, 	tempSensor.SetupDevice(), "error finding directory for sensor")
}

func mockReadDir(dir string) ([]os.FileInfo, error) {
  log.Print("using mock")
	if dir == "" {
		return nil, errors.New("error")
	}
	return nil, nil
}

func TestShoudlReadTemp(t *testing.T) {
	waterTempSensor := ds18b20 {
    id: "testSensor",
  }
	dataDirectory = "testdata/sensor_test"
	temp := waterTempSensor.GetState()
  assert.Equal(t, float64(10), temp)
}
