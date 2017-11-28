package HT16K33

import "fmt"

const (
	barMin uint = 0
	barMax uint = 23

	barOff    = 0x0
	barGreen  = 0x1
	barRed    = 0x2
	barYellow = 0x3
)

// TwentyFourBar models a 24-bar LED Bar Graph backpack.
type TwentyFourBar struct {
	Device
}

func NewTwentyFourBar(addr, bus byte) (t *TwentyFourBar, err error) {
	t = &TwentyFourBar{}
	err = t.Init(addr, bus)
	return
}

// Set bar sets the value of a bar (0-23) on the graph to the
// specified color. Color must be one of barOff, barGreen, barRed,
// barYellow.
func (t *TwentyFourBar) SetBar(bar, color byte) error {
	// Each Bar on the graph is mapped as two LEDs on the HT16K33
	if bar > 23 {
		return fmt.Errorf("Invalid Bar #: %v", bar)
	}

	// The row register is the same for both colors.
	var row byte
	if bar < 12 {
		row = bar / 4 // integer floor
	} else {
		row = (bar - 12) / 4
	}

	var offset uint16 = uint16(bar % 4)
	// The upper bars are shifted over 4 bits.
	if bar >= 12 {
		offset = offset + 4
	}

	var newBuffer uint16 = t.buffer[row]
	if color&barRed > 0 {
		newBuffer = newBuffer | 1<<offset
	} else {
		newBuffer = newBuffer & ^(1 << offset)
	}

	var offsetGreen uint16 = 8
	if color&barGreen > 0 {
		newBuffer = newBuffer | 1<<(offset+offsetGreen)
	} else {
		newBuffer = newBuffer & ^(1 << (offset + offsetGreen))
	}

	return t.SetBufferRow(row, newBuffer)
}
