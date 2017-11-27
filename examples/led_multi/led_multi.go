package main

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"time"

	"github.com/joeblubaugh/adafruit-led/HT16K33"
)

func main() {
	invalidCommandExit := func() {
		fmt.Printf("Usage: %s {addr}:{bus}[,{addr}:{bus}...]\n", os.Args[0])
		fmt.Printf("eg. %s 0x70:1,0x72:1\n", os.Args[0])
		os.Exit(1)
	}

	if len(os.Args) != 2 {
		invalidCommandExit()
	}

	devices, err := HT16K33.ParseDevices(os.Args[1])

	if err != nil {
		invalidCommandExit()
	}

	for _, device := range devices {
		go runTest(device)
	}

	// Loop forever until process is killed (eg. via CTRL-C).
	// Otherwise we will exit the process right away since
	// we're doing all the LED control on goroutines
	for {
		runtime.Gosched()
	}
}

func runTest(device *HT16K33.EightByEight) {
	for {
		if err := device.Clear(); err != nil {
			log.Panicf("led_multi: %v\n", err)
		}

		time.Sleep(500 * time.Millisecond)

		for x := byte(0); x < 8; x++ {
			for y := byte(0); y < 8; y++ {
				if err := device.SetPixel(x, y, true); err != nil {
					log.Panicf("led_multi: %v\n", err)
				}

				time.Sleep(50 * time.Millisecond)
			}
		}

		time.Sleep(500 * time.Millisecond)
	}
}
