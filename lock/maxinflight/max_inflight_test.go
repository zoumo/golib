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
	"math"
	"sync"
	"sync/atomic"
	"testing"
)

func TestMaxInflight_TryAcquire(t *testing.T) {
	tests := []struct {
		name         string
		lock         TokenBucketType
		max          uint32
		acquireTimes uint32
		want         uint32
	}{
		{
			"",
			Infinity,
			1000000,
			2000000,
			2000000,
		},
		{
			"",
			Atomic,
			0,
			2000000,
			0,
		},
		{
			"",
			Atomic,
			1000000,
			2000000,
			1000000,
		},
		{
			"",
			Atomic,
			0,
			2000000,
			0,
		},
		{
			"",
			Channel,
			0,
			20000,
			0,
		},
		{
			"",
			Channel,
			10000,
			20000,
			10000,
		},
		{
			"",
			Mutex,
			0,
			20000,
			0,
		},
		{
			"",
			Mutex,
			10000,
			20000,
			10000,
		},
	}
	for i := range tests {
		tt := tests[i]
		t.Run(tt.name, func(t *testing.T) {
			l := newBucket(tt.lock, tt.max)
			g := sync.WaitGroup{}

			g.Add(int(tt.acquireTimes))
			var got uint32
			for i := 0; i < int(tt.acquireTimes); i++ {
				go func() {
					defer g.Done()
					if l.TryAcquire() {
						atomic.AddUint32(&got, 1)
					}
				}()
			}
			g.Wait()
			if got != tt.want {
				t.Errorf("MaxInFlightLock.TryAcquire() type=%v, got=%v, want %v", tt.lock, got, tt.want)
			}
		})
	}
}

func Test_atomicLock_TryAcquire(t *testing.T) {
	want := uint32(10)
	lock := &atomicTokenBucket{
		count: int64(math.MaxUint32 - want),
		max:   math.MaxUint32,
	}

	var success uint32

	for i := 0; i < 100000; i++ {
		go func() {
			if lock.TryAcquire() {
				atomic.AddUint32(&success, 1)
			}
		}()
	}

	if success != want {
		t.Errorf("atomicMaxInFlightLock.TryAcquire() = %v, want %v", success, want)
	}
}
