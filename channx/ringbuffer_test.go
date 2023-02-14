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
	"reflect"
	"testing"
)

func TestSelfAdaptiveRingBuffer_growCap(t *testing.T) {
	type fields struct {
		buf      []interface{}
		maxSize  int
		initSize int
		size     int
		r        int
		w        int
		full     bool
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		{
			name: "[infinity] double",
			fields: fields{
				maxSize: 0,
				size:    200,
			},
			want: 400,
		},
		{
			name: "[infinity] smooth-ish transition to 1.25x",
			fields: fields{
				maxSize: 0,
				size:    2000,
			},
			want: 2000 + (2000+growThreshold*3)/4,
		},
		{
			name: "[infinity] overflow 1",
			fields: fields{
				maxSize: 0,
				size:    math.MaxInt - 10,
			},
			want: math.MaxInt,
		},
		{
			name: "[infinity] overflow the add number is not overflow but result is",
			fields: fields{
				maxSize: 0,
				size:    (math.MaxInt-growThreshold*3/4)*4/5 + 10,
			},
			want: math.MaxInt,
		},
		{
			name: "[Bounded] double but resize",
			fields: fields{
				maxSize: 5,
				size:    4,
			},
			want: 5,
		},
		{
			name: "[Bounded] double but resize",
			fields: fields{
				maxSize: 300,
				size:    200,
			},
			want: 300,
		},
		{
			name: "[Bounded] 1.25x but resize",
			fields: fields{
				maxSize: 2100,
				size:    2000,
			},
			want: 2100,
		},
		{
			name: "[Bounded] overflow but resize",
			fields: fields{
				maxSize: math.MaxInt - 5,
				size:    math.MaxInt - 10,
			},
			want: math.MaxInt - 5,
		},

		{
			name: "[Bounded] already over max",
			fields: fields{
				maxSize: 10,
				size:    10000,
			},
			want: 10000,
		},
	}
	for i := range tests {
		tt := tests[i]
		t.Run(tt.name, func(t *testing.T) {
			rb := &SelfAdaptiveRingBuffer{
				buf:      tt.fields.buf,
				maxSize:  tt.fields.maxSize,
				initSize: tt.fields.initSize,
				size:     tt.fields.size,
				r:        tt.fields.r,
				w:        tt.fields.w,
				full:     tt.fields.full,
			}
			if got := rb.growCap(); got != tt.want {
				t.Errorf("SelfAdaptiveRingBuffer.growCap() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSelfAdaptiveRingBuffer_IsEmpty(t *testing.T) {
	type fields struct {
		r    int
		w    int
		full bool
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "empty 0",
			fields: fields{
				r: 0,
				w: 0,
			},
			want: true,
		},
		{
			name: "empty 1",
			fields: fields{
				r: 1,
				w: 1,
			},
			want: true,
		},
		{
			name: "full 1",
			fields: fields{
				r:    1,
				w:    1,
				full: true,
			},
			want: false,
		},
		{
			name: "not empty",
			fields: fields{
				r: 2,
				w: 1,
			},
			want: false,
		},
	}
	for i := range tests {
		tt := tests[i]
		t.Run(tt.name, func(t *testing.T) {
			rb := &SelfAdaptiveRingBuffer{
				r:    tt.fields.r,
				w:    tt.fields.w,
				full: tt.fields.full,
			}
			if got := rb.IsEmpty(); got != tt.want {
				t.Errorf("SelfAdaptiveRingBuffer.IsEmpty() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSelfAdaptiveRingBuffer_LenAndCap(t *testing.T) {
	type fields struct {
		size int
		r    int
		w    int
		full bool
	}
	tests := []struct {
		name    string
		fields  fields
		wantLen int
		wantCap int
	}{
		{
			name: "empty",
			fields: fields{
				r:    0,
				w:    0,
				full: false,
				size: 10,
			},
			wantLen: 0,
			wantCap: 10,
		},
		{
			name: "full",
			fields: fields{
				size: 10,
				r:    4,
				w:    4,
				full: true,
			},
			wantLen: 10,
			wantCap: 10,
		},
		{
			name: "w > r",
			fields: fields{
				size: 10,
				r:    1,
				w:    6,
				full: false,
			},
			wantLen: 5,
			wantCap: 10,
		},
		{
			name: "w < r",
			fields: fields{
				size: 10,
				r:    6,
				w:    2,
				full: false,
			},
			wantLen: 6,
			wantCap: 10,
		},
	}
	for i := range tests {
		tt := tests[i]
		t.Run(tt.name, func(t *testing.T) {
			rb := &SelfAdaptiveRingBuffer{
				size: tt.fields.size,
				r:    tt.fields.r,
				w:    tt.fields.w,
				full: tt.fields.full,
			}
			if got := rb.Len(); got != tt.wantLen {
				t.Errorf("SelfAdaptiveRingBuffer.Len() = %v, want %v", got, tt.wantLen)
			}
			if got := rb.Cap(); got != tt.wantCap {
				t.Errorf("SelfAdaptiveRingBuffer.Cap() = %v, want %v", got, tt.wantLen)
			}
		})
	}
}

func TestSelfAdaptiveRingBuffer_Put(t *testing.T) {
	rb := NewSelfAdptiveRingBuffer(2, 5)

	noErrInput := []interface{}{1, 2, 3, 4, 5}
	for _, v := range noErrInput {
		if !rb.Put(v) {
			t.Errorf("SelfAdaptiveRingBuffer.Put() should be successful")
		}
	}
	wantErrInput := []interface{}{6, 7}
	for _, v := range wantErrInput {
		if rb.Put(v) {
			t.Errorf("SelfAdaptiveRingBuffer.Put should not failed")
		}
	}

	if rb.Len() != 5 {
		t.Error("SelfAdaptiveRingBuffer.Len() should be 5")
	}
	if rb.w != 0 {
		t.Error("SelfAdaptiveRingBuffer write index should be 0")
	}
	if !rb.IsFull() {
		t.Error("SelfAdaptiveRingBuffer should be full")
	}
	if !reflect.DeepEqual(rb.buf, noErrInput) {
		t.Errorf("SelfAdaptiveRingBuffer.buf not match, want = %v, got= %v", noErrInput, rb.buf)
	}

	rb = NewSelfAdptiveRingBuffer(2, 0)
	noErrInput = []interface{}{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	for _, v := range noErrInput {
		if !rb.Put(v) {
			t.Errorf("SelfAdaptiveRingBuffer.Put() should be successful")
		}
	}
	if rb.Len() != 10 {
		t.Error("SelfAdaptiveRingBuffer.Len() should be 10")
	}
	if rb.w != 10 {
		t.Error("SelfAdaptiveRingBuffer write index should be 10")
	}
	if rb.IsFull() {
		t.Error("SelfAdaptiveRingBuffer should not be full")
	}
}

func TestSelfAdaptiveRingBuffer_Pop(t *testing.T) {
	rb := NewSelfAdptiveRingBuffer(2, 5)
	_, ok := rb.Pop()
	if ok {
		t.Errorf("SelfAdaptiveRingBuffer.Pop() want false when buffer is empty")
	}

	input := []interface{}{1, 2, 3, 4, 5}
	for _, v := range input {
		rb.Put(v)
	}

	for i := range input {
		got, _ := rb.Pop()
		want := input[i]
		if !reflect.DeepEqual(got, want) {
			t.Errorf("SelfAdaptiveRingBuffer.Pop(), want = %v, got = %v", want, got)
		}
	}

	if rb.IsFull() {
		t.Errorf("ring buffer must bot be full after pop")
	}

	if !rb.IsEmpty() {
		t.Errorf("ring buffer must be empty")
	}
}
