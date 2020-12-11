package exec

import (
	"fmt"
	"os/exec"
	"time"
)

type ProbeError struct {
	err error
}

func (e *ProbeError) Error() string {
	return fmt.Sprintf("probe failed: %v", e.err)
}

type probeResult bool

const (
	success probeResult = true
	failure probeResult = false
)

// Probe describes a health check to be performed against a container to determine whether it is
// alive or ready to receive traffic.
type Probe struct {
	// The action taken to determine the health of a cmd
	// Default to IsCmdRunningHandler
	Handler func(running *exec.Cmd) error `json:"handler,inline" protobuf:"bytes,1,opt,name=handler"`
	// Number of seconds after the container has started before liveness probes are initiated.
	// +optional
	InitialDelaySeconds int `json:"initialDelaySeconds,omitempty" protobuf:"varint,2,opt,name=initialDelaySeconds"`
	// How often (in seconds) to perform the probe.
	// Default to 1 seconds. Minimum value is 1.
	// +optional
	PeriodSeconds int `json:"periodSeconds,omitempty" protobuf:"varint,4,opt,name=periodSeconds"`
	// Minimum consecutive successes for the probe to be considered successful after having failed.
	// Defaults to 2. Must be 1 for liveness. Minimum value is 1.
	// +optional
	SuccessThreshold int `json:"successThreshold,omitempty" protobuf:"varint,5,opt,name=successThreshold"`
	// Minimum consecutive failures for the probe to be considered failed after having succeeded.
	// Defaults to 3. Minimum value is 1.
	// +optional
	FailureThreshold int `json:"failureThreshold,omitempty" protobuf:"varint,6,opt,name=failureThreshold"`
}

type worker struct {
	runningCmd *exec.Cmd

	probe *Probe

	startAt time.Time

	lastResult probeResult
	resultRun  int

	lastErr error
	stopCh  <-chan struct{}
}

func newWorker(
	runningCmd *exec.Cmd,
	p *Probe,
	startTime time.Time,
	stopCh <-chan struct{},
) *worker {
	return &worker{
		runningCmd: runningCmd,
		probe:      p,
		startAt:    startTime,
		lastResult: failure,
		stopCh:     stopCh,
	}
}

func (w *worker) run() chan error {
	resultC := make(chan error)

	go func() {
		probeTickerPeriod := time.Duration(w.probe.PeriodSeconds) * time.Second
		probeTicker := time.NewTicker(probeTickerPeriod)
		defer func() {
			probeTicker.Stop()
		}()

	probeLoop:
		for w.doProbe() {
			select {
			case <-w.stopCh:
				break probeLoop
			case <-probeTicker.C:
				// continue
			}
		}

		if w.lastErr != nil {
			resultC <- &ProbeError{
				err: w.lastErr,
			}
		} else {
			resultC <- nil
		}
	}()

	return resultC
}

func (w *worker) doProbe() (keepGoing bool) {
	if int(time.Since(w.startAt).Seconds()) < w.probe.InitialDelaySeconds {
		return true
	}

	err := w.probe.Handler(w.runningCmd)
	result := success
	if err != nil {
		w.lastErr = err
		result = failure
	}

	if w.lastResult == result {
		w.resultRun++
	} else {
		w.lastResult = result
		w.resultRun = 1
	}

	if (result == failure && w.resultRun < w.probe.FailureThreshold) ||
		(result == success && w.resultRun < w.probe.SuccessThreshold) {
		return true
	}

	return false
}
