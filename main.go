/*
The untilMongod command provides a safe way to wait until a mongod (or mongos)
instance is ready to accept connections, without undue delay.
*/
package main

import (
	"flag"
	"fmt"
	"time"

	"os"

	"gopkg.in/mgo.v2"
)

type DialResult int

const (
	Success    DialResult = 0
	Timeout               = 1
	OtherError            = 2
)

// Provide an increasing delay from an existing one.
func backOff(current time.Duration) time.Duration {
	const BackoffFactor = 1.5

	floatNext := float64(current) * BackoffFactor
	return time.Duration(floatNext)
}

// Dial attempts to connect to the specified server for the specified duration.
//
// It returns a status code usable with os.Exit().
func Dial(url string, maxTimeout time.Duration, verbose bool) DialResult {
	const ExpectedErrorString = "no reachable servers"
	const Nanoseconds = float64(time.Second)

	t0 := time.Now()
	tMax := t0.Add(maxTimeout)

	timeout := 1 * time.Millisecond

	for {
		_, err := mgo.DialWithTimeout(url, timeout)

		// Success.
		if err == nil {
			return Success
		}

		// Timeout reached.
		t1 := time.Now()
		if t1.After(tMax) {
			return Timeout
		}

		if err.Error() != ExpectedErrorString {
			return OtherError
		}

		timeout = backOff(timeout)
		if verbose {
			fmt.Printf("Unavailable in %3f seconds, retrying in %.3f seconds.\n",
				float64(t1.Sub(t0))/Nanoseconds,
				float64(timeout)/Nanoseconds)
		}
		time.Sleep(timeout)
	}

}

func parseFlags(url *string, timeout *time.Duration, verbose *bool) {
	flag.StringVar(url, "url", "mongodb://localhost:27017",
		"The mongodb URL to which to connect.")
	flag.BoolVar(verbose, "v", true, "Make the command more verbose")

	uintTimeout := flag.Uint("timeout", 30,
		"The maximum delay in seconds after which the command will stop trying to connect")

	flag.Parse()

	*timeout = time.Duration(*uintTimeout) * time.Second
}

func main() {
	var url string
	var timeout time.Duration
	var verbose bool

	parseFlags(&url, &timeout, &verbose)

	os.Exit(int(Dial(url, timeout, verbose)))
}
