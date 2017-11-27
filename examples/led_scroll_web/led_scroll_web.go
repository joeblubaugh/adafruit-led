package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/joeblubaugh/adafruit-led/HT16K33"
)

// markit on demand json results look like this:
// {"Data":{"Status":"SUCCESS","Name":"Microsoft Corp","Symbol":"MSFT","LastPrice":26.56,"Change":0.00999999999999801,
//   "ChangePercent":0.0376647834274824,"Timestamp":"Fri Dec 28 15:59:59 UTC-05:00 2012","MarketCap":223541230720,
//   "Volume":3136569,"ChangeYTD":25.96,"ChangePercentYTD":2.31124807395993,"High":26.89,"Low":26.56,"Open":26.72}}
//
// So we create structs which match this format using the same field names, allowing for a simple json.Unmarshal call
// to parse the data for us
//
// Actual quote data
type quote struct {
	Symbol    string
	LastPrice float32
	Change    float32
}

// Wrapper struct to deal with the fact that the JSON responses are wrapped in a 'Data' response object
type dataWrapper struct {
	Data quote
}

func main() {
	invalidCommandExit := func() {
		fmt.Printf("Usage: %s {addr}:{bus}[,{addr}:{bus}...] [ticker symbol,ticker symbol...]\n", os.Args[0])
		fmt.Printf("eg. %s 0x70:1,0x71:1,0x72:1 msft,goog,amzn\n", os.Args[0])
		os.Exit(1)
	}

	if len(os.Args) < 3 {
		invalidCommandExit()
	}

	devices, err := HT16K33.ParseDevices(os.Args[1])

	if err != nil {
		invalidCommandExit()
	}

	padding := ""

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
		padding = fmt.Sprintf("  %s", padding)
	}

	symbols := os.Args[2]

	dataChan := make(chan *quote)

	for _, symbol := range strings.Split(symbols, ",") {
		go lookupLoop(symbol, dataChan)
	}

	quotes := make(map[string]*quote, 0)

	// Loop forever
	for {
		// See if any of the symbols have new quote data, if they
		// do it'll arrive on the dataChan. But only wait up until
		// 50 milliseconds for new data to avoid blocking the scroller
		select {
		case quote := <-dataChan:
			quotes[quote.Symbol] = quote

		case <-time.After(50 * time.Millisecond):
		}

		// Dump current contents of the quotes map into a string to be scrolled
		var buffer bytes.Buffer
		buffer.WriteString(padding)

		for _, quote := range quotes {
			buffer.WriteString(fmt.Sprintf("%v: %0.3f (CHG: %0.3f)   ", quote.Symbol, quote.LastPrice, quote.Change))
		}

		HT16K33.ScrollMessage(buffer.String(), devices, 210)
	}
}

func lookupLoop(symbol string, dataChan chan *quote) {
	for {
		response, err := http.Get(
			fmt.Sprintf("http://dev.markitondemand.com/Api/Quote/json?symbol=%s", symbol))

		if err != nil {
			log.Panicf("led_scroll_web: %v\n", err)
		}

		responseBytes, err := ioutil.ReadAll(response.Body)
		response.Body.Close()

		if err != nil {
			log.Panicf("led_scroll_web: %v\n", err)
		}

		// Allocate a new 'data' object each time because once
		// we send the child quote pointer across the dataChan
		// it is 'owned' (via golang channel sending convention)
		// by the main goroutine.  We don't want to write to the same
		// memory area next time through
		data := new(dataWrapper)

		if err = json.Unmarshal(responseBytes, data); err != nil {
			log.Panicf("led_scroll_web: %v\n", err)
		}

		// Write result to the dataChan channel so main loop
		// can pick up the changes
		dataChan <- &data.Data

		// Wait between 1 and 2 minutes to refresh data to avoid
		// slamming the API with requests. Randomized to increase the
		// spread of hits from all the different goroutines that may
		// be running lookups
		time.Sleep(time.Duration(60+rand.Int31n(60)) * time.Second)
	}
}
