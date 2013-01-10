package main

import (
	"bitbucket.org/gmcbay/i2c/HT16K33"
	"log"
	"time"
)

func main() {
	// Assumes a device at address 0x70 on I2C bus 1
	// Modify if your device uses a different address or bus number
	grid, err := HT16K33.NewEightByEight(0x70, 1)

	if err != nil {
		log.Panicf("led_simple: %v\n", err)
	}

	for {
		if err := grid.Clear(); err != nil {
			log.Panicf("led_simple: %v\n", err)
		}

		time.Sleep(500 * time.Millisecond)

		for x := byte(0); x < 8; x++ {
			for y := byte(0); y < 8; y++ {
				if err := grid.SetPixel(x, y, true); err != nil {
					log.Panicf("led_simple: %v\n", err)
				}
				time.Sleep(50 * time.Millisecond)
			}
		}

		time.Sleep(500 * time.Millisecond)
	}
}
