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
	"reflect"
	"sync"
	"testing"
	"time"
)

func TestChanX_Close(t *testing.T) {
	tests := []struct {
		name   string
		ch     *ChannX
		input  []interface{}
		output []interface{}
	}{
		{
			name: "close directly",
			ch: New(
				InChanSize(0),
				OutChanSzie(0),
				InitBufferSize(1),
				MaxBufferSize(1),
			),
		},
		{
			name: "close directly",
			ch: New(
				InChanSize(0),
				OutChanSzie(0),
				InitBufferSize(1),
				MaxBufferSize(1),
			),
			input:  []interface{}{1},
			output: []interface{}{1},
		},
		{
			name: "drop date after closed",
			ch: New(
				InChanSize(0),
				OutChanSzie(1),
				InitBufferSize(1),
				MaxBufferSize(1),
				DropClosedBufferData(),
			),
			input:  []interface{}{1, 2},
			output: []interface{}{1},
		},
	}

	for i := range tests {
		tt := tests[i]
		t.Run(tt.name, func(t *testing.T) {
			for _, v := range tt.input {
				tt.ch.In() <- v
			}
			tt.ch.Close()

			for _, want := range tt.output {
				got := <-tt.ch.Out()
				if !reflect.DeepEqual(want, got) {
					t.Errorf("want = %v but got = %v", want, got)
				}
			}
			_, ok := <-tt.ch.Out()
			if ok {
				t.Errorf("output channel should be closed")
			}
		})
	}
}

func TestChanX_Full(t *testing.T) {
	// input ->  buffer -> output
	//   1    +    1    +    1    =  3 + 1(poped)
	ch := New(
		InChanSize(1),
		OutChanSzie(1),
		InitBufferSize(1),
		MaxBufferSize(1),
	)
	for i := 0; i < 4; i++ {
		select {
		case ch.In() <- i:
		case <-time.After(1 * time.Millisecond):
			t.Errorf("should not block")
		}
	}

	select {
	case ch.In() <- 5:
		t.Errorf("should be blocked")
	default:
	}

	ch.Close()

	// wait close channel event
	time.Sleep(1 * time.Millisecond)

	for want := 0; want < 4; want++ {
		got, ok := <-ch.Out()
		if !ok {
			// closed
			return
		}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("get from output channel want = %v, got = %v", want, got)
		}
	}
}

func TestChanX_CustomerFirst(t *testing.T) {
	tests := []struct {
		name   string
		ch     *ChannX
		input  []interface{}
		output []interface{}
	}{
		{
			name: "no channel buffer",
			ch: New(
				InChanSize(0),
				OutChanSzie(0),
				InitBufferSize(1),
				MaxBufferSize(1),
			),
			input:  rangeIntSlice(0, 1000),
			output: rangeIntSlice(0, 1000),
		},
		{
			name: "no input channel buffer",
			ch: New(
				InChanSize(0),
				OutChanSzie(1),
				InitBufferSize(1),
				MaxBufferSize(1),
			),
			input:  rangeIntSlice(0, 1000),
			output: rangeIntSlice(0, 1000),
		},
		{
			name: "no output channel buffer",
			ch: New(
				InChanSize(0),
				OutChanSzie(1),
				InitBufferSize(1),
				MaxBufferSize(1),
			),
			input:  rangeIntSlice(0, 1000),
			output: rangeIntSlice(0, 1000),
		},
		{
			name: "with small channel buffer",
			ch: New(
				InChanSize(1),
				OutChanSzie(1),
				InitBufferSize(1),
				MaxBufferSize(1),
			),
			input:  rangeIntSlice(0, 1000),
			output: rangeIntSlice(0, 1000),
		},
		{
			name: "with medium channel buffer",
			ch: New(
				InChanSize(2),
				OutChanSzie(2),
				InitBufferSize(2),
				MaxBufferSize(6),
			),
			input:  rangeIntSlice(0, 1000),
			output: rangeIntSlice(0, 1000),
		},
		{
			name: "with large channel buffer",
			ch: New(
				InChanSize(10),
				OutChanSzie(10),
				InitBufferSize(10),
				MaxBufferSize(100),
			),
			input:  rangeIntSlice(0, 1000),
			output: rangeIntSlice(0, 1000),
		},
		{
			name: "with large2 channel buffer",
			ch: New(
				InChanSize(10),
				OutChanSzie(10),
				InitBufferSize(10),
				MaxBufferSize(100),
			),
			input:  rangeIntSlice(0, 1000),
			output: rangeIntSlice(0, 1000),
		},
	}

	for i := range tests {
		tt := tests[i]
		t.Run(tt.name, func(t *testing.T) {
			testSequenceScenario1CustomerFirst(t, tt.ch, tt.input, tt.output)
		})
	}
}

func TestChanX_ProducerFirst(t *testing.T) {
	tests := []struct {
		name   string
		ch     *ChannX
		input  []interface{}
		output []interface{}
	}{
		{
			name: "no channel buffer",
			ch: New(
				InChanSize(0),
				OutChanSzie(0),
				InitBufferSize(1),
				MaxBufferSize(1),
			),
			input:  rangeIntSlice(0, 1000),
			output: rangeIntSlice(0, 1000),
		},
		{
			name: "no input channel buffer",
			ch: New(
				InChanSize(0),
				OutChanSzie(1),
				InitBufferSize(1),
				MaxBufferSize(1),
			),
			input:  rangeIntSlice(0, 1000),
			output: rangeIntSlice(0, 1000),
		},
		{
			name: "no output channel buffer",
			ch: New(
				InChanSize(0),
				OutChanSzie(1),
				InitBufferSize(1),
				MaxBufferSize(1),
			),
			input:  rangeIntSlice(0, 1000),
			output: rangeIntSlice(0, 1000),
		},
		{
			name: "with small channel buffer",
			ch: New(
				InChanSize(1),
				OutChanSzie(1),
				InitBufferSize(1),
				MaxBufferSize(1),
			),
			input:  rangeIntSlice(0, 1000),
			output: rangeIntSlice(0, 1000),
		},
		{
			name: "with medium channel buffer",
			ch: New(
				InChanSize(2),
				OutChanSzie(2),
				InitBufferSize(2),
				MaxBufferSize(6),
			),
			input:  rangeIntSlice(0, 1000),
			output: rangeIntSlice(0, 1000),
		},
		{
			name: "with large channel buffer",
			ch: New(
				InChanSize(10),
				OutChanSzie(10),
				InitBufferSize(10),
				MaxBufferSize(100),
			),
			input:  rangeIntSlice(0, 1000),
			output: rangeIntSlice(0, 1000),
		},
		{
			name: "with large2 channel buffer",
			ch: New(
				InChanSize(10),
				OutChanSzie(10),
				InitBufferSize(10),
				MaxBufferSize(100),
			),
			input:  rangeIntSlice(0, 1000),
			output: rangeIntSlice(0, 1000),
		},
	}

	for i := range tests {
		tt := tests[i]
		t.Run(tt.name, func(t *testing.T) {
			testSequenceScenario1ProducerFirst(t, tt.ch, tt.input, tt.output)
		})
	}
}

// func TestChanX_Transform(t *testing.T) {
// 	tests := []struct {
// 		name   string
// 		ch     *ChannX
// 		input  []interface{}
// 		output []interface{}
// 	}{
// 		{
// 			name: "transform int to uint",
// 			ch: New(
// 				InChanSize(0),
// 				OutChanSzie(0),
// 				InitBufferSize(1),
// 				MaxBufferSize(1),
// 				WithTransform(func(i interface{}) (interface{}, error) {
// 					intI, ok := i.(int)
// 					if !ok {
// 						return nil, fmt.Errorf("want int")
// 					}
// 					return uint(intI), nil
// 				}),
// 			),
// 			input:  rangeIntSlice(0, 1000),
// 			output: rangeUintSlice(0, 1000),
// 		},
// 		{
// 			name: "transform add 1",
// 			ch: New(
// 				InChanSize(0),
// 				OutChanSzie(0),
// 				InitBufferSize(1),
// 				MaxBufferSize(1),
// 				WithTransform(func(i interface{}) (interface{}, error) {
// 					intI, ok := i.(int)
// 					if !ok {
// 						return nil, fmt.Errorf("want int")
// 					}
// 					return intI + 1, nil
// 				}),
// 			),
// 			input:  rangeIntSlice(0, 1000),
// 			output: rangeIntSlice(1, 1001),
// 		},
// 	}
// 	for i := range tests {
// 		tt := tests[i]
// 		t.Run(tt.name, func(t *testing.T) {
// 			testSequenceScenario1CustomerFirst(t, tt.ch, tt.input, tt.output)
// 		})
// 	}
// }

// func TestChanX_Filter(t *testing.T) {
// 	tests := []struct {
// 		name   string
// 		ch     *ChannX
// 		input  []interface{}
// 		output []interface{}
// 	}{
// 		{
// 			name: "filter i < 500",
// 			ch: New(
// 				InChanSize(0),
// 				OutChanSzie(0),
// 				InitBufferSize(1),
// 				MaxBufferSize(1),
// 				WithFilter(func(i interface{}) bool {
// 					intI, ok := i.(int)
// 					if !ok {
// 						return false
// 					}
// 					return intI < 500
// 				}),
// 			),
// 			input:  rangeIntSlice(0, 1000),
// 			output: rangeIntSlice(0, 500),
// 		},
// 	}
// 	for i := range tests {
// 		tt := tests[i]
// 		t.Run(tt.name, func(t *testing.T) {
// 			testSequenceScenario1CustomerFirst(t, tt.ch, tt.input, tt.output)
// 		})
// 	}
// }

func testSequenceScenario1CustomerFirst(t *testing.T, ch *ChannX, input, output []interface{}) {
	wg := sync.WaitGroup{}

	wg.Add(2)
	go func() {
		defer wg.Done()
		i := 0
		for {
			got, ok := <-ch.Out()
			if !ok {
				// closed
				return
			}
			want := output[i]
			if !reflect.DeepEqual(got, want) {
				t.Errorf("get from output channel want = %v, got = %v", want, got)
			}
			i++
		}
	}()

	go func() {
		defer wg.Done()
		for i := range input {
			ch.In() <- input[i]
		}
		ch.Close()
	}()

	wg.Wait()
}

func testSequenceScenario1ProducerFirst(t *testing.T, ch *ChannX, input, output []interface{}) {
	wg := sync.WaitGroup{}

	wg.Add(2)

	go func() {
		defer wg.Done()
		for i := range input {
			ch.In() <- input[i]
		}
		ch.Close()
	}()

	go func() {
		defer wg.Done()
		i := 0
		for {
			got, ok := <-ch.Out()
			if !ok {
				// closed
				return
			}
			want := output[i]
			if !reflect.DeepEqual(got, want) {
				t.Errorf("get from output channel want = %v, got = %v", want, got)
			}
			i++
		}
	}()

	wg.Wait()
}

func rangeIntSlice(start, end int) []interface{} {
	ret := []interface{}{}
	for i := start; i < end; i++ {
		ret = append(ret, i)
	}
	return ret
}
