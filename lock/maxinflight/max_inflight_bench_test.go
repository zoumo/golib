package maxinflight

import (
	"testing"
)

func BenchmarkMaxInflight_Channel(b *testing.B) {
	l := newChannel(uint32(b.N / 2))
	benchmarkMaxInFlight_TryAcquire(l, b)
}

func BenchmarkMaxInflight_Atomic(b *testing.B) {
	l := newAtomic(uint32(b.N / 2))
	benchmarkMaxInFlight_TryAcquire(l, b)
}

func BenchmarkMaxInflight_Mutex(b *testing.B) {
	l := newMutex(uint32(b.N / 2))
	benchmarkMaxInFlight_TryAcquire(l, b)
}

func benchmarkMaxInFlight_TryAcquire(l TokenBucket, b *testing.B) {
	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			l.TryAcquire()
		}
	})
}
