// Copyright 2023 jim.zoumo@gmail.com
// 
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// 
//     http://www.apache.org/licenses/LICENSE-2.0
// 
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package scp

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"sync/atomic"
	"time"
)

func newReader(ioReader io.Reader, timeout time.Duration) *reader {
	const processors = 3
	return &reader{
		reader:     ioReader,
		outputChan: make(chan []byte, processors),
		closedChan: make(chan bool, processors),
		errorChan:  make(chan error, processors),
		timeout:    timeout,
	}
}

type reader struct {
	processors int
	reader     io.Reader
	outputChan chan []byte
	closedChan chan bool
	errorChan  chan error
	closed     uint32
	timeout    time.Duration
	buf        bytes.Buffer
}

func (r *reader) sendCloseNotification() {
	for i := 0; i < r.processors; i++ {
		select {
		case r.closedChan <- true:
		case <-time.After(time.Millisecond):
		}
	}
}

func (r *reader) readStatus() error {
	data, err := r.read()
	if err != nil {
		return err
	}
	status := data[0]
	switch status {
	case StatusOK:
		return nil
	default:
		return errors.New(strings.TrimSpace(string(data[1:])))
	}
}

func (r *reader) readFile(info os.FileInfo) (io.Reader, error) {
	for r.buf.Len() <= int(info.Size()) {
		data, err := r.read()
		if err != nil {
			return nil, err
		}
		r.buf.Write(data)
	}
	defer r.buf.Reset()
	data := r.buf.Bytes()
	overflow := data[info.Size():]
	if len(overflow) != 1 || overflow[0] != StatusOK {
		return nil, fmt.Errorf("invalid statusOK, expected: %v, but got: %v ", []byte{StatusOK}, overflow)
	}
	data = data[:info.Size()]

	return bytes.NewReader(data), nil
}

func (r *reader) isClosed() bool {
	return atomic.LoadUint32(&r.closed) == 1
}

func (r *reader) close() {
	if atomic.CompareAndSwapUint32(&r.closed, 0, 1) {
		r.sendCloseNotification()
		close(r.errorChan)
		close(r.closedChan)
		close(r.outputChan)
	}
}

func (r *reader) read() ([]byte, error) {
	if r.isClosed() {
		return nil, fmt.Errorf("closed")
	}

	select {
	case data := <-r.outputChan:
		return data, nil
	case err := <-r.errorChan:
		return nil, err
	case <-time.Tick(r.timeout):
		return nil, fmt.Errorf("exceeded timeout %s", r.timeout)
	case <-r.closedChan:
		return nil, fmt.Errorf("closed")
	}
}

func (r *reader) readInBackground() {
	for {
		var buffer = make([]byte, 4096)
		n, err := r.reader.Read(buffer)
		if err != nil {
			r.closeWithError(err)
			return
		}
		if r.isClosed() {
			return
		}
		if n == 0 {
			continue
		}
		select {
		case r.outputChan <- buffer[:n]:
			continue
		case <-r.closedChan:
			return
		}
	}
}

func (r *reader) closeWithError(err error) {
	if r.isClosed() {
		return
	}
	r.errorChan <- err
	r.sendCloseNotification()
}
