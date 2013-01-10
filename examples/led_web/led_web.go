package main

import (
	"bitbucket.org/gmcbay/i2c/HT16K33"
	"fmt"
	"log"
	"net/http"
	"os"
)

var devices []*HT16K33.EightByEight

func main() {
	invalidCommandExit := func() {
		fmt.Printf("Usage: %s {addr}:{bus}[,{addr}:{bus}...]\n", os.Args[0])
		fmt.Printf("eg. %s 0x70:1,0x72:1\n", os.Args[0])
		os.Exit(1)
	}

	var err error

	if len(os.Args) != 2 {
		invalidCommandExit()
	}

	devices, err = HT16K33.ParseDevices(os.Args[1])

	if err != nil {
		invalidCommandExit()
	}

	for _, device := range devices {
		device.ImmediateUpdate = false
	}

	http.HandleFunc("/led", rootHandler)

	// Use port 80, we already need sudo access to write to I2C.
	// Change this to something else if you're running a real web
	// server on port 80 already
	log.Fatal(http.ListenAndServe(":80", nil))
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		log.Panicf("led_web: %v\n", err)
	}

	fmt.Fprintf(w, "<html><body><form method='POST'><input type='hidden' name='formSubmitted' value='true'>")

	printState := func(state bool) string {
		if state {
			return "checked"
		}

		return ""
	}

	var err error

	formSubmitted := len(r.FormValue("formSubmitted")) > 0

	for i, device := range devices {
		fmt.Fprintf(w, "<div style='float: left; padding: 1em'>")
		fmt.Fprintf(w, "<div>0x%x:%v</div>", device.Addr(), device.BusNum())
		for y := byte(0); y < 8; y++ {
			fmt.Fprintf(w, "<div>")

			for x := byte(0); x < 8; x++ {
				var state bool

				// If this page view is the result of the form being submitted, take the
				// state from the web form checkbox value, otherwise read it off the device
				if formSubmitted {
					state = r.FormValue(fmt.Sprintf("g_%d_%d", i, y*8+x)) == "on"
				} else {
					state = device.Pixel(x, y)
				}

				if err = device.SetPixel(x, y, state); err != nil {
					log.Panicf("led_web: %v\n", err)
				}

				fmt.Fprintf(w, "<input type='checkbox' name='g_%d_%d' %s onClick='this.form.submit();'>",
					i, y*8+x, printState(state))
			}

			fmt.Fprintf(w, "</div>")
		}

		fmt.Fprintf(w, "</div>")

		device.WriteDisplay()
	}

	fmt.Fprintf(w, "</form></body></html>")
}
