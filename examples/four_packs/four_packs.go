package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joeblubaugh/adafruit-led/HT16K33"
)

func main() {
	done := make(chan struct{})

	// Turn on all four items and run patterns on them in goroutines.
	grid0, err := HT16K33.NewEightByEight(0x70, 1)
	if err != nil {
		log.Panic("eight_by_eight: %v\n", err)
	}
	grid1, err := HT16K33.NewEightByEight(0x71, 1)
	if err != nil {
		log.Panic("eight_by_eight: %v\n", err)
	}

	go animateGrid(grid0, done)
	go animateGrid(grid1, done)

	bar0, err := HT16K33.NewTwentyFourBar(0x72, 1)
	if err != nil {
		log.Panic("twenty4bar: %v\n", err)
	}
	bar1, err := HT16K33.NewTwentyFourBar(0x73, 1)
	if err != nil {
		log.Panic("twenty4bar: %v\n", err)
	}

	go animateBar(bar0, done)
	go animateBar(bar1, done)

	// Wait for shutdown signal, then clean up.
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	go func(s chan os.Signal, stop chan struct{}) {
		sig := <-s
		log.Println("Sig: ", sig)
		close(stop)
	}(sigs, done)

	for {
		_, ok := <-done
		if !ok {
			time.Sleep(8 * time.Second)
			return
		}
	}
}

func animateGrid(grid *HT16K33.EightByEight, done chan struct{}) {
	for {
		select {
		case _, ok := <-done:
			if !ok {
				grid.Clear()
				return
			}
		default:
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
}

func animateBar(bar *HT16K33.TwentyFourBar, done chan struct{}) {
	for {
		select {
		case _, ok := <-done:
			if !ok {
				bar.Clear()
				return
			}
		default:
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
}
