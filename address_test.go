package mnemosyne

import (
	"testing"
	"unsafe"
)

func TestReadMemory(t *testing.T) {
	n := uint32(0x12345678)

	res := ReadMemory(unsafe.Pointer(&n), 4)
	expected := []uint8{0x78, 0x56, 0x34, 0x12}
	for i := range res {
		if res[i] != expected[i] {
			t.Errorf("Expected %x, got %x", expected, res)
		}
	}
}

func TestWriteMemory(t *testing.T) {
	n := uint32(0xdeadbeef)

	WriteMemory(unsafe.Pointer(&n), []uint8{0x78, 0x56, 0x34, 0x12})
	if n != 0x12345678 {
		t.Errorf("Expected 0x12345678, got %x", n)
	}
}

func TestFillMemory(t *testing.T) {
	n := uint32(0xdeadbeef)

	FillMemory(unsafe.Pointer(&n), 0x90, 3)
	if n != 0xde909090 {
		t.Errorf("Expected 0xde909090, got %x", n)
	}
}

func TestWrite(t *testing.T) {
	n := uint64(0x12345678deadbeef)

	Write(unsafe.Pointer(&n), uint8(0x90))
	if n != 0x12345678deadbe90 {
		t.Errorf("Expected 0x12345678deadbe90, got %x", n)
	}

	Write(unsafe.Pointer(&n), uint16(0xbaad))
	if n != 0x12345678deadbaad {
		t.Errorf("Expected 0x12345678deadbaad, got %x", n)
	}

	Write(unsafe.Pointer(&n), uint32(0xdeadbeef))
	if n != 0x12345678deadbeef {
		t.Errorf("Expected 0x12345678deadbeef, got %x", n)
	}

	Write(unsafe.Pointer(&n), uint64(0xdeadbeef12345678))
	if n != 0xdeadbeef12345678 {
		t.Errorf("Expected 0xdeadbeef12345678, got %x", n)
	}
}

func TestRead(t *testing.T) {
	n := uint64(0x12345678deadbeef)

	if val := Read[uint8](unsafe.Pointer(&n)); val != 0xef {
		t.Errorf("Expected 0xef, got %x", val)
	}

	if val := Read[uint16](unsafe.Pointer(&n)); val != 0xbeef {
		t.Errorf("Expected 0xbeef, got %x", val)
	}

	if val := Read[uint32](unsafe.Pointer(&n)); val != 0xdeadbeef {
		t.Errorf("Expected 0xdeadbeef, got %x", val)
	}

	if val := Read[uint64](unsafe.Pointer(&n)); val != 0x12345678deadbeef {
		t.Errorf("Expected 0x12345678deadbeef, got %x", val)
	}
}

type testStruct struct {
	a uint8
	b uint16
	c uint32
	d uint64
}

func TestWritePtrVal(t *testing.T) {
	obj := testStruct{
		a: 0x33,
		b: 0x9090,
		c: 0xbaadf00d,
		d: 0xdeadbeefdeadbeef,
	}
	ptr := unsafe.Pointer(&obj)
	ptrToPtr := unsafe.Pointer(&ptr)

	if WritePtrVal[uint8](ptrToPtr, unsafe.Offsetof(obj.a), 0x88); obj.a != 0x88 {
		t.Errorf("Expected 0x88, got %x", obj.a)
	}

	if WritePtrVal[uint16](ptrToPtr, unsafe.Offsetof(obj.b), 0xefef); obj.b != 0xefef {
		t.Errorf("Expected 0xefef, got %x", obj.b)
	}

	if WritePtrVal[uint32](ptrToPtr, unsafe.Offsetof(obj.c), 0x45454545); obj.c != 0x45454545 {
		t.Errorf("Expected 0x45454545, got %x", obj.c)
	}

	if WritePtrVal[uint64](ptrToPtr, unsafe.Offsetof(obj.d), 0x1234567887654321); obj.d != 0x1234567887654321 {
		t.Errorf("Expected 0x1234567887654321, got %x", obj.d)
	}
}

func TestReadPtrVal(t *testing.T) {
	obj := testStruct{
		a: 0x33,
		b: 0x9090,
		c: 0xbaadf00d,
		d: 0xdeadbeefdeadbeef,
	}
	ptr := unsafe.Pointer(&obj)
	ptrToPtr := unsafe.Pointer(&ptr)

	if val := ReadPtrVal[uint8](ptrToPtr, unsafe.Offsetof(obj.a)); val != obj.a {
		t.Errorf("Expected %x, got %x", obj.a, val)
	}

	if val := ReadPtrVal[uint16](ptrToPtr, unsafe.Offsetof(obj.b)); val != obj.b {
		t.Errorf("Expected %x, got %x", obj.b, val)
	}

	if val := ReadPtrVal[uint32](ptrToPtr, unsafe.Offsetof(obj.c)); val != obj.c {
		t.Errorf("Expected %x, got %x", obj.c, val)
	}

	if val := ReadPtrVal[uint64](ptrToPtr, unsafe.Offsetof(obj.d)); val != obj.d {
		t.Errorf("Expected %x, got %x", obj.d, val)
	}
}

type testStructMultilevelInner struct {
	x uint32
	y *uint32
	z uint32
}

type testStructMultilevel struct {
	a uint8
	b uint16
	c *uint32
	d *testStructMultilevelInner
	e uint64
}

func TestWriteMultilevelPtrVal(t *testing.T) {
	v1 := uint32(0xc0cac0ca)
	v2 := uint32(0xbaadf00d)

	innerObj := testStructMultilevelInner{
		x: 0x11223344,
		y: &v1,
		z: 0x56565656,
	}

	obj := testStructMultilevel{
		a: 0x33,
		b: 0x9090,
		c: &v2,
		d: &innerObj,
		e: 0xdeadbeefdeadbeef,
	}

	ptr := unsafe.Pointer(&obj)
	ptrToPtr := unsafe.Pointer(&ptr)

	if !WriteMultilevelPtrVal(ptrToPtr, []uintptr{unsafe.Offsetof(obj.a)}, uint8(0x88)) {
		t.Errorf("Failed to write value to a")
	}
	if obj.a != 0x88 {
		t.Errorf("Expected %v, got %v", 0x88, obj.a)
	}

	if !WriteMultilevelPtrVal(ptrToPtr, []uintptr{unsafe.Offsetof(obj.b)}, uint16(0xefef)) {
		t.Errorf("Failed to write value to b")
	}
	if obj.b != 0xefef {
		t.Errorf("Expected %v, got %v", 0xefef, obj.b)
	}

	if !WriteMultilevelPtrVal(ptrToPtr, []uintptr{unsafe.Offsetof(obj.c), 0}, uint32(0x45454545)) {
		t.Errorf("Failed to write value to c")
	}
	if *obj.c != 0x45454545 {
		t.Errorf("Expected %v, got %v", 0x45454545, *obj.c)
	}

	if !WriteMultilevelPtrVal(ptrToPtr, []uintptr{unsafe.Offsetof(obj.d), unsafe.Offsetof(innerObj.x)}, uint32(0x11111111)) {
		t.Errorf("Failed to write value to inner x")
	}
	if innerObj.x != 0x11111111 {
		t.Errorf("Expected %v, got %v", 0x11111111, innerObj.x)
	}

	if !WriteMultilevelPtrVal(ptrToPtr, []uintptr{unsafe.Offsetof(obj.d), unsafe.Offsetof(innerObj.y), 0}, uint32(0x77777777)) {
		t.Errorf("Failed to write value to inner y")
	}
	if *innerObj.y != 0x77777777 {
		t.Errorf("Expected %v, got %v", 0x77777777, *innerObj.y)
	}

	if !WriteMultilevelPtrVal(ptrToPtr, []uintptr{unsafe.Offsetof(obj.d), unsafe.Offsetof(innerObj.z)}, uint32(0x66666666)) {
		t.Errorf("Failed to write value to inner z")
	}
	if innerObj.z != 0x66666666 {
		t.Errorf("Expected %v, got %v", 0x66666666, innerObj.z)
	}

	if !WriteMultilevelPtrVal(ptrToPtr, []uintptr{unsafe.Offsetof(obj.e)}, uint64(0x1234567887654321)) {
		t.Errorf("Failed to write value to e")
	}
	if obj.e != 0x1234567887654321 {
		t.Errorf("Expected %v, got %v", 0x1234567887654321, obj.e)
	}
}

func TestReadMultilevelPtrVal(t *testing.T) {
	v1 := uint32(0xc0cac0ca)
	v2 := uint32(0xbaadf00d)

	innerObj := testStructMultilevelInner{
		x: 0x11223344,
		y: &v1,
		z: 0x56565656,
	}

	obj := testStructMultilevel{
		a: 0x33,
		b: 0x9090,
		c: &v2,
		d: &innerObj,
		e: 0xdeadbeefdeadbeef,
	}

	ptr := unsafe.Pointer(&obj)
	ptrToPtr := unsafe.Pointer(&ptr)

	if val := ReadMultilevelPtrVal[uint8](ptrToPtr, []uintptr{unsafe.Offsetof(obj.a)}); val != obj.a {
		t.Errorf("Expected %v, got %v", obj.a, val)
	}
	if val := ReadMultilevelPtrVal[uint16](ptrToPtr, []uintptr{unsafe.Offsetof(obj.b)}); val != obj.b {
		t.Errorf("Expected %v, got %v", obj.b, val)
	}
	if val := ReadMultilevelPtrVal[uint32](ptrToPtr, []uintptr{unsafe.Offsetof(obj.c), 0}); val != *obj.c {
		t.Errorf("Expected %v, got %v", *obj.c, val)
	}
	if val := ReadMultilevelPtrVal[uint32](ptrToPtr, []uintptr{unsafe.Offsetof(obj.d), unsafe.Offsetof(innerObj.x)}); val != innerObj.x {
		t.Errorf("Expected %v, got %v", innerObj.x, val)
	}
	if val := ReadMultilevelPtrVal[uint32](ptrToPtr, []uintptr{unsafe.Offsetof(obj.d), unsafe.Offsetof(innerObj.y), 0}); val != *innerObj.y {
		t.Errorf("Expected %v, got %v", *innerObj.y, val)
	}
	if val := ReadMultilevelPtrVal[uint32](ptrToPtr, []uintptr{unsafe.Offsetof(obj.d), unsafe.Offsetof(innerObj.z)}); val != innerObj.z {
		t.Errorf("Expected %v, got %v", innerObj.z, val)
	}
	if val := ReadMultilevelPtrVal[uint64](ptrToPtr, []uintptr{unsafe.Offsetof(obj.e)}); val != obj.e {
		t.Errorf("Expected %v, got %v", obj.e, val)
	}
}
