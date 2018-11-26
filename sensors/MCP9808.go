package sensors

import "log"

type MCP9808 struct {
	address string
	device i2cDevice
}

const (
	defaultI2CAddr = 0x18
	regTemp             = 0x05
)

func NewMCP9808() Sensor {
	return &MCP9808{}
}

func (mcp9808 *MCP9808) SetupDevice() (error) {
	var err error
	mcp9808.device, err = NewI2C(defaultI2CAddr, 1)
	if err != nil {
		log.Print(err)
		return err
	}
	return nil
}

func (mcp9808 MCP9808) GetState() (float32, error) {
	encodedTemp, err := mcp9808.device.ReadRegU16BE(regTemp)
	if err != nil {
		log.Printf("error reading temp from sensor: %s", err)
		return 0, err
	}
	temp, err := mcp9808.getTemp(encodedTemp)
	return temp, nil
}

func (mcp9808 MCP9808) getTemp(encodedTemp uint16) (float32, error) {
	temp := float32(encodedTemp&0x0fff) / 16
	return temp, nil
}
