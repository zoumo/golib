// Copyright 2022 jim.zoumo@gmail.com
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

package chanx

import (
	"sync"
)

// Options defines the functional option type for Channel
type Options func(*config)

type config struct {
	inChanSize           int
	outChanSize          int
	initBufferSize       int
	maxBufferSize        int
	dropClosedBufferData bool
}

// InChanSize sets input channel buffer size
// It means that input := make(chan interface{}, inChanSize)
func InChanSize(size int) Options {
	return func(c *config) {
		if size >= 0 {
			c.inChanSize = size
		}
	}
}

// OutChanSzie sets output channel buffer size
// It means that output := make(chan interface{}, outChanSize)
func OutChanSzie(size int) Options {
	return func(c *config) {
		if size >= 0 {
			c.outChanSize = size
		}
	}
}

// InitBufferSize sets the ring buffer initial size
func InitBufferSize(size int) Options {
	return func(c *config) {
		if size > 0 {
			c.initBufferSize = size
		}
	}
}

// MaxBufferSize sets the ring buffer max size
// If set to 0, it means no limit.
func MaxBufferSize(size int) Options {
	return func(c *config) {
		c.maxBufferSize = size
	}
}

// DropClosedBufferData specifies that the data in ring buffer and
// input channel buffer will be dropped after Close() called.
// The data in output channel buffer can still be accessed.
func DropClosedBufferData() Options {
	return func(c *config) {
		c.dropClosedBufferData = true
	}
}

func newDefuerConfig() *config {
	return &config{
		initBufferSize:       2,
		dropClosedBufferData: false,
	}
}

// ChannX is a self adaptive channel with a ring buffer.
// The channel buffer capacity will automatically increase according
// to excessive input and restore to original when buffer is empty.
// It can be used as an Unbounded Channel.
type ChannX struct {
	in        chan interface{}
	out       chan interface{}
	close     chan struct{}
	clsoeOnce sync.Once
	cfg       *config
	buffer    *SelfAdaptiveRingBuffer
}

func New(opts ...Options) *ChannX {
	cfg := newDefuerConfig()
	for _, opt := range opts {
		opt(cfg)
	}

	ch := &ChannX{
		cfg:   cfg,
		close: make(chan struct{}),
	}
	ch.in = make(chan interface{}, cfg.inChanSize)
	ch.out = make(chan interface{}, cfg.outChanSize)
	ch.buffer = NewSelfAdptiveRingBuffer(cfg.initBufferSize, cfg.maxBufferSize)

	go ch.process()
	return ch
}

func (ch *ChannX) process() {
	var v interface{}
	var ok bool
	for {
		// buffer is empty
		select {
		case v, ok = <-ch.in:
			if !ok {
				panic("chanx: input channel can not be closed")
			}
			if !ch.processObjectFromInput(v) {
				return
			}
		case <-ch.close:
			ch.processTermination(nil)
			return
		}

		// v has already processed
		// we need to deal with objects in buffer

		for !ch.buffer.IsEmpty() {
			peek, _ := ch.buffer.Peek()
			select {
			case v, ok = <-ch.in:
				if !ok {
					panic("chanx: input channel can not be closed")
				}
				if !ch.processObjectFromInput(v) {
					return
				}
			case ch.out <- peek:
				ch.buffer.Pop() // nolint
				if ch.buffer.NeedReset() {
					ch.buffer.Reset()
				}
			case <-ch.close:
				ch.processTermination(nil)
				return
			}
		}
	}
}

// process object from input channel, it will performs transformation and filter
func (ch *ChannX) processObjectFromInput(v interface{}) bool {
	if ch.buffer.IsEmpty() {
		// try to send v through channel directly
		select {
		case ch.out <- v:
			return true
		default:
			// output channel is full, put item to buffer
		}
	}

	// try send to buffer
	if !ch.mustPutToBuffer(v) {
		ch.processTermination(v)
		return false
	}
	return true
}

func (ch *ChannX) processTermination(poped interface{}) {
	close(ch.in)
	defer close(ch.out)

	if ch.cfg.dropClosedBufferData {
		// drop all data after closed
		ch.buffer.Reset()
		return
	}

	// We need to send data still in ringbuffer and
	// input channel buffer sequentially

	// send all item in buffer
	for !ch.buffer.IsEmpty() {
		v, _ := ch.buffer.Pop()
		ch.out <- v
	}
	ch.buffer.Reset()

	// poped is the latest poped item from input channel
	// it should be processed before others in channel buffer
	if poped != nil {
		ch.out <- poped
	}

	// send all item in input channel
	for v := range ch.in {
		ch.out <- v
	}
}

// Try to put item into buffer.
// If buffer is full, it wait util the peek of buffer is sent
// to output channel.
func (ch *ChannX) mustPutToBuffer(v interface{}) bool {
	if ch.buffer.Put(v) {
		return true
	}

	// buffer is full
	peek, _ := ch.buffer.Peek()

	select {
	case ch.out <- peek:
	case <-ch.close:
		return false
	}

	ch.buffer.Pop() //nolint
	ch.buffer.Put(v)
	return true
}

func (ch *ChannX) In() chan<- interface{} {
	return ch.in
}

func (ch *ChannX) Out() <-chan interface{} {
	return ch.out
}

func (ch *ChannX) Close() {
	ch.clsoeOnce.Do(func() {
		close(ch.close)
	})
}
