// +build !ignore_autogenerated

// Code generated by controller-gen. DO NOT EDIT.

package testdata

import "reflect"

type MyBool bool

type MyInt int

type MyInt8 int8

type MyInt16 int16

type MyInt32 int32

type MyInt64 int64

type MyUint uint

type MyUint8 uint8

type MyUint16 uint16

type MyUint32 uint32

type MyUint64 uint64

type MyUintptr uintptr

type MyFloat32 float32

type MyFloat64 float64

type MyComplex64 complex64

type MyComplex128 complex128

type MyString string

type MySlice []string

type MyMap map[string]string

type Struct struct {
	Struct                   Predeclared
	StructPtr                *Predeclared
	UnexportedFieldStruct    UnexportedFieldStruct
	UnexportedFieldStructPtr *UnexportedFieldStruct
	Anonymous                struct {
		String string
	}
	AnonymousPtr *struct {
		String string
	}
}

func NewStruct() Struct {
	p := NewPredeclared()
	uep := NewUnexportedFieldStruct()
	return Struct{
		Anonymous: struct {
			String string
		}{String: "string"},
		AnonymousPtr: &struct {
			String string
		}{String: "string"},
		Struct:                   p,
		StructPtr:                &p,
		UnexportedFieldStruct:    uep,
		UnexportedFieldStructPtr: &uep,
	}
}

func (p Struct) Values() []reflect.Value {
	var values []reflect.Value
	values = append(values, reflect.ValueOf(p))
	return values
}

type UnexportedFieldStruct struct {
	Predeclared
	privateBool       bool
	privateInt        int
	privateInt8       int8
	privateInt16      int16
	privateInt32      int32
	privateInt64      int64
	privateUint       uint
	privateUint8      uint8
	privateUint16     uint16
	privateUint32     uint32
	privateUint64     uint64
	privateUintptr    uintptr
	privateFloat32    float32
	privateFloat64    float64
	privateComplex64  complex64
	privateComplex128 complex128
	privateString     string
}

func NewUnexportedFieldStruct() UnexportedFieldStruct {
	return UnexportedFieldStruct{
		Predeclared:       NewPredeclared(),
		privateBool:       true,
		privateComplex128: complex128((0 + 1i)),
		privateComplex64:  complex64((0 + 1i)),
		privateFloat32:    float32(3.2),
		privateFloat64:    float64(3.2),
		privateInt:        int(1),
		privateInt16:      int16(1),
		privateInt32:      int32(1),
		privateInt64:      int64(1),
		privateInt8:       int8(1),
		privateString:     "Predeclared",
		privateUint:       uint(1),
		privateUint16:     uint16(1),
		privateUint32:     uint32(1),
		privateUint64:     uint64(1),
		privateUint8:      uint8(1),
		privateUintptr:    uintptr(1),
	}
}

func (p UnexportedFieldStruct) Values() []reflect.Value {
	var values []reflect.Value
	values = append(values, reflect.ValueOf(p))
	return values
}

type Predeclared struct {
	Bool         bool
	MyBool       MyBool
	Int          int
	MyInt        MyInt
	Int8         int8
	MyInt8       MyInt8
	Int16        int16
	MyInt16      MyInt16
	Int32        int32
	MyInt32      MyInt32
	Int64        int64
	MyInt64      MyInt64
	Uint         uint
	MyUint       MyUint
	Uint8        uint8
	MyUint8      MyUint8
	Uint16       uint16
	MyUint16     MyUint16
	Uint32       uint32
	MyUint32     MyUint32
	Uint64       uint64
	MyUint64     MyUint64
	Uintptr      uintptr
	MyUintptr    MyUintptr
	Float32      float32
	MyFloat32    MyFloat32
	Float64      float64
	MyFloat64    MyFloat64
	Complex64    complex64
	MyComplex64  MyComplex64
	Complex128   complex128
	MyComplex128 MyComplex128
	String       string
	MyString     MyString
}

func NewPredeclared() Predeclared {
	return Predeclared{
		Bool:         true,
		Complex128:   complex128((0 + 1i)),
		Complex64:    complex64((0 + 1i)),
		Float32:      float32(3.2),
		Float64:      float64(3.2),
		Int:          int(1),
		Int16:        int16(1),
		Int32:        int32(1),
		Int64:        int64(1),
		Int8:         int8(1),
		MyBool:       MyBool(true),
		MyComplex128: MyComplex128(complex128((0 + 1i))),
		MyComplex64:  MyComplex64(complex64((0 + 1i))),
		MyFloat32:    MyFloat32(float32(3.2)),
		MyFloat64:    MyFloat64(float64(3.2)),
		MyInt:        MyInt(int(1)),
		MyInt16:      MyInt16(int16(1)),
		MyInt32:      MyInt32(int32(1)),
		MyInt64:      MyInt64(int64(1)),
		MyInt8:       MyInt8(int8(1)),
		MyString:     MyString("Predeclared"),
		MyUint:       MyUint(uint(1)),
		MyUint16:     MyUint16(uint16(1)),
		MyUint32:     MyUint32(uint32(1)),
		MyUint64:     MyUint64(uint64(1)),
		MyUint8:      MyUint8(uint8(1)),
		MyUintptr:    MyUintptr(uintptr(1)),
		String:       "Predeclared",
		Uint:         uint(1),
		Uint16:       uint16(1),
		Uint32:       uint32(1),
		Uint64:       uint64(1),
		Uint8:        uint8(1),
		Uintptr:      uintptr(1),
	}
}

func (p Predeclared) Values() []reflect.Value {
	var values []reflect.Value
	values = append(values, reflect.ValueOf(p.Bool), reflect.ValueOf(p.MyBool))
	values = append(values, reflect.ValueOf(p.Int), reflect.ValueOf(p.MyInt))
	values = append(values, reflect.ValueOf(p.Int8), reflect.ValueOf(p.MyInt8))
	values = append(values, reflect.ValueOf(p.Int16), reflect.ValueOf(p.MyInt16))
	values = append(values, reflect.ValueOf(p.Int32), reflect.ValueOf(p.MyInt32))
	values = append(values, reflect.ValueOf(p.Int64), reflect.ValueOf(p.MyInt64))
	values = append(values, reflect.ValueOf(p.Uint), reflect.ValueOf(p.MyUint))
	values = append(values, reflect.ValueOf(p.Uint8), reflect.ValueOf(p.MyUint8))
	values = append(values, reflect.ValueOf(p.Uint16), reflect.ValueOf(p.MyUint16))
	values = append(values, reflect.ValueOf(p.Uint32), reflect.ValueOf(p.MyUint32))
	values = append(values, reflect.ValueOf(p.Uint64), reflect.ValueOf(p.MyUint64))
	values = append(values, reflect.ValueOf(p.Uintptr), reflect.ValueOf(p.MyUintptr))
	values = append(values, reflect.ValueOf(p.Float32), reflect.ValueOf(p.MyFloat32))
	values = append(values, reflect.ValueOf(p.Float64), reflect.ValueOf(p.MyFloat64))
	values = append(values, reflect.ValueOf(p.Complex64), reflect.ValueOf(p.MyComplex64))
	values = append(values, reflect.ValueOf(p.Complex128), reflect.ValueOf(p.MyComplex128))
	values = append(values, reflect.ValueOf(p.String), reflect.ValueOf(p.MyString))
	return values
}

type Slice struct {
	Bool                     []bool
	BoolPtr                  []*bool
	MyBool                   []MyBool
	MyBoolPtr                []*MyBool
	Int                      []int
	IntPtr                   []*int
	MyInt                    []MyInt
	MyIntPtr                 []*MyInt
	Int8                     []int8
	Int8Ptr                  []*int8
	MyInt8                   []MyInt8
	MyInt8Ptr                []*MyInt8
	Int16                    []int16
	Int16Ptr                 []*int16
	MyInt16                  []MyInt16
	MyInt16Ptr               []*MyInt16
	Int32                    []int32
	Int32Ptr                 []*int32
	MyInt32                  []MyInt32
	MyInt32Ptr               []*MyInt32
	Int64                    []int64
	Int64Ptr                 []*int64
	MyInt64                  []MyInt64
	MyInt64Ptr               []*MyInt64
	Uint                     []uint
	UintPtr                  []*uint
	MyUint                   []MyUint
	MyUintPtr                []*MyUint
	Uint8                    []uint8
	Uint8Ptr                 []*uint8
	MyUint8                  []MyUint8
	MyUint8Ptr               []*MyUint8
	Uint16                   []uint16
	Uint16Ptr                []*uint16
	MyUint16                 []MyUint16
	MyUint16Ptr              []*MyUint16
	Uint32                   []uint32
	Uint32Ptr                []*uint32
	MyUint32                 []MyUint32
	MyUint32Ptr              []*MyUint32
	Uint64                   []uint64
	Uint64Ptr                []*uint64
	MyUint64                 []MyUint64
	MyUint64Ptr              []*MyUint64
	Uintptr                  []uintptr
	UintptrPtr               []*uintptr
	MyUintptr                []MyUintptr
	MyUintptrPtr             []*MyUintptr
	Float32                  []float32
	Float32Ptr               []*float32
	MyFloat32                []MyFloat32
	MyFloat32Ptr             []*MyFloat32
	Float64                  []float64
	Float64Ptr               []*float64
	MyFloat64                []MyFloat64
	MyFloat64Ptr             []*MyFloat64
	Complex64                []complex64
	Complex64Ptr             []*complex64
	MyComplex64              []MyComplex64
	MyComplex64Ptr           []*MyComplex64
	Complex128               []complex128
	Complex128Ptr            []*complex128
	MyComplex128             []MyComplex128
	MyComplex128Ptr          []*MyComplex128
	String                   []string
	StringPtr                []*string
	MyString                 []MyString
	MyStringPtr              []*MyString
	MySlice                  MySlice
	Struct                   []Predeclared
	StructPtr                []*Predeclared
	UnexportedFieldStruct    []UnexportedFieldStruct
	UnexportedFieldStructPtr []*UnexportedFieldStruct
}

func NewSlice() Slice {
	p := NewPredeclared()
	uep := NewUnexportedFieldStruct()
	return Slice{
		Bool:                     []bool{p.Bool},
		BoolPtr:                  []*bool{&p.Bool},
		Complex128:               []complex128{p.Complex128},
		Complex128Ptr:            []*complex128{&p.Complex128},
		Complex64:                []complex64{p.Complex64},
		Complex64Ptr:             []*complex64{&p.Complex64},
		Float32:                  []float32{p.Float32},
		Float32Ptr:               []*float32{&p.Float32},
		Float64:                  []float64{p.Float64},
		Float64Ptr:               []*float64{&p.Float64},
		Int:                      []int{p.Int},
		Int16:                    []int16{p.Int16},
		Int16Ptr:                 []*int16{&p.Int16},
		Int32:                    []int32{p.Int32},
		Int32Ptr:                 []*int32{&p.Int32},
		Int64:                    []int64{p.Int64},
		Int64Ptr:                 []*int64{&p.Int64},
		Int8:                     []int8{p.Int8},
		Int8Ptr:                  []*int8{&p.Int8},
		IntPtr:                   []*int{&p.Int},
		MyBool:                   []MyBool{p.MyBool},
		MyBoolPtr:                []*MyBool{&p.MyBool},
		MyComplex128:             []MyComplex128{p.MyComplex128},
		MyComplex128Ptr:          []*MyComplex128{&p.MyComplex128},
		MyComplex64:              []MyComplex64{p.MyComplex64},
		MyComplex64Ptr:           []*MyComplex64{&p.MyComplex64},
		MyFloat32:                []MyFloat32{p.MyFloat32},
		MyFloat32Ptr:             []*MyFloat32{&p.MyFloat32},
		MyFloat64:                []MyFloat64{p.MyFloat64},
		MyFloat64Ptr:             []*MyFloat64{&p.MyFloat64},
		MyInt:                    []MyInt{p.MyInt},
		MyInt16:                  []MyInt16{p.MyInt16},
		MyInt16Ptr:               []*MyInt16{&p.MyInt16},
		MyInt32:                  []MyInt32{p.MyInt32},
		MyInt32Ptr:               []*MyInt32{&p.MyInt32},
		MyInt64:                  []MyInt64{p.MyInt64},
		MyInt64Ptr:               []*MyInt64{&p.MyInt64},
		MyInt8:                   []MyInt8{p.MyInt8},
		MyInt8Ptr:                []*MyInt8{&p.MyInt8},
		MyIntPtr:                 []*MyInt{&p.MyInt},
		MySlice:                  MySlice{"myslice"},
		MyString:                 []MyString{p.MyString},
		MyStringPtr:              []*MyString{&p.MyString},
		MyUint:                   []MyUint{p.MyUint},
		MyUint16:                 []MyUint16{p.MyUint16},
		MyUint16Ptr:              []*MyUint16{&p.MyUint16},
		MyUint32:                 []MyUint32{p.MyUint32},
		MyUint32Ptr:              []*MyUint32{&p.MyUint32},
		MyUint64:                 []MyUint64{p.MyUint64},
		MyUint64Ptr:              []*MyUint64{&p.MyUint64},
		MyUint8:                  []MyUint8{p.MyUint8},
		MyUint8Ptr:               []*MyUint8{&p.MyUint8},
		MyUintPtr:                []*MyUint{&p.MyUint},
		MyUintptr:                []MyUintptr{p.MyUintptr},
		MyUintptrPtr:             []*MyUintptr{&p.MyUintptr},
		String:                   []string{p.String},
		StringPtr:                []*string{&p.String},
		Struct:                   []Predeclared{p},
		StructPtr:                []*Predeclared{&p},
		Uint:                     []uint{p.Uint},
		Uint16:                   []uint16{p.Uint16},
		Uint16Ptr:                []*uint16{&p.Uint16},
		Uint32:                   []uint32{p.Uint32},
		Uint32Ptr:                []*uint32{&p.Uint32},
		Uint64:                   []uint64{p.Uint64},
		Uint64Ptr:                []*uint64{&p.Uint64},
		Uint8:                    []uint8{p.Uint8},
		Uint8Ptr:                 []*uint8{&p.Uint8},
		UintPtr:                  []*uint{&p.Uint},
		Uintptr:                  []uintptr{p.Uintptr},
		UintptrPtr:               []*uintptr{&p.Uintptr},
		UnexportedFieldStruct:    []UnexportedFieldStruct{uep},
		UnexportedFieldStructPtr: []*UnexportedFieldStruct{&uep},
	}
}

func (p Slice) Values() []reflect.Value {
	var values []reflect.Value
	values = append(values, reflect.ValueOf(p.Bool), reflect.ValueOf(p.BoolPtr), reflect.ValueOf(p.MyBool), reflect.ValueOf(p.MyBoolPtr))
	values = append(values, reflect.ValueOf(p.Int), reflect.ValueOf(p.IntPtr), reflect.ValueOf(p.MyInt), reflect.ValueOf(p.MyIntPtr))
	values = append(values, reflect.ValueOf(p.Int8), reflect.ValueOf(p.Int8Ptr), reflect.ValueOf(p.MyInt8), reflect.ValueOf(p.MyInt8Ptr))
	values = append(values, reflect.ValueOf(p.Int16), reflect.ValueOf(p.Int16Ptr), reflect.ValueOf(p.MyInt16), reflect.ValueOf(p.MyInt16Ptr))
	values = append(values, reflect.ValueOf(p.Int32), reflect.ValueOf(p.Int32Ptr), reflect.ValueOf(p.MyInt32), reflect.ValueOf(p.MyInt32Ptr))
	values = append(values, reflect.ValueOf(p.Int64), reflect.ValueOf(p.Int64Ptr), reflect.ValueOf(p.MyInt64), reflect.ValueOf(p.MyInt64Ptr))
	values = append(values, reflect.ValueOf(p.Uint), reflect.ValueOf(p.UintPtr), reflect.ValueOf(p.MyUint), reflect.ValueOf(p.MyUintPtr))
	values = append(values, reflect.ValueOf(p.Uint8), reflect.ValueOf(p.Uint8Ptr), reflect.ValueOf(p.MyUint8), reflect.ValueOf(p.MyUint8Ptr))
	values = append(values, reflect.ValueOf(p.Uint16), reflect.ValueOf(p.Uint16Ptr), reflect.ValueOf(p.MyUint16), reflect.ValueOf(p.MyUint16Ptr))
	values = append(values, reflect.ValueOf(p.Uint32), reflect.ValueOf(p.Uint32Ptr), reflect.ValueOf(p.MyUint32), reflect.ValueOf(p.MyUint32Ptr))
	values = append(values, reflect.ValueOf(p.Uint64), reflect.ValueOf(p.Uint64Ptr), reflect.ValueOf(p.MyUint64), reflect.ValueOf(p.MyUint64Ptr))
	values = append(values, reflect.ValueOf(p.Uintptr), reflect.ValueOf(p.UintptrPtr), reflect.ValueOf(p.MyUintptr), reflect.ValueOf(p.MyUintptrPtr))
	values = append(values, reflect.ValueOf(p.Float32), reflect.ValueOf(p.Float32Ptr), reflect.ValueOf(p.MyFloat32), reflect.ValueOf(p.MyFloat32Ptr))
	values = append(values, reflect.ValueOf(p.Float64), reflect.ValueOf(p.Float64Ptr), reflect.ValueOf(p.MyFloat64), reflect.ValueOf(p.MyFloat64Ptr))
	values = append(values, reflect.ValueOf(p.Complex64), reflect.ValueOf(p.Complex64Ptr), reflect.ValueOf(p.MyComplex64), reflect.ValueOf(p.MyComplex64Ptr))
	values = append(values, reflect.ValueOf(p.Complex128), reflect.ValueOf(p.Complex128Ptr), reflect.ValueOf(p.MyComplex128), reflect.ValueOf(p.MyComplex128Ptr))
	values = append(values, reflect.ValueOf(p.String), reflect.ValueOf(p.StringPtr), reflect.ValueOf(p.MyString), reflect.ValueOf(p.MyStringPtr))
	values = append(values, reflect.ValueOf(p.MySlice), reflect.ValueOf(p.Struct), reflect.ValueOf(p.StructPtr), reflect.ValueOf(p.UnexportedFieldStruct), reflect.ValueOf(p.UnexportedFieldStructPtr))
	return values
}

type Map struct {
	Bool                     map[bool]bool
	BoolPtr                  map[bool]*bool
	MyBool                   map[MyBool]MyBool
	MyBoolPtr                map[MyBool]*MyBool
	Int                      map[int]int
	IntPtr                   map[int]*int
	MyInt                    map[MyInt]MyInt
	MyIntPtr                 map[MyInt]*MyInt
	Int8                     map[int8]int8
	Int8Ptr                  map[int8]*int8
	MyInt8                   map[MyInt8]MyInt8
	MyInt8Ptr                map[MyInt8]*MyInt8
	Int16                    map[int16]int16
	Int16Ptr                 map[int16]*int16
	MyInt16                  map[MyInt16]MyInt16
	MyInt16Ptr               map[MyInt16]*MyInt16
	Int32                    map[int32]int32
	Int32Ptr                 map[int32]*int32
	MyInt32                  map[MyInt32]MyInt32
	MyInt32Ptr               map[MyInt32]*MyInt32
	Int64                    map[int64]int64
	Int64Ptr                 map[int64]*int64
	MyInt64                  map[MyInt64]MyInt64
	MyInt64Ptr               map[MyInt64]*MyInt64
	Uint                     map[uint]uint
	UintPtr                  map[uint]*uint
	MyUint                   map[MyUint]MyUint
	MyUintPtr                map[MyUint]*MyUint
	Uint8                    map[uint8]uint8
	Uint8Ptr                 map[uint8]*uint8
	MyUint8                  map[MyUint8]MyUint8
	MyUint8Ptr               map[MyUint8]*MyUint8
	Uint16                   map[uint16]uint16
	Uint16Ptr                map[uint16]*uint16
	MyUint16                 map[MyUint16]MyUint16
	MyUint16Ptr              map[MyUint16]*MyUint16
	Uint32                   map[uint32]uint32
	Uint32Ptr                map[uint32]*uint32
	MyUint32                 map[MyUint32]MyUint32
	MyUint32Ptr              map[MyUint32]*MyUint32
	Uint64                   map[uint64]uint64
	Uint64Ptr                map[uint64]*uint64
	MyUint64                 map[MyUint64]MyUint64
	MyUint64Ptr              map[MyUint64]*MyUint64
	Uintptr                  map[uintptr]uintptr
	UintptrPtr               map[uintptr]*uintptr
	MyUintptr                map[MyUintptr]MyUintptr
	MyUintptrPtr             map[MyUintptr]*MyUintptr
	Float32                  map[float32]float32
	Float32Ptr               map[float32]*float32
	MyFloat32                map[MyFloat32]MyFloat32
	MyFloat32Ptr             map[MyFloat32]*MyFloat32
	Float64                  map[float64]float64
	Float64Ptr               map[float64]*float64
	MyFloat64                map[MyFloat64]MyFloat64
	MyFloat64Ptr             map[MyFloat64]*MyFloat64
	Complex64                map[complex64]complex64
	Complex64Ptr             map[complex64]*complex64
	MyComplex64              map[MyComplex64]MyComplex64
	MyComplex64Ptr           map[MyComplex64]*MyComplex64
	Complex128               map[complex128]complex128
	Complex128Ptr            map[complex128]*complex128
	MyComplex128             map[MyComplex128]MyComplex128
	MyComplex128Ptr          map[MyComplex128]*MyComplex128
	String                   map[string]string
	StringPtr                map[string]*string
	MyString                 map[MyString]MyString
	MyStringPtr              map[MyString]*MyString
	MyMap                    MyMap
	Struct                   map[string]Predeclared
	StructPtr                map[string]*Predeclared
	UnexportedFieldStruct    map[string]UnexportedFieldStruct
	UnexportedFieldStructPtr map[string]*UnexportedFieldStruct
}

func NewMap() Map {
	p := NewPredeclared()
	uep := NewUnexportedFieldStruct()
	return Map{
		Bool:                     map[bool]bool{p.Bool: p.Bool},
		BoolPtr:                  map[bool]*bool{p.Bool: &p.Bool},
		Complex128:               map[complex128]complex128{p.Complex128: p.Complex128},
		Complex128Ptr:            map[complex128]*complex128{p.Complex128: &p.Complex128},
		Complex64:                map[complex64]complex64{p.Complex64: p.Complex64},
		Complex64Ptr:             map[complex64]*complex64{p.Complex64: &p.Complex64},
		Float32:                  map[float32]float32{p.Float32: p.Float32},
		Float32Ptr:               map[float32]*float32{p.Float32: &p.Float32},
		Float64:                  map[float64]float64{p.Float64: p.Float64},
		Float64Ptr:               map[float64]*float64{p.Float64: &p.Float64},
		Int:                      map[int]int{p.Int: p.Int},
		Int16:                    map[int16]int16{p.Int16: p.Int16},
		Int16Ptr:                 map[int16]*int16{p.Int16: &p.Int16},
		Int32:                    map[int32]int32{p.Int32: p.Int32},
		Int32Ptr:                 map[int32]*int32{p.Int32: &p.Int32},
		Int64:                    map[int64]int64{p.Int64: p.Int64},
		Int64Ptr:                 map[int64]*int64{p.Int64: &p.Int64},
		Int8:                     map[int8]int8{p.Int8: p.Int8},
		Int8Ptr:                  map[int8]*int8{p.Int8: &p.Int8},
		IntPtr:                   map[int]*int{p.Int: &p.Int},
		MyBool:                   map[MyBool]MyBool{p.MyBool: p.MyBool},
		MyBoolPtr:                map[MyBool]*MyBool{p.MyBool: &p.MyBool},
		MyComplex128:             map[MyComplex128]MyComplex128{p.MyComplex128: p.MyComplex128},
		MyComplex128Ptr:          map[MyComplex128]*MyComplex128{p.MyComplex128: &p.MyComplex128},
		MyComplex64:              map[MyComplex64]MyComplex64{p.MyComplex64: p.MyComplex64},
		MyComplex64Ptr:           map[MyComplex64]*MyComplex64{p.MyComplex64: &p.MyComplex64},
		MyFloat32:                map[MyFloat32]MyFloat32{p.MyFloat32: p.MyFloat32},
		MyFloat32Ptr:             map[MyFloat32]*MyFloat32{p.MyFloat32: &p.MyFloat32},
		MyFloat64:                map[MyFloat64]MyFloat64{p.MyFloat64: p.MyFloat64},
		MyFloat64Ptr:             map[MyFloat64]*MyFloat64{p.MyFloat64: &p.MyFloat64},
		MyInt:                    map[MyInt]MyInt{p.MyInt: p.MyInt},
		MyInt16:                  map[MyInt16]MyInt16{p.MyInt16: p.MyInt16},
		MyInt16Ptr:               map[MyInt16]*MyInt16{p.MyInt16: &p.MyInt16},
		MyInt32:                  map[MyInt32]MyInt32{p.MyInt32: p.MyInt32},
		MyInt32Ptr:               map[MyInt32]*MyInt32{p.MyInt32: &p.MyInt32},
		MyInt64:                  map[MyInt64]MyInt64{p.MyInt64: p.MyInt64},
		MyInt64Ptr:               map[MyInt64]*MyInt64{p.MyInt64: &p.MyInt64},
		MyInt8:                   map[MyInt8]MyInt8{p.MyInt8: p.MyInt8},
		MyInt8Ptr:                map[MyInt8]*MyInt8{p.MyInt8: &p.MyInt8},
		MyIntPtr:                 map[MyInt]*MyInt{p.MyInt: &p.MyInt},
		MyMap:                    MyMap{"mymap": "mymap"},
		MyString:                 map[MyString]MyString{p.MyString: p.MyString},
		MyStringPtr:              map[MyString]*MyString{p.MyString: &p.MyString},
		MyUint:                   map[MyUint]MyUint{p.MyUint: p.MyUint},
		MyUint16:                 map[MyUint16]MyUint16{p.MyUint16: p.MyUint16},
		MyUint16Ptr:              map[MyUint16]*MyUint16{p.MyUint16: &p.MyUint16},
		MyUint32:                 map[MyUint32]MyUint32{p.MyUint32: p.MyUint32},
		MyUint32Ptr:              map[MyUint32]*MyUint32{p.MyUint32: &p.MyUint32},
		MyUint64:                 map[MyUint64]MyUint64{p.MyUint64: p.MyUint64},
		MyUint64Ptr:              map[MyUint64]*MyUint64{p.MyUint64: &p.MyUint64},
		MyUint8:                  map[MyUint8]MyUint8{p.MyUint8: p.MyUint8},
		MyUint8Ptr:               map[MyUint8]*MyUint8{p.MyUint8: &p.MyUint8},
		MyUintPtr:                map[MyUint]*MyUint{p.MyUint: &p.MyUint},
		MyUintptr:                map[MyUintptr]MyUintptr{p.MyUintptr: p.MyUintptr},
		MyUintptrPtr:             map[MyUintptr]*MyUintptr{p.MyUintptr: &p.MyUintptr},
		String:                   map[string]string{p.String: p.String},
		StringPtr:                map[string]*string{p.String: &p.String},
		Struct:                   map[string]Predeclared{"struct": p},
		StructPtr:                map[string]*Predeclared{"struct": &p},
		Uint:                     map[uint]uint{p.Uint: p.Uint},
		Uint16:                   map[uint16]uint16{p.Uint16: p.Uint16},
		Uint16Ptr:                map[uint16]*uint16{p.Uint16: &p.Uint16},
		Uint32:                   map[uint32]uint32{p.Uint32: p.Uint32},
		Uint32Ptr:                map[uint32]*uint32{p.Uint32: &p.Uint32},
		Uint64:                   map[uint64]uint64{p.Uint64: p.Uint64},
		Uint64Ptr:                map[uint64]*uint64{p.Uint64: &p.Uint64},
		Uint8:                    map[uint8]uint8{p.Uint8: p.Uint8},
		Uint8Ptr:                 map[uint8]*uint8{p.Uint8: &p.Uint8},
		UintPtr:                  map[uint]*uint{p.Uint: &p.Uint},
		Uintptr:                  map[uintptr]uintptr{p.Uintptr: p.Uintptr},
		UintptrPtr:               map[uintptr]*uintptr{p.Uintptr: &p.Uintptr},
		UnexportedFieldStruct:    map[string]UnexportedFieldStruct{"struct": uep},
		UnexportedFieldStructPtr: map[string]*UnexportedFieldStruct{"struct": &uep},
	}
}

func (p Map) Values() []reflect.Value {
	var values []reflect.Value
	values = append(values, reflect.ValueOf(p.Bool), reflect.ValueOf(p.BoolPtr), reflect.ValueOf(p.MyBool), reflect.ValueOf(p.MyBoolPtr))
	values = append(values, reflect.ValueOf(p.Int), reflect.ValueOf(p.IntPtr), reflect.ValueOf(p.MyInt), reflect.ValueOf(p.MyIntPtr))
	values = append(values, reflect.ValueOf(p.Int8), reflect.ValueOf(p.Int8Ptr), reflect.ValueOf(p.MyInt8), reflect.ValueOf(p.MyInt8Ptr))
	values = append(values, reflect.ValueOf(p.Int16), reflect.ValueOf(p.Int16Ptr), reflect.ValueOf(p.MyInt16), reflect.ValueOf(p.MyInt16Ptr))
	values = append(values, reflect.ValueOf(p.Int32), reflect.ValueOf(p.Int32Ptr), reflect.ValueOf(p.MyInt32), reflect.ValueOf(p.MyInt32Ptr))
	values = append(values, reflect.ValueOf(p.Int64), reflect.ValueOf(p.Int64Ptr), reflect.ValueOf(p.MyInt64), reflect.ValueOf(p.MyInt64Ptr))
	values = append(values, reflect.ValueOf(p.Uint), reflect.ValueOf(p.UintPtr), reflect.ValueOf(p.MyUint), reflect.ValueOf(p.MyUintPtr))
	values = append(values, reflect.ValueOf(p.Uint8), reflect.ValueOf(p.Uint8Ptr), reflect.ValueOf(p.MyUint8), reflect.ValueOf(p.MyUint8Ptr))
	values = append(values, reflect.ValueOf(p.Uint16), reflect.ValueOf(p.Uint16Ptr), reflect.ValueOf(p.MyUint16), reflect.ValueOf(p.MyUint16Ptr))
	values = append(values, reflect.ValueOf(p.Uint32), reflect.ValueOf(p.Uint32Ptr), reflect.ValueOf(p.MyUint32), reflect.ValueOf(p.MyUint32Ptr))
	values = append(values, reflect.ValueOf(p.Uint64), reflect.ValueOf(p.Uint64Ptr), reflect.ValueOf(p.MyUint64), reflect.ValueOf(p.MyUint64Ptr))
	values = append(values, reflect.ValueOf(p.Uintptr), reflect.ValueOf(p.UintptrPtr), reflect.ValueOf(p.MyUintptr), reflect.ValueOf(p.MyUintptrPtr))
	values = append(values, reflect.ValueOf(p.Float32), reflect.ValueOf(p.Float32Ptr), reflect.ValueOf(p.MyFloat32), reflect.ValueOf(p.MyFloat32Ptr))
	values = append(values, reflect.ValueOf(p.Float64), reflect.ValueOf(p.Float64Ptr), reflect.ValueOf(p.MyFloat64), reflect.ValueOf(p.MyFloat64Ptr))
	values = append(values, reflect.ValueOf(p.Complex64), reflect.ValueOf(p.Complex64Ptr), reflect.ValueOf(p.MyComplex64), reflect.ValueOf(p.MyComplex64Ptr))
	values = append(values, reflect.ValueOf(p.Complex128), reflect.ValueOf(p.Complex128Ptr), reflect.ValueOf(p.MyComplex128), reflect.ValueOf(p.MyComplex128Ptr))
	values = append(values, reflect.ValueOf(p.String), reflect.ValueOf(p.StringPtr), reflect.ValueOf(p.MyString), reflect.ValueOf(p.MyStringPtr))
	values = append(values, reflect.ValueOf(p.MyMap), reflect.ValueOf(p.Struct), reflect.ValueOf(p.StructPtr), reflect.ValueOf(p.UnexportedFieldStruct), reflect.ValueOf(p.UnexportedFieldStructPtr))
	return values
}
