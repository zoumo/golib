package signal

import (
	"os"
	"os/signal"
	"syscall"
)

// Signal is a helper struct to let you handle signal eazily
type Signal struct {
	signalChan chan os.Signal
	exit       bool
}

// New returns a new Signal
func New(exit bool, sig ...os.Signal) *Signal {
	s := &Signal{
		signalChan: make(chan os.Signal, 1),
		exit:       exit,
	}

	signal.Notify(s.signalChan, sig...)
	return s
}

// Handle receives a signal and let the handler to handle it
// If exit is true, the program will exit with the exitCode.
// Otherwise, it will wait for the next signal arrivalling.
func (s *Signal) Handle(handler func(os.Signal) (exitCode int)) {
	go func() {
		for {
			sig := <-s.signalChan
			code := handler(sig)
			if s.exit {
				os.Exit(code)
			}
		}
	}()
}

// HandleSigterm helps you handle the SIGTERM
func HandleSigterm(handler func(os.Signal) (exitCode int)) {
	s := New(true, syscall.SIGTERM)
	s.Handle(handler)
}

// HandleSigint helps you handle the SIGINT
func HandleSigint(handler func(os.Signal) (exitCode int)) {
	s := New(true, syscall.SIGINT)
	s.Handle(handler)
}
