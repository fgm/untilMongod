package dial

import "time"

const backoffFactor = 1.5

// See time.maxDuration. Cannot use it since it is private.
const maxDuration time.Duration = 1<<63 - 1

// A valid maximum timeout must be positive and not too large.
func isMaxTimeoutValid(t time.Duration) bool {
	if t <= 0*time.Nanosecond {
		return false
	}

	limit := float64(maxDuration) / backoffFactor
	if t >= time.Duration(limit) {
		return false
	}

	return true
}

// FIXME Provide an increasing delay from an existing one.
func backOff(current time.Duration) time.Duration {

	floatNext := float64(current) * backoffFactor
	return time.Duration(floatNext)
}
