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

package maxinflight

import (
	"sync"
	"sync/atomic"
)

var (
	InfinityTokenBucket = NewInfinity()
)

type TokenBucket interface {
	// TryAccept returns true if a token is taken immediately. Otherwise,
	// it returns false.
	TryAcquire() bool
	// Release add a token back to the lock
	Release()
	// Resize changes the max in flight lock's capacity
	Resize(n uint32)
}

type TokenBucketType string

const (
	Atomic   TokenBucketType = "atomic"
	Channel  TokenBucketType = "channel"
	Mutex    TokenBucketType = "mutex"
	Infinity TokenBucketType = "infinity"
)

func New(size uint32) TokenBucket {
	return newBucket(Atomic, size)
}

func NewInfinity() TokenBucket {
	return &infinity{}
}

func newBucket(t TokenBucketType, size uint32) TokenBucket {
	switch t {
	case Atomic:
		return newAtomic(size)
	case Channel:
		return newChannel(size)
	case Mutex:
		return newMutex(size)
	case Infinity:
		return NewInfinity()
	}
	return nil
}

// Infinity is a special lock, it always return true when
// you try to acquire one token
type infinity struct {
}

func (*infinity) TryAcquire() bool {
	return true
}

// Release add a token back to the lock
func (*infinity) Release() {
}

// Resize changes the max in flight lock's capacity
func (*infinity) Resize(n uint32) {
}

// use a larger range of values than max to avoid overflow when increacing count
type atomicTokenBucket struct {
	max   uint32 // range of 0 ~ 4,294,967,295
	count int64  // range of -9,223,372,036,854,775,808 ~ 9,223,372,036,854,775,807
}

func newAtomic(n uint32) *atomicTokenBucket {
	return &atomicTokenBucket{
		max:   n,
		count: 0,
	}
}

func (f *atomicTokenBucket) TryAcquire() bool {
	count := atomic.LoadInt64(&f.count)
	max := int64(atomic.LoadUint32(&f.max))
	if count < 0 {
		if atomic.CompareAndSwapInt64(&f.count, count, 1) {
			// reset count to 0 and add 1
			return true
		}
		// if not swapped, f.count must be set to zero or bigger
		// than one, so we just need to increace it
	} else if count >= max {
		return false
	}
	count = atomic.AddInt64(&f.count, 1)
	if count > max {
		atomic.AddInt64(&f.count, -1)
		return false
	}
	return true
}

func (f *atomicTokenBucket) Release() {
	if f.count <= 0 {
		return
	}
	count := atomic.AddInt64(&f.count, -1)
	if count < 0 {
		atomic.StoreInt64(&f.count, 0)
	}
}

func (f *atomicTokenBucket) Resize(n uint32) {
	if f.max != n {
		atomic.StoreUint32(&f.max, n)
	}
}

type channelTokenBucket struct {
	ch chan bool
}

func newChannel(n uint32) *channelTokenBucket {
	return &channelTokenBucket{
		ch: make(chan bool, n),
	}
}

func (l *channelTokenBucket) TryAcquire() bool {
	select {
	case l.ch <- true:
		return true
	default:
		return false
	}
}

func (l *channelTokenBucket) Release() {
	select {
	case <-l.ch:
	default:
	}
}

func (l *channelTokenBucket) Resize(n uint32) {
	// not implement
}

type mutexTokenBucket struct {
	count int64
	max   uint32
	m     sync.Mutex
}

func newMutex(n uint32) *mutexTokenBucket {
	return &mutexTokenBucket{
		max: n,
	}
}

func (f *mutexTokenBucket) TryAcquire() bool {
	if f.count >= int64(f.max) {
		return false
	}

	f.m.Lock()
	defer f.m.Unlock()

	if f.count >= int64(f.max) {
		return false
	}

	f.count++
	return true
}

func (f *mutexTokenBucket) Release() {
	if f.count == 0 {
		return
	}

	f.m.Lock()
	defer f.m.Unlock()

	if f.count == 0 {
		return
	}

	f.count--
}

func (f *mutexTokenBucket) Resize(n uint32) {
	if f.max == n {
		return
	}

	f.m.Lock()
	defer f.m.Unlock()

	if f.max != n {
		f.max = n
	}
}
