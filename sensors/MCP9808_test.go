package sensors

import (
	"testing"
	"errors"
	"gotest.tools/assert"
)

func TestShouldReturn0TempIfCantGetReading(t *testing.T) {
	mockI2C := &mockI2cDevice{}
	err := errors.New("test error")
	mockI2C.On("ReadRegU16BE", uint8(regTemp)).Return(uint16(0x0), err).Once()
	mcp9808 := MCP9808 {
		device: mockI2C,
	}
	temp, _ := mcp9808.GetState()
	assert.Equal(t, temp, float32(0))
}

func TestShouldGet25point25CelciusFrom0x0020(t *testing.T) {
	mockI2C := &mockI2cDevice{}
	mockI2C.On("ReadRegU16BE", uint8(regTemp)).Return(uint16(0x0194),nil).Once()
	mcp9808 := MCP9808 {
		device: mockI2C,
	}
	temp, _ := mcp9808.GetState()
	assert.Equal(t, temp, float32(25.25))
}

func TestReturnNegativeTempWhenBelowZero(t *testing.T) {
	mockI2C := &mockI2cDevice{}
	mockI2C.On("ReadRegU16BE", uint8(regTemp)).Return(uint16(0x1194),nil).Once()
	mcp9808 := MCP9808 {
		device: mockI2C,
	}
	temp, _ := mcp9808.GetState()
	assert.Equal(t, temp, float32(-25.25))
}
