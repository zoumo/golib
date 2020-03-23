package heap

import (
	"testing"
)

// This file was copied from k8s.io/client-go/tools/cache/heap.go and modified
// for our heap

func testHeapObjectKeyFunc(obj interface{}) (string, error) {
	return obj.(testHeapObject).name, nil
}

type testHeapObject struct {
	name string
	val  interface{}
}

func mkHeapObj(name string, val interface{}) testHeapObject {
	return testHeapObject{name: name, val: val}
}

func compareInts(val1 interface{}, val2 interface{}) bool {
	first := val1.(testHeapObject).val.(int)
	second := val2.(testHeapObject).val.(int)
	return first < second
}

// TestHeapBasic tests Heap invariant
func TestHeapBasic(t *testing.T) {
	h := New(testHeapObjectKeyFunc, compareInts)
	const amount = 500
	var i int

	nothing := h.Pop()
	if nothing != nil {
		t.Errorf("unexpected item %v", nothing)
	}
	for i = amount; i > 0; i-- {
		h.AddIfNotPresent(mkHeapObj(string([]rune{'a', rune(i)}), i))
	}

	// Make sure that the numbers are popped in ascending order.
	prevNum := 0
	for i := 0; i < amount; i++ {
		obj := h.Pop()
		if obj == nil {
			break
		}
		num := obj.(testHeapObject).val.(int)
		// All the items must be sorted.
		if prevNum > num {
			t.Errorf("got %v out of order, last was %v", obj, prevNum)
		}
		prevNum = num
	}
}

// Tests Heap.AddOrUpdate and ensures that heap invariant is preserved after adding items.
func TestHeap_AddOrUpdate(t *testing.T) {
	h := New(testHeapObjectKeyFunc, compareInts)
	h.AddOrUpdate(mkHeapObj("foo", 10))
	h.AddOrUpdate(mkHeapObj("bar", 1))
	h.AddOrUpdate(mkHeapObj("baz", 11))
	h.AddOrUpdate(mkHeapObj("zab", 30))
	h.AddOrUpdate(mkHeapObj("foo", 13)) // This updates "foo".

	if val := h.data.items["foo"].obj.(testHeapObject).val; val != 13 {
		t.Errorf("unexpected value: %d", val)
	}
	item := h.Pop()
	if e, a := 1, item.(testHeapObject).val; a != e {
		t.Fatalf("expected %d, got %d", e, a)
	}
	item = h.Pop()
	if e, a := 11, item.(testHeapObject).val; a != e {
		t.Fatalf("expected %d, got %d", e, a)
	}
	h.Remove(mkHeapObj("baz", 11))      // Nothing is deleted.
	h.AddOrUpdate(mkHeapObj("foo", 14)) // foo is updated.
	item = h.Pop()
	if e, a := 14, item.(testHeapObject).val; a != e {
		t.Fatalf("expected %d, got %d", e, a)
	}
	item = h.Pop()
	if e, a := 30, item.(testHeapObject).val; a != e {
		t.Fatalf("expected %d, got %d", e, a)
	}
}

// TestHeap_AddIfNotPresent tests Heap.AddIfNotPresent and ensures that heap
// invariant is preserved after adding items.
func TestHeap_AddIfNotPresent(t *testing.T) {
	h := New(testHeapObjectKeyFunc, compareInts)
	h.AddIfNotPresent(mkHeapObj("foo", 10))
	h.AddIfNotPresent(mkHeapObj("bar", 1))
	h.AddIfNotPresent(mkHeapObj("baz", 11))
	h.AddIfNotPresent(mkHeapObj("zab", 30))
	h.AddIfNotPresent(mkHeapObj("foo", 13)) // This is not added.

	if len := len(h.data.items); len != 4 {
		t.Errorf("unexpected number of items: %d", len)
	}
	if val := h.data.items["foo"].obj.(testHeapObject).val; val != 10 {
		t.Errorf("unexpected value: %d", val)
	}
	item := h.Pop()
	if e, a := 1, item.(testHeapObject).val; a != e {
		t.Fatalf("expected %d, got %d", e, a)
	}
	item = h.Pop()
	if e, a := 10, item.(testHeapObject).val; a != e {
		t.Fatalf("expected %d, got %d", e, a)
	}
	// bar is already popped. Let's add another one.
	h.AddIfNotPresent(mkHeapObj("bar", 14))
	item = h.Pop()
	if e, a := 11, item.(testHeapObject).val; a != e {
		t.Fatalf("expected %d, got %d", e, a)
	}
	item = h.Pop()
	if e, a := 14, item.(testHeapObject).val; a != e {
		t.Fatalf("expected %d, got %d", e, a)
	}
}

// TestHeap_AddIfNotPresent tests Heap.AddIfNotPresent and ensures that heap
// invariant is preserved after adding items.
func TestHeap_UpdateIfPresent(t *testing.T) {
	h := New(testHeapObjectKeyFunc, compareInts)
	h.AddIfNotPresent(mkHeapObj("foo", 10))
	h.AddIfNotPresent(mkHeapObj("bar", 1))
	h.AddIfNotPresent(mkHeapObj("baz", 11))
	h.AddIfNotPresent(mkHeapObj("zab", 30))

	h.UpdateIfPresent(mkHeapObj("foo", 13)) // updated
	h.UpdateIfPresent(mkHeapObj("no", 13))  // This is not updated.

	if len := len(h.data.items); len != 4 {
		t.Errorf("unexpected number of items: %d", len)
	}
	if val := h.data.items["foo"].obj.(testHeapObject).val; val != 13 {
		t.Errorf("unexpected value: %d", val)
	}
}

// TestHeap_Delete tests Heap.Delete and ensures that heap invariant is
// preserved after deleting items.
func TestHeap_Delete(t *testing.T) {
	h := New(testHeapObjectKeyFunc, compareInts)
	h.AddIfNotPresent(mkHeapObj("foo", 10))
	h.AddIfNotPresent(mkHeapObj("bar", 1))
	h.AddIfNotPresent(mkHeapObj("bal", 31))
	h.AddIfNotPresent(mkHeapObj("baz", 11))

	// Remove head. Remove should work with "key" and doesn't care about the value.
	if err := h.Remove(mkHeapObj("bar", 200)); err != nil {
		t.Fatalf("Failed to delete head.")
	}
	item := h.Pop()
	if e, a := 10, item.(testHeapObject).val; a != e {
		t.Fatalf("expected %d, got %d", e, a)
	}
	h.AddIfNotPresent(mkHeapObj("zab", 30))
	h.AddIfNotPresent(mkHeapObj("faz", 30))
	len := h.data.Len()
	var err error
	// Delete non-existing item.
	if err = h.Remove(mkHeapObj("non-existent", 10)); err != nil || len != h.data.Len() {
		t.Fatalf("Didn't expect any item removal")
	}
	// Delete tail.
	if err = h.Remove(mkHeapObj("bal", 31)); err != nil {
		t.Fatalf("Failed to delete tail.")
	}
	// Delete one of the items with value 30.
	if err = h.Remove(mkHeapObj("zab", 30)); err != nil {
		t.Fatalf("Failed to delete item.")
	}
	item = h.Pop()
	if e, a := 11, item.(testHeapObject).val; a != e {
		t.Fatalf("expected %d, got %d", e, a)
	}
	item = h.Pop()
	if e, a := 30, item.(testHeapObject).val; a != e {
		t.Fatalf("expected %d, got %d", e, a)
	}
	if h.data.Len() != 0 {
		t.Fatalf("expected an empty heap.")
	}
}

// // TestHeap_Get tests Heap.Get.
// func TestHeap_Get(t *testing.T) {
// 	h := New(testHeapObjectKeyFunc, compareInts)
// 	h.Add(mkHeapObj("foo", 10))
// 	h.Add(mkHeapObj("bar", 1))
// 	h.Add(mkHeapObj("bal", 31))
// 	h.Add(mkHeapObj("baz", 11))

// 	// Get works with the key.
// 	obj, exists, err := h.Get(mkHeapObj("baz", 0))
// 	if exists == false || obj.(testHeapObject).val != 11 {
// 		t.Fatalf("unexpected error in getting element")
// 	}
// 	// Get non-existing object.
// 	_, exists, err = h.Get(mkHeapObj("non-existing", 0))
// 	if exists == true {
// 		t.Fatalf("didn't expect to get any object")
// 	}
// }

// TestHeap_GetByKey tests Heap.GetByKey and is very similar to TestHeap_Get.
func TestHeap_GetByKey(t *testing.T) {
	h := New(testHeapObjectKeyFunc, compareInts)
	h.AddIfNotPresent(mkHeapObj("foo", 10))
	h.AddIfNotPresent(mkHeapObj("bar", 1))
	h.AddIfNotPresent(mkHeapObj("bal", 31))
	h.AddIfNotPresent(mkHeapObj("baz", 11))

	obj, exists := h.GetByKey("baz")
	if exists == false || obj.(testHeapObject).val != 11 {
		t.Fatalf("unexpected error in getting element")
	}
	// Get non-existing object.
	_, exists = h.GetByKey("non-existing")
	if exists == true {
		t.Fatalf("didn't expect to get any object")
	}
}

// TestHeap_List tests Heap.List function.
func TestHeap_List(t *testing.T) {
	h := New(testHeapObjectKeyFunc, compareInts)
	list := h.List()
	if len(list) != 0 {
		t.Errorf("expected an empty list")
	}

	items := map[string]int{
		"foo": 10,
		"bar": 1,
		"bal": 30,
		"baz": 11,
		"faz": 30,
	}
	for k, v := range items {
		h.AddIfNotPresent(mkHeapObj(k, v))
	}
	list = h.List()
	if len(list) != len(items) {
		t.Errorf("expected %d items, got %d", len(items), len(list))
	}
	for _, obj := range list {
		heapObj := obj.(testHeapObject)
		v, ok := items[heapObj.name]
		if !ok || v != heapObj.val {
			t.Errorf("unexpected item in the list: %v", heapObj)
		}
	}
}

func TestHeap_Peek(t *testing.T) {
	h := New(testHeapObjectKeyFunc, compareInts)
	head := h.Peek()
	if head != nil {
		t.Errorf("uexpected head %v", head)
	}
	items := map[string]int{
		"a": 0,
		"b": 2,
		"c": 1,
		"d": 6,
		"e": 3,
		"f": 5,
	}
	for k, v := range items {
		h.AddIfNotPresent(mkHeapObj(k, v))
	}
	head = h.Peek()
	if e, a := 0, head.(testHeapObject).val.(int); a != e {
		t.Fatalf("expected %d, got %d", e, a)
	}
}

func TestHeap_PeekSecond(t *testing.T) {
	h := New(testHeapObjectKeyFunc, compareInts)
	second := h.PeekSecond()
	if second != nil {
		t.Errorf("uexpected second %v", second)
	}

	h.AddIfNotPresent(mkHeapObj("a", 0))
	second = h.PeekSecond()
	if second != nil {
		t.Errorf("uexpected second %v", second)
	}
	h.AddIfNotPresent(mkHeapObj("b", 2))
	second = h.PeekSecond()
	if e, a := 2, second.(testHeapObject).val.(int); a != e {
		t.Fatalf("expected %d, got %d", e, a)
	}
	h.AddIfNotPresent(mkHeapObj("c", 3))
	second = h.PeekSecond()
	if e, a := 2, second.(testHeapObject).val.(int); a != e {
		t.Fatalf("expected %d, got %d", e, a)
	}
	h.AddOrUpdate(mkHeapObj("c", 1))
	second = h.PeekSecond()
	if e, a := 1, second.(testHeapObject).val.(int); a != e {
		t.Fatalf("expected %d, got %d", e, a)
	}

}
