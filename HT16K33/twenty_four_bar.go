package HT16K33

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
	t := &TwentyFourBar{}
	err = t.Init(addr, bus)
	return
}

func (t *TwentyFourBar) SetBar(bar int, color byte) error {
	// TODO: Figure out mapping of Device.buffer rows to bytes for this device.
	return nil
}
