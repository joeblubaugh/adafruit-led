package HT16K33

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

// Control 8x8 LED

type EightByEight struct {
	Device
}

func NewEightByEight(addr, bus byte) (e *EightByEight, err error) {
	e = new(EightByEight)
	err = e.Init(addr, bus)
	return
}

func (e *EightByEight) Pixel(x, y byte) bool {
	return ((e.buffer[y]) & (1 << ((x + 7) % 8))) != 0
}

func (e *EightByEight) SetPixel(x, y byte, on bool) (err error) {
	if x < 0 || x > 7 || y < 0 || y > 7 {
		err = fmt.Errorf("i2c: Invalid pixel location: %v, %v\n", x, y)
		return
	}

	x = (x + 7) % 8

	if on {
		err = e.SetBufferRow(y, e.buffer[y]|1<<x)
	} else {
		err = e.SetBufferRow(y, e.buffer[y] & ^(1<<x))
	}

	return err
}

// Helper func to generate a list of devices based on an input string
// matching the spec addr:bus,addr:bus...
func ParseDevices(deviceList string) (devices []*EightByEight, err error) {
	splitDeviceList := strings.Split(deviceList, ",")

	devices = make([]*EightByEight, len(splitDeviceList))

	for i, deviceSpec := range splitDeviceList {
		deviceSpecParts := strings.Split(deviceSpec, ":")

		if len(deviceSpecParts) != 2 {
			err = errors.New(fmt.Sprintf("Invalid device list: %s", deviceList))
			return
		}

		var addrInt, busInt int64

		if addrInt, err = strconv.ParseInt(deviceSpecParts[0], 0, 8); err != nil {
			return
		} else {
			if busInt, err = strconv.ParseInt(deviceSpecParts[1], 0, 8); err != nil {
				return
			} else {
				if devices[i], err = NewEightByEight(byte(addrInt), byte(busInt)); err != nil {
					return
				}
			}
		}
	}

	return
}

func ScrollMessage(msg string, devices []*EightByEight, speed byte) {
	// Pre-calculate how may total horizontal pixels our scroller needs
	msgLen := len(msg)
	horizPixelCount := msgLen * FONT_WIDTH
	displayWidth := 8 * len(devices)

	for i := 0; i < horizPixelCount; i++ {
		// each 'frame', loop through the currently visible x/y
		// values and map their character on/off 'bits' to LED matrix on/off bits

		for x := 0; x < displayWidth; x++ {

			xOffset := i + x
			charIndex := xOffset / FONT_WIDTH
			deviceNum, deviceX := int(x/8), byte(x%8)

			if charIndex < msgLen && charIndex > -1 {
				charOffset := xOffset % FONT_WIDTH
				charBlock := chars[msg[charIndex]]

				// For characters not defined in the font, use space
				if charBlock == nil {
					charBlock = chars[' ']
				}

				for y := 0; y < 7; y++ {
					on := charOffset < (FONT_WIDTH-1) && charBlock[y][charOffset] == '*'

					if err := devices[deviceNum].SetPixel(deviceX, byte(y), on); err != nil {
						log.Panicf("i2c_test_scroll: %v\n", err)
					}
				}
			}
		}

		// flush device buffers out to their respective devices
		for _, device := range devices {
			if err := device.WriteDisplay(); err != nil {
				log.Panicf("i2c_test_scroll: %v\n", err)
			}
		}

		// Sleep between 'frames' to vary the speed of the scroller
		time.Sleep(time.Duration(255-speed) * time.Millisecond)
	}
}
