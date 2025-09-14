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
	"io"
	"os"
	"time"

	"github.com/fgm/untilMongod/dial"
)

func parseFlags(fs *flag.FlagSet, args []string, url *string, timeout *time.Duration, verbose *bool) {
	fs.StringVar(url, "url", "mongodb://localhost:27017",
		"The mongodb URL to which to connect.")
	fs.BoolVar(verbose, "v", false,
		"Make the command more verbose")

	uintTimeout := fs.Uint("timeout", 30,
		"The maximum delay in seconds after which the command will stop trying to connect")

	_ = fs.Parse(args)

	*timeout = time.Duration(*uintTimeout) * time.Second
}

func realMain(w io.Writer, fs *flag.FlagSet, args []string, dialer dial.DriverDial) int {
	var (
		timeout time.Duration
		url     string
		verbose bool
	)

	parseFlags(fs, args, &url, &timeout, &verbose)
	reporter := dial.NewReporter(verbose, w)
	dialResult := dial.Dial(url, timeout, dialer, reporter)
	return int(dialResult)
}

func main() {
	os.Exit(realMain(os.Stderr, flag.CommandLine, os.Args[1:], dial.NewMongoDbDial()))
}
