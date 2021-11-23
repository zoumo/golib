package maxinflight

import (
	"sync"
	"sync/atomic"
)

type Lock interface {
	// TryAccept returns true if a token is taken immediately. Otherwise,
	// it returns false.
	TryAcquire() bool
	// Release add a token back to the lock
	Release()
	// Resize changes the max in flight lock's capacity
	Resize(n uint32)
}

type LockType string

const (
	Atomic  LockType = "atomic"
	Channel LockType = "channel"
	Mutex   LockType = "mutex"
)

func New(size uint32) Lock {
	return newLock(Atomic, size)
}

func newLock(t LockType, size uint32) Lock {
	switch t {
	case Atomic:
		return newAtomic(size)
	case Channel:
		return newChannelLock(size)
	case Mutex:
		return newMutexLock(size)
	}
	return nil
}

// use a larger range of values than max to avoid overflow when increacing count
type atomicLock struct {
	max   uint32 // range of 0 ~ 4,294,967,295
	count int64  // range of -9,223,372,036,854,775,808 ~ 9,223,372,036,854,775,807
}

func newAtomic(n uint32) *atomicLock {
	return &atomicLock{
		max:   n,
		count: 0,
	}
}

func (f *atomicLock) TryAcquire() bool {
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

func (f *atomicLock) Release() {
	count := atomic.AddInt64(&f.count, -1)
	if count < 0 {
		atomic.StoreInt64(&f.count, 0)
	}
}

func (f *atomicLock) Resize(n uint32) {
	if f.max != n {
		atomic.StoreUint32(&f.max, n)
	}
}

type channelLock struct {
	ch chan bool
}

func newChannelLock(n uint32) *channelLock {
	return &channelLock{
		ch: make(chan bool, n),
	}
}

func (l *channelLock) TryAcquire() bool {
	select {
	case l.ch <- true:
		return true
	default:
		return false
	}
}

func (l *channelLock) Release() {
	select {
	case <-l.ch:
	default:
	}
}

func (l *channelLock) Resize(n uint32) {
	// not implement
}

type mutexLock struct {
	count int64
	max   uint32
	m     sync.Mutex
}

func newMutexLock(n uint32) *mutexLock {
	return &mutexLock{
		max: n,
	}
}

func (f *mutexLock) TryAcquire() bool {
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

func (f *mutexLock) Release() {
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

func (f *mutexLock) Resize(n uint32) {
	if f.max == n {
		return
	}

	f.m.Lock()
	defer f.m.Unlock()

	if f.max != n {
		f.max = n
	}
}
