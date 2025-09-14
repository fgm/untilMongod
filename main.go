/*
Package main provides the untilMongod command, a safe way to wait
until a mongod (or mongos) instance is ready to accept connections,
without undue delay.

Usage

	untilMongod [flags]

The following optional flags are available:

	-url string
		The MongoDB URL to which to connect.
		(default "mongodb://localhost:270117")
	-timeout uint
		The maximum delay in seconds after which the command will stop trying to connect.
		(default 30)
	-v	Make the command more verbose (default: false)

Example

	untilMongod -url mongodb://example.com:11117?directConnection=true -timeout 60

Return codes

	0: Success
	1: Timeout
	2: Other errors
*/
package main

import (
	"flag"
	"os"
	"time"

	"github.com/fgm/untilMongod/dial"
)

func parseFlags(url *string, timeout *time.Duration, verbose *bool) {
	flag.StringVar(url, "url", "mongodb://localhost:27017",
		"The mongodb URL to which to connect.")
	flag.BoolVar(verbose, "v", false,
		"Make the command more verbose")

	uintTimeout := flag.Uint("timeout", 30,
		"The maximum delay in seconds after which the command will stop trying to connect")

	flag.Parse()

	*timeout = time.Duration(*uintTimeout) * time.Second
}

func main() {
	var (
		dialer  = dial.NewMongoDbDial()
		timeout time.Duration
		url     string
		verbose bool
	)

	parseFlags(&url, &timeout, &verbose)

	reporter := dial.NewReporter(verbose)
	dialResult := dial.Dial(url, timeout, dialer, reporter)
	os.Exit(int(dialResult))
}
