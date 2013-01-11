package main

import (
	"bitbucket.org/gmcbay/i2c/MMA7455"
	"log"
	"math"
)

// Used as a stupidly coarse jitter filter, if the value on an
// axis doesn't change by this amount, ignore it
const REPORT_THRESHOLD = 12

type accelValue struct {
	reg   byte
	value int8
}

func main() {
	accel, err := MMA7455.New()

	if err != nil {
		log.Panicf("accel_dump: %v\n", err)
	}

	valueChan := make(chan *accelValue)

	go readValues(valueChan, accel, MMA7455.REGISTER_XOUT)
	go readValues(valueChan, accel, MMA7455.REGISTER_YOUT)
	go readValues(valueChan, accel, MMA7455.REGISTER_ZOUT)

	// [x, y, z] tracking variables
	values := make([]int8, 3)

	for {
		select {
		case v := <-valueChan:
			// This relies on the assumption that REGISTER_XOUT is the lowest
			// value and the subsequent values are contiguous
			valueIndex := v.reg - MMA7455.REGISTER_XOUT

			if math.Abs(float64(v.value-values[valueIndex])) > REPORT_THRESHOLD {
				log.Printf("accel_dump: (X: %v, Y: %v, Z: %v)\n", values[0], values[1], values[2])
				values[valueIndex] = v.value
			}
		}
	}
}

func readValues(valueChan chan *accelValue, accel *MMA7455.Device, reg byte) {
	for {
		if value, err := accel.Read(reg); err != nil {
			log.Panicf("accel_dump: %v\n", err)
		} else {
			valueChan <- &accelValue{reg, value}
		}
	}
}
