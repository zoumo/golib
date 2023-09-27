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

package heap

import (
	"container/heap"
	"fmt"
)

type KeyError struct {
	Obj interface{}
	Err error
}

// Error gives a human-readable description of the error.
func (k KeyError) Error() string {
	return fmt.Sprintf("couldn't create key for object %+v: %v", k.Obj, k.Err)
}

// KeyFunc is a function type to get the key from an object.
type KeyFunc func(obj interface{}) (string, error)

// LessFunc is a function that receives two items and returns true if the first
// item should be placed before the second one when the list is sorted.
type LessFunc = func(x, y interface{}) bool

type containerHeapItem struct {
	// The index of the object's key in the Heap.queue.
	index int
	// The key of the object
	key string
	// The object which is stored in the heap.
	obj interface{}
}

// containerHeap is an struct that implements the standard container/heap interface
// and keeps the containerHeap stored in the heap.
type containerHeap struct {
	// items is a map from key of the objects to the objects and their index.
	// We depend on the property that items in the map are in the queue and vice versa.
	items map[string]*containerHeapItem
	// ordered implements a heap data structure and keeps the order of elements
	// according to the heap invariant. The ordered keeps the keys of objects stored
	// in "items".
	ordered []string
	// lessFunc is used to compare two objects in the heap.
	lessFunc LessFunc
}

var (
	_ = heap.Interface(&containerHeap{}) // heapData is a standard heap
)

// Len returns the number of items in the Heap.
// Implement standard heap.Interface.
func (h *containerHeap) Len() int {
	return len(h.ordered)
}

// Less reports whether the element with
// index i should sort before the element with index j.
// Implement standard heap.Interface.
func (h *containerHeap) Less(i, j int) bool {
	if i > len(h.ordered) || j > len(h.ordered) {
		return false
	}
	x, ok := h.items[h.ordered[i]]
	if !ok {
		return false
	}
	y, ok := h.items[h.ordered[j]]
	if !ok {
		return false
	}
	return h.lessFunc(x.obj, y.obj)
}

// Swap swaps the elements with indexes i and j.
// Implement standard heap.Interface
func (h *containerHeap) Swap(i, j int) {
	h.ordered[i], h.ordered[j] = h.ordered[j], h.ordered[i]
	item := h.items[h.ordered[i]]
	item.index = i
	item = h.items[h.ordered[j]]
	item.index = j
}

// Push pushes the element kv onto the heap.
// Implement standard heap.Interface.
func (h *containerHeap) Push(kv interface{}) {
	item := kv.(*containerHeapItem)
	item.index = len(h.ordered)
	h.items[item.key] = item
	h.ordered = append(h.ordered, item.key)
}

// Pop removes and returns the minimum element (according to Less) from the heap,
// Implement standard heap.Interface.
func (h *containerHeap) Pop() interface{} {
	if len(h.ordered) == 0 {
		return nil
	}
	n := len(h.ordered)
	key := h.ordered[n-1]
	item := h.items[key]

	h.ordered = h.ordered[0 : n-1]
	delete(h.items, key)
	return item.obj
}

// Peek returns the head of the heap without removing it.
func (h *containerHeap) Peek() interface{} {
	if len(h.ordered) > 0 {
		return h.items[h.ordered[0]].obj
	}
	return nil
}

// PeekSecond returns the second item of heap without removing it.
func (h *containerHeap) PeekSecond() interface{} {
	if len(h.ordered) < 2 {
		return nil
	}
	if len(h.ordered) == 2 {
		return h.items[h.ordered[1]].obj
	}
	// compare left and right child
	// index 0 is head, 1 is the left child and 2 is the right.
	if h.Less(1, 2) {
		return h.items[h.ordered[1]].obj
	}
	return h.items[h.ordered[2]].obj
}

func (h *containerHeap) GetByKey(key string) (interface{}, bool) {
	item, ok := h.items[key]
	if !ok {
		return nil, false
	}
	return item.obj, true
}

// Heap is a producer/consumer queue that implements a heap data structure.
// It can be used to implement priority queues and similar data structures.
type Heap struct {
	// data stores objects and has a queue that keeps their ordering according
	// to the heap invariant.
	data *containerHeap
	// keyFunc is used to make the key used for queued item insertion and retrieval, and
	// should be deterministic.
	keyFunc KeyFunc
}

func New(keyfunc KeyFunc, lessfunc LessFunc) *Heap {
	return &Heap{
		data: &containerHeap{
			items:    make(map[string]*containerHeapItem),
			ordered:  make([]string, 0),
			lessFunc: lessfunc,
		},
		keyFunc: keyfunc,
	}
}

func (h *Heap) Len() int {
	return h.data.Len()
}

// Add inserts an item, and puts it in the queue. The item is updated if it
// already exists.
func (h *Heap) AddOrUpdate(obj interface{}) error {
	key, err := h.keyFunc(obj)
	if err != nil {
		return KeyError{Obj: obj, Err: err}
	}
	if _, exists := h.data.items[key]; exists {
		h.data.items[key].obj = obj
		heap.Fix(h.data, h.data.items[key].index)
	} else {
		heap.Push(h.data, &containerHeapItem{key: key, obj: obj})
	}
	return nil
}

// AddIfNotPresent inserts an item, and puts it in the queue. If an item with
// the key is present in the heap, no changes is made to the item.
func (h *Heap) AddIfNotPresent(obj interface{}) error {
	key, err := h.keyFunc(obj)
	if err != nil {
		return KeyError{Obj: obj, Err: err}
	}
	if _, exists := h.data.items[key]; !exists {
		heap.Push(h.data, &containerHeapItem{key: key, obj: obj})
	}
	return nil
}

// UpdateIfPresent update an item's obj and fix the order if it is present in the heap.
func (h *Heap) UpdateIfPresent(obj interface{}) error {
	key, err := h.keyFunc(obj)
	if err != nil {
		return KeyError{Obj: obj, Err: err}
	}
	if _, exists := h.data.items[key]; exists {
		h.data.items[key].obj = obj
		heap.Fix(h.data, h.data.items[key].index)
	}
	return nil
}

// Delete removes an item.
func (h *Heap) Remove(obj interface{}) error {
	key, err := h.keyFunc(obj)
	if err != nil {
		return KeyError{Obj: obj, Err: err}
	}
	if item, ok := h.data.items[key]; ok {
		heap.Remove(h.data, item.index)
		return nil
	}
	return nil
}

// Pop returns the head of the heap and removes it.
func (h *Heap) Pop() interface{} {
	if len(h.data.ordered) == 0 {
		return nil
	}
	return heap.Pop(h.data)
}

// Peek returns the head of the heap without removing it.
func (h *Heap) Peek() interface{} {
	return h.data.Peek()
}

// PeekSecond returns the second item of heap without removing it.
func (h *Heap) PeekSecond() interface{} {
	return h.data.PeekSecond()
}

// GetByKey returns the requested item, or sets exists=false.
func (h *Heap) GetByKey(key string) (interface{}, bool) {
	return h.data.GetByKey(key)
}

// Range calls f sequentially for each key and value present in the heap.
// If f returns false, range stops the iteration.
//
// Range does not guarantee the order.
func (h *Heap) Range(f func(i int, key string, obj interface{}) bool) {
	for _, item := range h.data.items {
		f(item.index, item.key, item.obj)
	}
}

// List returns a list of all the items.
func (h *Heap) List() []interface{} {
	if len(h.data.items) == 0 {
		return []interface{}{}
	}
	list := make([]interface{}, len(h.data.items))
	i := 0
	for _, item := range h.data.items {
		list[i] = item.obj
		i++
	}
	return list
}
