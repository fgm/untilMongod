package dial

import (
	"fmt"
	"time"

	"github.com/globalsign/mgo"
	"io"
	"os"
)

// DialResult is the type of the Dial() results.
type DialResult int

const (
	Success DialResult = iota
	Timeout
	OtherError
)

// DriverDial is the type of the lower-level functions actually implementing
// the "connect to a server" feature, without the backoff mechanism.
//
// Its implementations are typically thin wrappers around a driver-provided
// mechanism.
type DriverDial = func(url string, duration time.Duration) error

// NewMongoV2Dial returns a new DriverDial function using the MGOv2 driver
// available at https://github.com/globalsign/mgo
//
// That driver replaces the original Labix driver and the now-unmaintained
// https://github.com/go-mgo/mgo/tree/v2 driver.
func NewMgoV2Dial() DriverDial {
	return func(url string, duration time.Duration) error {
		_, err := mgo.DialWithTimeout(url, duration)
		return err
	}
}

// NewMongoDbDial returns a new DriverDial function using the official MongoDB
// driver available at https://github.com/mongodb/mongo-go-driver
func NewMongoDbDial() DriverDial {
	// BUG(FGM): NewMongoDBDial() needs to be actually implemented, when that driver
	// reaches at least beta status and gets some documentation.
	return func(url string, duration time.Duration) error {
		panic("MongoDB official driver is not implemented yet")
	}
}

// The string used by MongoDB to signal it could not connect.
const ExpectedErrorString = "no reachable servers"

// Dial attempts to connect to the specified server for the specified duration,
// using the specified connection method.
//
// It returns a status code usable with os.Exit().
//
// The final variadic io.Writer argument allows passing ONE specific writer,
// and defaulting to os.Stderr when none is passed.
func Dial(url string, maxTimeout time.Duration, verbose bool, dialer DriverDial, w ...io.Writer) DialResult {
	var writer io.Writer
	if len(w) == 0 {
		writer = os.Stderr
	} else {
		writer = w[0]
	}
	const Nanoseconds = float64(time.Second)

	t0 := time.Now()
	tMax := t0.Add(maxTimeout)

	timeout := 1 * time.Millisecond

	for {
		err := dialer(url, timeout)

		// Success.
		if err == nil {
			return Success
		}

		// Unexpected errors.
		if err.Error() != ExpectedErrorString {
			return OtherError
		}

		// Timeout reached.
		t1 := time.Now()
		if t1.After(tMax) {
			return Timeout
		}

		timeout = backOff(timeout)
		if verbose {
			fmt.Fprintf(writer, "Unavailable in %3f seconds, retrying in %.3f seconds.\n",
				float64(t1.Sub(t0))/Nanoseconds,
				float64(timeout)/Nanoseconds)
		}
		time.Sleep(timeout)
	}

}
