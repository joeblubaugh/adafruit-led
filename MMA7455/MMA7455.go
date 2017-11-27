package MMA7455

import (
	"github.com/joeblubaugh/adafruit-led/i2c"
)

const (
	ADDR          = 0x3A
	REGISTER_XOUT = 0x06
	REGISTER_YOUT = 0x07
	REGISTER_ZOUT = 0x08
)

type Device struct {
	bus *i2c.I2CBus
}

func New() (bp *Device, err error) {
	bp = new(Device)
	bp.bus, err = i2c.Bus(0)
	return
}

func (bp *Device) Read(reg byte) (value int8, err error) {
	var bytes []byte

	bytes, err = bp.bus.ReadByteBlock(ADDR, reg, 1)

	if err != nil {
		return
	}

	value = int8(bytes[0])

	return
}
