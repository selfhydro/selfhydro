package sensors

import (
	i2c "github.com/d2r2/go-i2c"
)

type i2cDevice interface {
	ReadRegU16BE(reg byte) (uint16, error)
}

func NewI2C(addr uint8, bus int) (*i2c.I2C, error) {
	return i2c.NewI2C(addr, bus)
}
