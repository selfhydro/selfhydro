package sensors

import (
	"errors"
	"gotest.tools/assert"
	"os"
	"testing"
)

func TestGetID(t *testing.T) {
	tempSensor := DS18b20{}
	dataDirectory = "./testdata/sensor_test/"
	tempSensor.SetupDevice()
	assert.Equal(t, "testSensor", tempSensor.Id)
}

func TestShouldLogErrorIfCantGetId(t *testing.T) {
	dataDirectory = ""
	getReadDir = mockReadDir
	tempSensor := DS18b20{}
	assert.ErrorContains(t, tempSensor.SetupDevice(), "error finding directory for sensor")
}

func mockReadDir(dir string) ([]os.FileInfo, error) {
	if dir == "" {
		return nil, errors.New("error")
	}
	return nil, nil
}

func TestShouldReadTemp(t *testing.T) {
	waterTempSensor := DS18b20{
		Id: "testSensor",
	}
	dataDirectory = "testdata/sensor_test"
	temp, _ := waterTempSensor.GetState()
	assert.Equal(t, float32(10), temp)
}

func TestShouldReturnErrorWhenFailsToReadTemp(t *testing.T) {
	waterTempSensor := DS18b20{
		Id: "fakeSensor",
	}
	_, err := waterTempSensor.GetState()
	assert.ErrorContains(t, err, "failed to read sensor temperature")
}
