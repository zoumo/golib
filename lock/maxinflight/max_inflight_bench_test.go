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
