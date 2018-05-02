package dial

import "time"

// BackoffFactor is the interval increase coefficient for retries.
const BackoffFactor = 1.5

// MaxDuration is the largest Duration the time.Duration type can represent.
//
// Cannot use time.maxDuration since it is private in time package.
const MaxDuration time.Duration = 1<<63 - 1

// IsMaxTimeoutValid ensures timeouts passed to dialers are positive and not too large.
func IsMaxTimeoutValid(t time.Duration) bool {
	if t <= 0*time.Nanosecond {
		return false
	}

	limit := float64(MaxDuration) / BackoffFactor
	if t >= time.Duration(limit) {
		return false
	}

	return true
}

// FIXME Provide an increasing delay from an existing one.
func backOff(current time.Duration) time.Duration {

	floatNext := float64(current) * BackoffFactor
	return time.Duration(floatNext)
}
