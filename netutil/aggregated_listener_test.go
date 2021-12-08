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
	"net"
	"os"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// testUnixAddr uses os.CreateTemp to get a name that is unique.
func testUnixAddr() string {
	f, err := os.CreateTemp("", "go-nettest")
	if err != nil {
		panic(err)
	}
	addr := f.Name()
	f.Close()
	os.Remove(addr)
	return addr
}

func TestNewAggregatedLister(t *testing.T) {
	type args struct {
		listeners []net.Listener
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"empty lister is not allowed",
			args{
				listeners: nil,
			},
			true,
		},
		{
			"one lister is not allowed",
			args{
				listeners: []net.Listener{nil},
			},
			true,
		},
	}
	for i := range tests {
		tt := tests[i]
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewAggregatedListener(tt.args.listeners...)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewAggregatedLister() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func createTestAggregatedLister(t *testing.T) (AggregatedListener, *net.TCPListener, *net.UnixListener) {
	tcpLn, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}

	unixName := testUnixAddr()
	unixLn, err := net.Listen("unix", unixName)
	if err != nil {
		t.Fatal(err)
	}

	ln, err := NewAggregatedListener(tcpLn, unixLn)
	if err != nil {
		t.Fatal(err)
	}

	return ln, tcpLn.(*net.TCPListener), unixLn.(*net.UnixListener)
}

func TestAggregatedListener_Accept(t *testing.T) {
	ln, tcpLn, unixLn := createTestAggregatedLister(t)

	const N = 10
	var wg sync.WaitGroup
	wg.Add(N)
	var got int32
	for i := 0; i < N; i++ {
		go func() {
			defer wg.Done()
			for {
				c, err := ln.Accept()
				if errors.Is(err, ErrAccecptClosed) {
					return
				}
				if err != nil {
					t.Logf("err: %v", err)
					continue
				}
				atomic.AddInt32(&got, 1)
				c.Close()
			}
		}()
	}

	attempts := N * 10
	fails := 0
	d := &net.Dialer{Timeout: 200 * time.Millisecond}
	for i := 0; i < attempts; i++ {
		c, err := d.Dial("tcp", tcpLn.Addr().String())
		if err != nil {
			fails++
		} else {
			c.Close()
		}
	}
	for i := 0; i < attempts; i++ {
		c, err := d.Dial("unix", unixLn.Addr().String())
		if err != nil {
			fails++
		} else {
			c.Close()
		}
	}

	// time.Sleep(time.Second)
	ln.Close()
	wg.Wait()

	if fails > 0 {
		t.Logf("# of failed Dials: %v", fails)
	}

	if got != int32(attempts*2) {
		t.Fatalf("got = %v, want = %v", got, attempts*2)
	}
}

func TestAggregatedListener_AcceptTCP(t *testing.T) {
	ln, tcpLn, unixLn := createTestAggregatedLister(t)

	const N = 10
	var wg sync.WaitGroup
	wg.Add(N)
	var got int32
	for i := 0; i < N; i++ {
		go func() {
			defer wg.Done()
			for {
				c, err := ln.AcceptTCP()
				if errors.Is(err, ErrAccecptClosed) {
					return
				}
				if err != nil {
					t.Logf("err: %v", err)
					continue
				}
				atomic.AddInt32(&got, 1)
				c.Close()
			}
		}()
	}

	attempts := N * 10
	fails := 0
	d := &net.Dialer{Timeout: 200 * time.Millisecond}
	for i := 0; i < attempts; i++ {
		c, err := d.Dial("tcp", tcpLn.Addr().String())
		if err != nil {
			fails++
		} else {
			c.Close()
		}
	}
	for i := 0; i < attempts; i++ {
		c, err := d.Dial("unix", unixLn.Addr().String())
		if err != nil {
			fails++
		} else {
			c.Close()
		}
	}
	// time.Sleep(time.Second)
	ln.Close()
	wg.Wait()

	if fails > 0 {
		t.Logf("# of failed Dials: %v", fails)
	}

	if got != int32(attempts) {
		t.Fatalf("got = %v, want = %v", got, attempts)
	}
}

func TestAggregatedListener_AcceptUnix(t *testing.T) {
	ln, tcpLn, unixLn := createTestAggregatedLister(t)

	const N = 10
	var wg sync.WaitGroup
	wg.Add(N)
	var got int32
	for i := 0; i < N; i++ {
		go func() {
			defer wg.Done()
			for {
				c, err := ln.AcceptUnix()
				if errors.Is(err, ErrAccecptClosed) {
					return
				}
				if err != nil {
					t.Logf("err: %v", err)
					continue
				}
				atomic.AddInt32(&got, 1)
				c.Close()
			}
		}()
	}

	attempts := N * 10
	fails := 0
	d := &net.Dialer{Timeout: 200 * time.Millisecond}
	for i := 0; i < attempts; i++ {
		c, err := d.Dial("tcp", tcpLn.Addr().String())
		if err != nil {
			fails++
		} else {
			c.Close()
		}
	}
	for i := 0; i < attempts; i++ {
		c, err := d.Dial("unix", unixLn.Addr().String())
		if err != nil {
			fails++
		} else {
			c.Close()
		}
	}
	// time.Sleep(time.Second)
	ln.Close()
	wg.Wait()

	if fails > 0 {
		t.Logf("# of failed Dials: %v", fails)
	}

	if got != int32(attempts) {
		t.Fatalf("got = %v, want = %v", got, attempts)
	}
}
