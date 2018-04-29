package dial

import (
	"testing"
	"time"

	"errors"
	"regexp"
	"strings"
)

func TestDial_happy(t *testing.T) {
	var actual DialResult

	var happyDialer DriverDial = func(url string, duration time.Duration) error {
		return nil
	}

	var attempts = []struct {
		timeout  time.Duration
		expected DialResult
		message  string
	}{
		{1, Success, "Valid timeout"},
		{0, Success, "Zero timeout"},
		{-1, Success, "Negative timeout"},
	}

	t.Parallel()
	for _, attempt := range attempts {
		t.Run(attempt.message, func(t *testing.T) {
			t.Parallel()
			actual = Dial("", attempt.timeout, false, happyDialer)
			if actual != attempt.expected {
				t.Logf("Expected %v, got %v\n", attempt.expected, actual)
				t.Fail()
			}
		})
	}
}

func TestDial_sad(t *testing.T) {
	var actual DialResult

	var sadDialer DriverDial = func(url string, duration time.Duration) error {
		return errors.New("Sad")
	}

	var attempts = []struct {
		timeout  time.Duration
		expected DialResult
		message  string
	}{
		{1, OtherError, "Valid timeout"},
		{0, OtherError, "Zero timeout"},
		{-1, OtherError, "Negative timeout"},
	}

	t.Parallel()
	for _, attempt := range attempts {
		t.Run(attempt.message, func(t *testing.T) {
			t.Parallel()
			actual = Dial("", attempt.timeout, false, sadDialer)
			if actual != attempt.expected {
				t.Logf("Expected %v, got %v\n", attempt.expected, actual)
				t.Fail()
			}
		})
	}
}

func TestDial_slow(t *testing.T) {
	var actual DialResult

	var slowDialer DriverDial = func(url string, timeout time.Duration) error {
		time.Sleep(1 * time.Millisecond)
		if timeout > 10*time.Millisecond {
			return nil
		}
		return errors.New(ExpectedErrorString)
	}

	var attempts = []struct {
		timeout  time.Duration
		expected DialResult
		message  string
	}{
		{1 * time.Nanosecond, Timeout, "Short timeout"},
		{1 * time.Second, Success, "Long-enough timeout"},
		{0, Timeout, "Zero timeout"},
		{-1, Timeout, "Negative timeout"},
	}

	t.Parallel()
	for _, attempt := range attempts {
		t.Run(attempt.message, func(t *testing.T) {
			t.Parallel()
			actual = Dial("", attempt.timeout, false, slowDialer)
			if actual != attempt.expected {
				t.Logf("Expected %v, got %v\n", attempt.expected, actual)
				t.Fail()
			}
		})
	}
}

func TestDial_verbose(t *testing.T) {
	t.Parallel()

	// This dialer only succeeds when given at least 10 msec as a timeout.
	var slowDialer DriverDial = func(url string, timeout time.Duration) error {
		if timeout > 10*time.Millisecond {
			time.Sleep(timeout)
			return nil
		}
		return errors.New(ExpectedErrorString)
	}

	var b strings.Builder

	b = strings.Builder{}
	Dial("", 1, true, slowDialer, &b)
	if b.Len() != 0 {
		t.Error("Success without a retry should not output anything.")
	}

	// Retried response should output a retry message.
	Dial("", 10000*time.Millisecond, true, slowDialer, &b)
	var retryRegex = regexp.MustCompile(`^Unavailable in .* seconds, retrying in .* seconds\.`)
	var output = b.String()
	if !retryRegex.MatchString(output) {
		t.Error("Verbose output on a retry did not match the expected pattern")
	}

}
