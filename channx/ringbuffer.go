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
	"math"
)

const (
	growThreshold = 1024
)

type SelfAdaptiveRingBuffer struct {
	buf      []interface{}
	maxSize  int
	initSize int
	size     int
	r        int // read position
	w        int // write position
	full     bool
}

// NewSelfAdptiveRingBuffer creates a self adaptive ringbuffer with init and max size.
// - The initSize must be greater than 0.
// - If maxSize <= 0, it means unlimited
func NewSelfAdptiveRingBuffer(initSize, maxSize int) *SelfAdaptiveRingBuffer {
	if initSize < 0 {
		initSize = 1
	}
	if maxSize < 0 {
		// unbounded ringbuffer
		maxSize = 0
	}
	return &SelfAdaptiveRingBuffer{
		buf:      make([]interface{}, initSize),
		initSize: initSize,
		maxSize:  maxSize,
		size:     initSize,
		r:        0,
		w:        0,
	}
}

func (rb *SelfAdaptiveRingBuffer) Put(v interface{}) bool {
	if rb.IsFull() {
		return false
	}

	rb.buf[rb.w] = v
	rb.w++

	if rb.w == rb.size {
		// out of range
		rb.w = 0
	}

	if rb.w == rb.r {
		// need grow
		if !rb.grow() {
			// can not grow any more
			rb.full = true
		}
	}
	return true
}

func (rb *SelfAdaptiveRingBuffer) grow() bool {
	newcap := rb.growCap()
	if newcap <= rb.size {
		return false
	}

	// copy data
	buf := make([]interface{}, newcap)

	// full buffer w == r
	// [ . . . . . . . . . . . . ]
	//         ^
	//       /  \
	//      w    r
	copy(buf[0:], rb.buf[rb.r:])             // copy old buf from r to end, new buf from 0 to size-r-1
	copy(buf[rb.size-rb.r:], rb.buf[0:rb.r]) // copy old buf from 0 to r-1, new buf from size-r to szie-1

	rb.r = 0       // read from new index
	rb.w = rb.size // next writeable index
	rb.size = newcap
	rb.buf = buf

	return true
}

func (rb *SelfAdaptiveRingBuffer) growCap() int {
	if rb.maxSize > 0 && rb.size >= rb.maxSize {
		// can not grow any more
		return rb.size
	}

	// reference to go18/runtime/slice/growslice
	newcap := rb.size
	doubleCap := newcap + newcap
	if doubleCap <= 0 {
		// overflow
		doubleCap = math.MaxInt
	}

	if rb.size < growThreshold {
		newcap = doubleCap
	} else {
		// Transition from growing 2x for small slices
		// to growing 1.25x for large slices. This formula
		// gives a smooth-ish transition between the two.
		add := newcap + 3*growThreshold
		if add <= 0 {
			// overflow
			newcap = math.MaxInt
		} else {
			newcap += add / 4
		}

		// Set newcap to math.MaxInt when
		// the newcap calculation overflowed.
		if newcap <= 0 {
			newcap = math.MaxInt
		}
	}

	if rb.maxSize > 0 && newcap > rb.maxSize {
		// resize to maxSize
		newcap = rb.maxSize
	}

	return newcap
}

func (rb *SelfAdaptiveRingBuffer) Peek() (interface{}, bool) {
	if rb.IsEmpty() {
		return nil, false
	}

	v := rb.buf[rb.r]
	return v, true
}

func (rb *SelfAdaptiveRingBuffer) Pop() (interface{}, bool) {
	v, ok := rb.Peek()
	if !ok {
		return nil, false
	}

	rb.buf[rb.r] = nil // de-reference
	rb.r++
	if rb.r == rb.size {
		// out of range
		rb.r = 0
	}

	if rb.full {
		rb.full = false
	}

	return v, true
}

func (rb *SelfAdaptiveRingBuffer) IsEmpty() bool {
	return !rb.full && rb.r == rb.w
}

func (rb *SelfAdaptiveRingBuffer) NeedReset() bool {
	return rb.IsEmpty() && rb.size > rb.initSize
}

func (rb *SelfAdaptiveRingBuffer) IsFull() bool {
	return rb.full
}

func (rb *SelfAdaptiveRingBuffer) Len() int {
	if rb.IsEmpty() {
		return 0
	}

	if rb.w > rb.r {
		return rb.w - rb.r
	}

	return rb.size - rb.r + rb.w
}

func (rb *SelfAdaptiveRingBuffer) Cap() int {
	return rb.size
}

func (rb *SelfAdaptiveRingBuffer) Reset() {
	rb.r, rb.w = 0, 0
	rb.size = rb.initSize
	rb.buf = make([]interface{}, rb.initSize)
}
