/*
Copyright 2020 Jim Zhang (jim.zoumo@gmail.com). All rights reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package netutil

import (
	"errors"
	"fmt"
	"net"
	"sync"

	utilerrors "k8s.io/apimachinery/pkg/util/errors"
)

var (
	ErrAccecptClosed     = &net.OpError{Op: "accept", Err: fmt.Errorf("use of closed network connection")}
	ErrNoEnoughListeners = errors.New("must supply at least two listeners")
)

// AggregatedListener is a listener aggregated by other listeners in
// order to satisfy net.Listener interface.
//
// It makes http server can accept connections from several listeners in the
// meanwhile.
type AggregatedListener interface {
	net.Listener

	// Addrs returns all listeners' network addresses.
	Addrs() []net.Addr

	// TCPAddrs returns all tcp listeners' network addresses.
	TCPAddrs() []*net.TCPAddr

	// UnixAddrs returns all unix listeners' network addresses.
	UnixAddrs() []*net.UnixAddr

	// AcceptTCP accepts the next tcp incoming call and returns the new
	// tcp connection.
	AcceptTCP() (*net.TCPConn, error)

	// AcceptUnix accepts the next unix incoming call and returns the new
	// unix connection.
	AcceptUnix() (*net.UnixConn, error)
}

type acceptResult struct {
	conn net.Conn
	err  error
}

type aggregatedListener struct {
	major net.Listener

	acceptC     chan *acceptResult
	acceptTCPC  chan *acceptResult
	acceptUnixC chan *acceptResult

	tcpLns  []*net.TCPListener
	unixLns []*net.UnixListener
	lns     []net.Listener

	closeOnece   sync.Once
	closeAcceptC chan struct{}
	closeC       chan struct{}
}

// NewAggregatedListener aggregate all input listeners into one to
// satisfy net.Listener interface.
//
// Must supply at least two listeners.
//
// It takes the first listener as major to expose the address and
// accepts all listeners on background to get network connections.
func NewAggregatedListener(listeners ...net.Listener) (AggregatedListener, error) {
	if len(listeners) < 2 {
		return nil, ErrNoEnoughListeners
	}
	l := &aggregatedListener{
		acceptC:      make(chan *acceptResult),
		acceptTCPC:   make(chan *acceptResult),
		acceptUnixC:  make(chan *acceptResult),
		closeC:       make(chan struct{}),
		closeAcceptC: make(chan struct{}),
		major:        listeners[0],
	}

	for i := range listeners {
		switch ln := listeners[i].(type) {
		case *net.TCPListener:
			l.tcpLns = append(l.tcpLns, ln)
		case *net.UnixListener:
			l.unixLns = append(l.unixLns, ln)
		default:
			l.lns = append(l.lns, ln)
		}
	}
	l.acceptBackgroup()
	return l, nil
}

func (l *aggregatedListener) acceptBackgroup() {
	wg := sync.WaitGroup{}

	for i := range l.lns {
		ll := l.lns[i]
		wg.Add(1)
		go func() {
			l.acceptFromListener(ll, l.acceptC)
			wg.Done()
		}()
	}
	for i := range l.tcpLns {
		ll := l.tcpLns[i]
		wg.Add(1)
		go func() {
			l.acceptFromListener(ll, l.acceptTCPC)
			wg.Done()
		}()
	}
	for i := range l.unixLns {
		ll := l.unixLns[i]
		wg.Add(1)
		go func() {
			l.acceptFromListener(ll, l.acceptUnixC)
			wg.Done()
		}()
	}

	go func() {
		wg.Wait()
		// close channel after all listener accecpt goroutine stopped
		close(l.closeC)
	}()
}

func (l *aggregatedListener) acceptFromListener(ln net.Listener, resultChan chan<- *acceptResult) {
	for {
		conn, err := ln.Accept()

		needBreak := false
		if err != nil {
			needBreak = true
			if nerr, ok := err.(net.Error); ok && nerr.Temporary() {
				needBreak = false
			}
		}

		select {
		case resultChan <- &acceptResult{conn, err}:
			if needBreak {
				return
			}
		case <-l.closeAcceptC:
			return
		}
	}
}

// AcceptTCP accepts the next tcp incoming call and returns the new
// tcp connection.
func (l *aggregatedListener) AcceptTCP() (*net.TCPConn, error) {
	select {
	case result := <-l.acceptTCPC:
		if result.err != nil {
			return nil, result.err
		}
		return result.conn.(*net.TCPConn), nil
	case <-l.closeC:
		return nil, ErrAccecptClosed
	}
}

// AcceptUnix accepts the next unix incoming call and returns the new
// unix connection.
func (l *aggregatedListener) AcceptUnix() (*net.UnixConn, error) {
	select {
	case result := <-l.acceptUnixC:
		if result.err != nil {
			return nil, result.err
		}
		return result.conn.(*net.UnixConn), nil
	case <-l.closeC:
		return nil, ErrAccecptClosed
	}
}

// Accept implements the Accept method in the Listener interface; it
// waits for the next call and returns a generic Conn.
func (l *aggregatedListener) Accept() (net.Conn, error) {
	var result *acceptResult
	select {
	case result = <-l.acceptC:
	case result = <-l.acceptTCPC:
	case result = <-l.acceptUnixC:
	case <-l.closeC:
		return nil, ErrAccecptClosed
	}
	return result.conn, result.err
}

func (l *aggregatedListener) Addr() net.Addr {
	return l.major.Addr()
}

func (l *aggregatedListener) Close() error {
	var closeErr error
	l.closeOnece.Do(func() {
		var errors []error
		for _, ln := range l.lns {
			if err := ln.Close(); err != nil {
				errors = append(errors, err)
			}
		}
		for _, ln := range l.tcpLns {
			if err := ln.Close(); err != nil {
				errors = append(errors, err)
			}
		}
		for _, ln := range l.unixLns {
			if err := ln.Close(); err != nil {
				errors = append(errors, err)
			}
		}
		closeErr = utilerrors.NewAggregate(errors)
		close(l.closeAcceptC)
	})
	return closeErr
}

// Addrs returns all listeners' network addresses.
func (l *aggregatedListener) Addrs() []net.Addr {
	var addrs []net.Addr

	for _, ln := range l.lns {
		addrs = append(addrs, ln.Addr())
	}
	for _, ln := range l.tcpLns {
		addrs = append(addrs, ln.Addr())
	}
	for _, ln := range l.unixLns {
		addrs = append(addrs, ln.Addr())
	}
	return addrs
}

// TCPAddrs returns all tcp listeners' network addresses.
func (l *aggregatedListener) TCPAddrs() []*net.TCPAddr {
	var addrs []*net.TCPAddr
	for _, ln := range l.tcpLns {
		addrs = append(addrs, ln.Addr().(*net.TCPAddr))
	}
	return addrs
}

// UnixAddrs returns all unix listeners' network addresses.
func (l *aggregatedListener) UnixAddrs() []*net.UnixAddr {
	var addrs []*net.UnixAddr
	for _, ln := range l.tcpLns {
		addrs = append(addrs, ln.Addr().(*net.UnixAddr))
	}
	return addrs
}
