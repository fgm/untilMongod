package dial

import (
	"testing"
	"time"

	"errors"
	"regexp"
	"strings"
)

func TestDial_happy(t *testing.T) {
	var actual Result
	var happyDialer DriverDial

	happyDialer = func(url string, duration time.Duration) error {
		return nil
	}

	var attempts = []struct {
		timeout  time.Duration
		expected Result
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
			reporter := NewReporter(false)
			actual = Dial("", attempt.timeout, happyDialer, reporter)
			if actual != attempt.expected {
				t.Logf("Expected %v, got %v\n", attempt.expected, actual)
				t.Fail()
			}
		})
	}
}

func TestDial_sad(t *testing.T) {
	var actual Result
	var sadDialer DriverDial

	sadDialer = func(url string, duration time.Duration) error {
		return errors.New("sad")
	}

	var attempts = []struct {
		timeout  time.Duration
		expected Result
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
			reporter := NewReporter(false)
			actual = Dial("", attempt.timeout, sadDialer, reporter)
			if actual != attempt.expected {
				t.Logf("Expected %v, got %v\n", attempt.expected, actual)
				t.Fail()
			}
		})
	}
}

func TestDial_slow(t *testing.T) {
	var actual Result
	var slowDialer DriverDial

	slowDialer = func(url string, timeout time.Duration) error {
		time.Sleep(1 * time.Millisecond)
		if timeout > 10*time.Millisecond {
			return nil
		}
		return errors.New(ExpectedErrorString)
	}

	var attempts = []struct {
		timeout  time.Duration
		expected Result
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
			reporter := NewReporter(false)
			actual = Dial("", attempt.timeout, slowDialer, reporter)
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
	var slowDialer DriverDial

	slowDialer = func(url string, timeout time.Duration) error {
		if timeout > 10*time.Millisecond {
			time.Sleep(timeout)
			return nil
		}
		return errors.New(ExpectedErrorString)
	}

	var b strings.Builder

	b = strings.Builder{}
	reporter := NewReporter(true, &b)

	Dial("", 1, slowDialer, reporter)
	if b.Len() != 0 {
		t.Error("Success without a retry should not output anything.")
	}

	// Retried response should output a retry message.
	Dial("", 10*time.Second, slowDialer, reporter)
	var retryRegex = regexp.MustCompile(`^Unavailable in .* seconds, retrying in .* seconds\.`)
	var output = b.String()
	if !retryRegex.MatchString(output) {
		t.Error("Verbose output on a retry did not match the expected pattern")
	}

}
