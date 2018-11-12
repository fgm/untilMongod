package dial

import (
	"testing"
	"time"
)

func TestIsTimeoutValid(t *testing.T) {
	t.Parallel()
	limit := time.Duration(float64(MaxDuration) / BackoffFactor).Truncate(1 * time.Second)

	attempts := []struct {
		duration time.Duration
		expected bool
		message  string
	}{
		{0, false, "Zero duration is invalid"},
		{-1 * time.Second, false, "Negative duration cannot be waited for"},
		{MaxDuration, false, "Maximum duration cannot be extended"},
		{1, true, "Minimum duration is valid"},
		{limit, true, "Limit duration is valid"},
	}

	for _, attempt := range attempts {
		t.Run(attempt.message, func(t *testing.T) {
			// t.Parallel()
			actual := IsMaxTimeoutValid(attempt.duration)
			if actual != attempt.expected {
				t.Logf("Expected %v, got %v\n", attempt.expected, actual)
				t.Fail()
			}
		})
	}
}
