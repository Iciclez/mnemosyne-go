package mnemosyne

import (
	"unsafe"

	"golang.org/x/exp/constraints"
)

type MemoryEdit interface {
	Edit()
	Revert()
}

type MemoryPatch struct {
	ptr          unsafe.Pointer
	replaceBytes []uint8
	retainBytes  []uint8
}

func NewMemoryPatch(ptr unsafe.Pointer, bytes []uint8) *MemoryPatch {
	return &MemoryPatch{
		ptr:          ptr,
		replaceBytes: bytes,
		retainBytes:  ReadMemory(ptr, len(bytes)),
	}
}

func (mp *MemoryPatch) Edit() {
	WriteMemory(mp.ptr, mp.replaceBytes)
}

func (mp *MemoryPatch) Revert() {
	WriteMemory(mp.ptr, mp.retainBytes)
}

type MemoryDataEdit[T constraints.Integer] struct {
	ptr         unsafe.Pointer
	replaceData T
	retainData  T
}

func NewMemoryDataEdit[T constraints.Integer](ptr unsafe.Pointer, data T) *MemoryDataEdit[T] {
	return &MemoryDataEdit[T]{
		ptr:         ptr,
		replaceData: data,
		retainData:  Read[T](ptr),
	}
}

func (mde *MemoryDataEdit[T]) Edit() {
	Write(mde.ptr, mde.replaceData)
}

func (mde *MemoryDataEdit[T]) Revert() {
	Write(mde.ptr, mde.retainData)
}
