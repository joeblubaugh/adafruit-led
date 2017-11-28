package HT16K33

import (
	"fmt"

	"github.com/joeblubaugh/adafruit-led/i2c"
)

const (
	REGISTER_DISPLAY_SETUP = 0x80
	REGISTER_SYSTEM_SETUP  = 0x20
	REGISTER_DIMMING       = 0xE0
	BLINKRATE_OFF          = 0x00
	BLINKRATE_2HZ          = 0x01
	BLINKRATE_1HZ          = 0x02
	BLINKRATE_HALFHZ       = 0x03
)

// Device is the General HT16K33 controller interface
type Device struct {
	ImmediateUpdate bool
	bus             *i2c.I2CBus
	busNum          byte
	addr            byte

	// buffer mirrors the internal device Display RAM. See:
	// https://cdn-shop.adafruit.com/datasheets/ht16K33v110.pdf
	//
	// Each entry in the buffer is a complete "Row" in the
	// Device, for ROW0 - ROW15. However, each column must be
	// written on the I2C interface as two bytes.
	buffer [8]uint16
}

func (bp *Device) Init(addr, busNum byte) (err error) {
	if bp.bus, err = i2c.Bus(busNum); err != nil {
		return
	}

	bp.ImmediateUpdate = true
	bp.busNum = busNum
	bp.addr = addr

	if err = bp.bus.WriteByte(addr, REGISTER_SYSTEM_SETUP|0x01, 0x00); err != nil {
		return
	}

	bp.SetBlinkRate(BLINKRATE_OFF)
	bp.SetBrightness(15)
	bp.ReadDisplay()

	return
}

func (bp *Device) Addr() byte {
	return bp.addr
}

func (bp *Device) BusNum() byte {
	return bp.busNum
}

func (bp *Device) SetBrightness(brightness byte) (err error) {
	if brightness < 0 || brightness > 15 {
		err = fmt.Errorf("i2c: Invalid brightness: %v\n", brightness)
		return
	}

	err = bp.bus.WriteByte(bp.addr, REGISTER_DIMMING|brightness, 0x00)

	return
}

func (bp *Device) SetBlinkRate(blinkRate byte) (err error) {
	if blinkRate < BLINKRATE_OFF || blinkRate > BLINKRATE_HALFHZ {
		err = fmt.Errorf("i2c: Invalid blinkRate: %v\n", blinkRate)
		return
	}

	err = bp.bus.WriteByte(bp.addr, REGISTER_DISPLAY_SETUP|0x01|(blinkRate<<1), 0x00)

	return
}

func (bp *Device) SetBufferRow(row byte, word uint16) (err error) {
	if row < 0 || row > 7 {
		err = fmt.Errorf("i2c: Invalid row: %v\n", row)
		return
	}

	bp.buffer[row] = word

	if bp.ImmediateUpdate {
		err = bp.WriteDisplay()
	}

	return
}

func (bp *Device) Clear() (err error) {
	for i := range bp.buffer {
		bp.buffer[i] = 0x00
	}

	if bp.ImmediateUpdate {
		err = bp.WriteDisplay()
	}

	return
}

func (bp *Device) ReadDisplay() (err error) {
	var bytes []byte

	if bytes, err = bp.bus.ReadByteBlock(bp.addr, 0x00, 16); err != nil {
		return
	}

	j := 0

	for i, _ := range bp.buffer {
		bp.buffer[i] = uint16((bytes[j] & 0xff) | (bytes[j+1] >> 8))
		j += 2
	}

	return
}

func (bp *Device) WriteDisplay() (err error) {
	bytes := make([]byte, 16)

	i := 0

	for _, item := range bp.buffer {
		bytes[i], bytes[i+1] = byte(item&0xff), byte(item>>8&0xff)
		i += 2
	}

	err = bp.bus.WriteByteBlock(bp.addr, 0x00, bytes)

	return
}
