package retry

import (
	"time"

	"k8s.io/apimachinery/pkg/util/wait"
)

// ErrWaitTimeout is returned when the condition exited without success.
var ErrWaitTimeout = wait.ErrWaitTimeout

// Backoff holds parameters applied to a Backoff function.
// redefined here for convenience
type Backoff = wait.Backoff

// DefaultRetry is the recommended retry for a conflict where multiple clients
// are making changes to the same resource.
var DefaultRetry = Backoff{
	Steps:    5,
	Duration: 10 * time.Millisecond,
	Factor:   1.0,
	Jitter:   0.1,
}

// DefaultBackoff is the recommended backoff for a conflict where a client
// may be attempting to make an unrelated modification to a resource under
// active management by one or more controllers.
var DefaultBackoff = Backoff{
	Steps:    4,
	Duration: 10 * time.Millisecond,
	Factor:   5.0,
	Jitter:   0.1,
}

// Retry executes the provided condition func repeatedly, retrying with exponential
// backoff if the condition func returns an error.
//
// It checks the condition up to Steps times, increasing the wait by multiplying
// the previous duration by Factor.
//
// If Jitter is greater than zero, a random amount of each duration is added
// (between duration and duration*(1+jitter)).
//
// If the retrying timeout, the last error of condition will be returned
func Retry(backoff Backoff, condition func() error) error {
	return retry(backoff, condition, nil, nil)
}

// RetryOn does the same thing with Retry() except that it will keep trying if
// the error returned by condition is within the expected by continueOn function
//
// Deprecated: Use RetryContined
func RetryOn(backoff Backoff, condition func() error, continueOn func(error) bool) error {
	return retry(backoff, condition, nil, continueOn)
}

// RetryContined does the same thing with Retry() except that it will keep trying if
// the error returned by condition is within the expected by continued function
func RetryContined(backoff Backoff, condition func() error, continued func(error) bool) error {
	return retry(backoff, condition, nil, continued)
}

// RetryIgnored does the same thing with Retry() but it will stop retrying when conidtion returns
// an error which will be ignored if it is within the expected.
func RetryIgnored(backoff Backoff, condition func() error, ignored func(error) bool) error {
	return retry(backoff, condition, ignored, nil)
}

// retry executes the provided condition func repeatedly, retrying with exponential
// backoff if the condition func returns an error.
//
// If the ignored function returns true, it will interrupte the retrying loop.
//
// If the conitnued fucntion returns true, retrying will keep going.
//
// If the retrying timeout, the last error of condition will be returned.
func retry(backoff Backoff, condition func() error, ignored, continued func(error) bool) error {
	var lastErr error
	err := wait.ExponentialBackoff(backoff, func() (bool, error) {
		err := condition()
		if err == nil {
			return true, nil
		}
		if ignored != nil && ignored(err) {
			return true, nil
		}
		// record lastErr
		lastErr = err
		if continued != nil && continued(err) {
			return false, nil
		}
		return false, err
	})
	if err == wait.ErrWaitTimeout {
		err = lastErr
	}
	return err
}
