package dial

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

const (
	// DefaultURL is the default MongoDB server URL.
	DefaultURL = "mongodb://localhost:27017"
	// DefaultBaseTimeout is the initial timeout used by Dial() before backoff starts.
	DefaultBaseTimeout = 1 * time.Millisecond
)

// Result is the type of the Dial() results.
type Result int

// These constants represent the result of Dial().
const (
	Success Result = iota
	Timeout
	OtherError
)

// DriverDial is the type of the lower-level functions actually implementing
// the "connect to a server" feature, without the backoff mechanism.
//
// Its implementations are typically thin wrappers around a driver-provided
// mechanism.
type DriverDial = func(url string, duration time.Duration) error

// NewMongoDbDial returns a new DriverDial function using the official MongoDB
// driver available at https://github.com/mongodb/mongo-go-driver
func NewMongoDbDial() DriverDial {
	return func(url string, duration time.Duration) error {
		ctx, cancel := context.WithTimeout(context.Background(), duration)
		c, err := mongo.Connect(ctx, options.Client().ApplyURI(url))
		if err != nil {
			cancel()
			return err
		}
		err = c.Ping(ctx, readpref.Primary())
		if err != nil {
			cancel()
			return err
		}
		err = c.Disconnect(ctx)
		cancel()
		return err
	}
}

// ExpectedErrorString is used by MongoDB clients to signal it could not
// connect, starting with the mongo shell.
const ExpectedErrorString = "no reachable servers"

// Reporter prints to a writer, defaulting to stderr, only if verbose is true.
type Reporter struct {
	writer  io.Writer
	verbose bool
}

// Printf only prints if r.verbose is true.
func (r Reporter) Printf(format string, v ...interface{}) {
	if r.verbose {
		_, _ = fmt.Fprintf(r.writer, format, v...)
	}
}

// NewReporter creates an instance of Reporter.
//
// The final variadic io.Writer argument allows passing ONE specific writer
// and defaulting to os.Stderr when none is passed.
func NewReporter(verbose bool, w ...io.Writer) Reporter {
	var writer io.Writer
	if len(w) == 0 {
		writer = os.Stderr
	} else {
		writer = w[0]
	}

	return Reporter{
		verbose: verbose,
		writer:  writer,
	}
}

// Dial attempts to connect to the specified server for the specified duration,
// using the specified connection method.
//
// It returns a status code usable with os.Exit().
func Dial(url string, maxTimeout time.Duration, dialer DriverDial, r Reporter) Result {
	const Nanoseconds = float64(time.Second)

	t0 := time.Now()
	tMax := t0.Add(maxTimeout)

	timeout := DefaultBaseTimeout

	for {
		err := dialer(url, timeout)

		// Success.
		if err == nil {
			r.Printf("Connected in %d msec.\n", time.Since(t0)/1.0e6)
			return Success
		}

		// Ignore expected errors, but return on others.
		if err.Error() != ExpectedErrorString && !errors.Is(err, context.DeadlineExceeded) {
			r.Printf("Unexpected error %v.\n", err)
			return OtherError
		}

		// Timeout reached.
		t1 := time.Now()
		if t1.After(tMax) {
			return Timeout
		}

		timeout = backOff(timeout)
		r.Printf("Unavailable in %3f seconds, retrying in %.3f seconds.\n",
			float64(t1.Sub(t0))/Nanoseconds,
			float64(timeout)/Nanoseconds)
		time.Sleep(timeout)
	}
}
