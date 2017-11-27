package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/joeblubaugh/adafruit-led/HT16K33"
)

func main() {
	invalidCommandExit := func() {
		fmt.Printf("Usage: %s {addr}:{bus}[,{addr}:{bus}...] {message}\n", os.Args[0])
		fmt.Printf("eg. %s 0x70:1,0x71:1,0x72:1 hello, world!\n", os.Args[0])
		os.Exit(1)
	}

	if len(os.Args) < 3 {
		invalidCommandExit()
	}

	devices, err := HT16K33.ParseDevices(os.Args[1])

	if err != nil {
		invalidCommandExit()
	}

	msg := strings.Join(os.Args[2:], " ")

	for _, device := range devices {
		// Make sure device buffer is cleared of any pixels
		// currently in the on state.
		device.Clear()

		// Don't flush data out to I2C device on every pixel write,
		// instead flush it all out at once manually when a 'frame' is ready
		device.ImmediateUpdate = false

		// Pad message with an appropriate number of spaces such that it initally
		// scrolls in from the right instead of just popping on to the LED(s).  Two
		// spaces per device works well.
		msg = fmt.Sprintf("  %s", msg)
	}

	// Loop forever
	for {
		HT16K33.ScrollMessage(msg, devices, 205)
	}
}
