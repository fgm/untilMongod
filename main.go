/*
The untilMongod command provides a safe way to wait until a mongod (or mongos)
instance is ready to accept connections, without undue delay.

Syntax

Here is how to invoke the command:
  untilMongo -url mongodb://example.com:11117 -timeout 60 -v

  • -url is a typical MongoDB URL
  • -timeout is the maximum delay (in seconds) the command will wait before aborting
  • -v enable verbose mode

untilMongod returns 0 on success, 1 on timeout, 2 on other errors.
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
	var url string
	var timeout time.Duration
	var verbose bool
	var dialer = dial.NewMongoDbDial()

	parseFlags(&url, &timeout, &verbose)

	reporter := dial.NewReporter(verbose)
	dialResult := dial.Dial(url, timeout, dialer, reporter)
	os.Exit(int(dialResult))
}
