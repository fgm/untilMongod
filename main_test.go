package main

import (
	"bytes"
	"errors"
	"flag"
	"strings"
	"testing"
	"time"

	"github.com/fgm/untilMongod/dial"
)

// happyDialer always succeeds
func happyDialer(_ string, _ time.Duration) error {
	return nil
}

// sadDialer always fails with a "no reachable servers" error
func sadDialer(_ string, _ time.Duration) error {
	return errors.New(dial.ExpectedErrorString)
}

// otherErrorDialer always fails with a generic error
func otherErrorDialer(_ string, _ time.Duration) error {
	return errors.New("some other error")
}

func TestRealMain(t *testing.T) {
	t.Parallel()
	tests := [...]struct {
		name         string
		args         []string
		dialer       dial.DriverDial
		expectedCode int
		expectedOut  string
	}{
		{
			name:         "Success",
			args:         []string{"-url", "mongodb://localhost:27017"},
			dialer:       happyDialer,
			expectedCode: 0,
			expectedOut:  "", // Not verbose, so no output
		},
		{
			name:         "SuccessVerbose",
			args:         []string{"-url", "mongodb://localhost:27017", "-v"},
			dialer:       happyDialer,
			expectedCode: 0,
			expectedOut:  "Connected in", // Verbose, so should have output
		},
		{
			name:         "Timeout",
			args:         []string{"-url", "mongodb://localhost:27017", "-timeout", "1"},
			dialer:       sadDialer,
			expectedCode: 1,  // dial.Timeout
			expectedOut:  "", // Not verbose
		},
		{
			name:         "TimeoutVerbose",
			args:         []string{"-url", "mongodb://localhost:27017", "-timeout", "1", "-v"},
			dialer:       sadDialer,
			expectedCode: 1, // dial.Timeout
			expectedOut:  "Unavailable in",
		},
		{
			name:         "OtherError",
			args:         []string{"-url", "mongodb://localhost:27017"},
			dialer:       otherErrorDialer,
			expectedCode: 2, // dial.OtherError
			expectedOut:  "",
		},
		{
			name:         "OtherErrorVerbose",
			args:         []string{"-url", "mongodb://localhost:27017", "-v"},
			dialer:       otherErrorDialer,
			expectedCode: 2, // dial.OtherError
			expectedOut:  "Unexpected error",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			var out bytes.Buffer
			fs := flag.NewFlagSet(test.name, flag.ContinueOnError)

			code := realMain(&out, fs, test.args, test.dialer)

			if code != test.expectedCode {
				t.Errorf("expected exit code %d, got %d", test.expectedCode, code)
			}

			if test.expectedOut != "" && !strings.Contains(out.String(), test.expectedOut) {
				t.Errorf("expected output to contain %q, got %q", test.expectedOut, out.String())
			}

			if test.expectedOut == "" && out.String() != "" {
				t.Errorf("expected no output, got %q", out.String())
			}
		})
	}
}
