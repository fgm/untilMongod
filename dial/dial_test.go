package dial

import (
	"errors"
	"regexp"
	"strings"
	"testing"
	"time"
)

func TestDial_happy(t *testing.T) {
	t.Parallel()

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
	t.Parallel()

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
	t.Parallel()

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
	assertSilence := func(b *strings.Builder, format string) {
		if b.Len() != 0 {
			t.Errorf(format, b.String())
		}
	}
	assertRetries := func(b *strings.Builder, format string) {
		var retryRegex = regexp.MustCompile(`^Unavailable in .* seconds, retrying in .* seconds\.`)
		actual := b.String()
		if !retryRegex.MatchString(actual) {
			t.Errorf(format, actual)
		}
	}

	checks := []struct {
		name    string
		timeout time.Duration
		test    func(b *strings.Builder, format string)
		format  string
	}{
		// Immediate success should not output anything
		{"short", 1, assertSilence, "Success without a retry should not output anything: %q"},
		// Retried response should output a retry message.
		{"normal", time.Second, assertRetries, "Verbose output on a retry did not match the expected pattern: %#v"},
	}

	for _, check := range checks {
		t.Run(check.name, func(t *testing.T) {
			b := &strings.Builder{}
			reporter := NewReporter(true, b)
			defer func() {
				if t.Failed() {
					t.Log(t.Name(), b.String())
				}
			}()
			Dial(DefaultURL, check.timeout, slowDialer, reporter)
			check.test(b, check.format)
		})
	}
}

func TestDialIntegration(t *testing.T) {
	const Timeout = 100 * time.Millisecond
	if testing.Short() {
		t.Skip("Skipping slow test for NewMongoDbV2Dial")
	}
	t.Parallel()
	dialer := NewMongoDbDial()
	checks := [...]struct {
		name     string
		url      string
		expected string
	}{
		// err will be non-nil since there is no server at "": this is not an error.
		{"empty url", "", "error parsing uri"},
		// This may succeed if there is a server on the default MongoDB URL.
		{"default url", DefaultURL, "server selection"},
	}

	for _, check := range checks {
		t.Run(check.name, func(t *testing.T) {
			err := dialer(check.url, Timeout)
			if err != nil {
				if !regexp.MustCompile(check.expected).MatchString(err.Error()) {
					t.Errorf("NewMongoDbDial returned an unexpected error: %v", err)
				}
			}
		})
	}
}
