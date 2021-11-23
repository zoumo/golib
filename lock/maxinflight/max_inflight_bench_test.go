package maxinflight

import (
	"testing"
)

func BenchmarkMaxInflightLock_Channel(b *testing.B) {
	l := newChannelLock(uint32(b.N / 2))
	benchmarkMaxInFlightLock_TryAcquire(l, b)
}

func BenchmarkMaxInflightLock_Atomic(b *testing.B) {
	l := newAtomic(uint32(b.N / 2))
	benchmarkMaxInFlightLock_TryAcquire(l, b)
}

func BenchmarkMaxInflightLock_Mutex(b *testing.B) {
	l := newMutexLock(uint32(b.N / 2))
	benchmarkMaxInFlightLock_TryAcquire(l, b)
}

func benchmarkMaxInFlightLock_TryAcquire(l Lock, b *testing.B) {
	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			l.TryAcquire()
		}
	})
}
