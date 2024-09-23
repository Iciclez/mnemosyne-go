package mnemosyne

import (
	"unsafe"

	"golang.org/x/exp/constraints"
)

func ReadMemory(ptr unsafe.Pointer, size int) []uint8 {
	memory := make([]uint8, 0, size)
	for i := 0; i < size; i++ {
		memory = append(memory, *(*uint8)(unsafe.Add(ptr, i)))
	}
	return memory
}

func WriteMemory(ptr unsafe.Pointer, bytes []uint8) {
	for i := 0; i < len(bytes); i++ {
		*(*uint8)(unsafe.Add(ptr, i)) = bytes[i]
	}
}

func FillMemory(ptr unsafe.Pointer, u8 uint8, size int) {
	for i := 0; i < size; i++ {
		*(*uint8)(unsafe.Add(ptr, i)) = u8
	}
}

func Write[T constraints.Integer](ptr unsafe.Pointer, data T) {
	*(*T)(ptr) = data
}

func Read[T constraints.Integer](ptr unsafe.Pointer) T {
	return *(*T)(ptr)
}

func WritePtrVal[T constraints.Integer](ptr unsafe.Pointer, offset uintptr, value T) bool {
	if ptr == nil {
		return false
	}

	ptrToVal := unsafe.Pointer(*(*uintptr)(ptr) + offset)
	*(*T)(ptrToVal) = value
	return true
}

func ReadPtrVal[T constraints.Integer](ptr unsafe.Pointer, offset uintptr) T {
	if ptr == nil {
		return 0
	}

	ptrToVal := unsafe.Pointer(*(*uintptr)(ptr) + offset)
	return *(*T)(ptrToVal)
}

func WriteMultilevelPtrVal[T constraints.Integer](ptr unsafe.Pointer, offsets []uintptr, value T) bool {
	if ptr == nil {
		return false
	}

	base := *(*uintptr)(ptr)
	for i, offset := range offsets {
		if i == len(offsets)-1 {
			ptrToVal := unsafe.Pointer(base + offset)
			*(*T)(ptrToVal) = value
			return true
		} else {
			base = *(*uintptr)(unsafe.Pointer(base + offset))
		}
	}

	return false
}

func ReadMultilevelPtrVal[T constraints.Integer](ptr unsafe.Pointer, offsets []uintptr) T {
	if ptr == nil {
		return 0
	}

	base := *(*uintptr)(ptr)
	for i, offset := range offsets {
		if i == len(offsets)-1 {
			ptrToVal := unsafe.Pointer(base + offset)
			return *(*T)(ptrToVal)
		} else {
			base = *(*uintptr)(unsafe.Pointer(base + offset))
		}
	}

	return 0
}
