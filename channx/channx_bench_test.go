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
	"testing"
)

func BenchmarkBuiltin0(b *testing.B) {
	ch := make(chan interface{})
	benchmarkChannel(b, ch, ch, func() { close(ch) })
}

func BenchmarkBuiltin10(b *testing.B) {
	ch := make(chan interface{}, 10)
	benchmarkChannel(b, ch, ch, func() { close(ch) })
}

func BenchmarkBuiltin100(b *testing.B) {
	ch := make(chan interface{}, 1000)
	benchmarkChannel(b, ch, ch, func() { close(ch) })
}

func BenchmarkBuiltin1000(b *testing.B) {
	ch := make(chan interface{}, 1000)
	benchmarkChannel(b, ch, ch, func() { close(ch) })
}

func BenchmarkBuiltin10000(b *testing.B) {
	ch := make(chan interface{}, 10000)
	benchmarkChannel(b, ch, ch, func() { close(ch) })
}

func BenchmarkPipe0(b *testing.B) {
	in := make(chan interface{})
	out := make(chan interface{})
	go func() {
		for v := range in {
			out <- v
		}
		close(out)
	}()
	benchmarkChannel(b, in, out, func() { close(in) })
}

func BenchmarkPipe10(b *testing.B) {
	in := make(chan interface{}, 10)
	out := make(chan interface{}, 10)
	go func() {
		for v := range in {
			out <- v
		}
		close(out)
	}()
	benchmarkChannel(b, in, out, func() { close(in) })
}

func BenchmarkPipe100(b *testing.B) {
	in := make(chan interface{}, 10)
	out := make(chan interface{}, 10)
	go func() {
		for v := range in {
			out <- v
		}
		close(out)
	}()
	benchmarkChannel(b, in, out, func() { close(in) })
}

func BenchmarkChanx0(b *testing.B) {
	ch := New(InChanSize(0), OutChanSzie(0), InitBufferSize(1000), MaxBufferSize(10000))
	benchmarkChannel(b, ch.In(), ch.Out(), ch.Close)
}

func BenchmarkChanx10(b *testing.B) {
	ch := New(InChanSize(5), OutChanSzie(5), InitBufferSize(1000), MaxBufferSize(10000))
	benchmarkChannel(b, ch.In(), ch.Out(), ch.Close)
}

func BenchmarkChanx1000(b *testing.B) {
	ch := New(InChanSize(500), OutChanSzie(500), InitBufferSize(1000), MaxBufferSize(10000))
	benchmarkChannel(b, ch.In(), ch.Out(), ch.Close)
}

func BenchmarkChanx2000(b *testing.B) {
	ch := New(InChanSize(500), OutChanSzie(500), InitBufferSize(10), MaxBufferSize(10000))
	benchmarkChannel(b, ch.In(), ch.Out(), ch.Close)
}

func BenchmarkChanx10000(b *testing.B) {
	ch := New(InChanSize(5000), OutChanSzie(5000), InitBufferSize(1000), MaxBufferSize(10000))
	benchmarkChannel(b, ch.In(), ch.Out(), ch.Close)
}

func benchmarkChannel(b *testing.B, in chan<- interface{}, out <-chan interface{}, closeFn func()) {
	wg := sync.WaitGroup{}
	wg.Add(2)
	count := 0
	b.ResetTimer()
	go func() {
		defer wg.Done()
		for i := 0; i < b.N; i++ {
			in <- i
		}
		closeFn()
	}()
	go func() {
		defer wg.Done()
		for range out {
			count++
		}
	}()

	wg.Wait()

	if count != b.N {
		b.Errorf("want = %v, got = %v", b.N, count)
	}
}
