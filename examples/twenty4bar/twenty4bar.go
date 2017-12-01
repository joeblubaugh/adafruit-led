package main

import (
	"log"
	"time"

	"github.com/joeblubaugh/adafruit-led/HT16K33"
)

func main() {
	// Assumes a device at address 0x70 of I2C bus 1
	// Modify if your device is on a different bus or address.
	bar, err := HT16K33.NewTwentyFourBar(0x73, 1)
	if err != nil {
		log.Panic("twenty4bar: %v\n", err)
	}

	for {
		if err := bar.Clear(); err != nil {
			log.Panic("twenty4bar: %v\n", err)
		}
		time.Sleep(500 * time.Millisecond)

		for _, c := range []byte{HT16K33.BarRed, HT16K33.BarYellow, HT16K33.BarGreen} {
			for b := byte(0); b <= 23; b++ {
				if err := bar.SetBar(b, c); err != nil {
					log.Panic("twenty4bar: %v\n", err)
				}
				time.Sleep(50 * time.Millisecond)
			}
			time.Sleep(500 * time.Millisecond)
		}
	}
}
