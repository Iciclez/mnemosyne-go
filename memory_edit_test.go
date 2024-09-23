package mnemosyne

import (
	"testing"
	"unsafe"
)

func TestMemoryPatchEdit(t *testing.T) {
	n := uint32(0xdeadbeef)
	patch := NewMemoryPatch(unsafe.Pointer(&n), []uint8{0x78, 0x56, 0x34, 0x12})

	patch.Edit()
	if n != 0x12345678 {
		t.Errorf("Expected 0x12345678, got 0x%x", n)
	}

	patch.Revert()
	patch.Edit()
	if n != 0x12345678 {
		t.Errorf("Expected 0x12345678, got 0x%x", n)
	}
}

func TestMemoryPatchRevert(t *testing.T) {
	n := uint32(0xdeadbeef)
	patch := NewMemoryPatch(unsafe.Pointer(&n), []uint8{0x78, 0x56, 0x34, 0x12})

	patch.Edit()
	patch.Revert()
	if n != 0xdeadbeef {
		t.Errorf("Expected 0xdeadbeef, got 0x%x", n)
	}

	patch.Edit()
	patch.Revert()
	if n != 0xdeadbeef {
		t.Errorf("Expected 0xdeadbeef, got 0x%x", n)
	}
}

func TestMemoryDataEditEdit(t *testing.T) {
	n := uint32(0xdeadbeef)
	dataEdit := NewMemoryDataEdit(unsafe.Pointer(&n), uint32(0x12345678))

	dataEdit.Edit()
	if n != 0x12345678 {
		t.Errorf("Expected 0x12345678, got 0x%x", n)
	}

	dataEdit.Revert()
	dataEdit.Edit()
	if n != 0x12345678 {
		t.Errorf("Expected 0x12345678, got 0x%x", n)
	}
}

func TestMemoryDataEditRevert(t *testing.T) {
	n := uint32(0xdeadbeef)
	dataEdit := NewMemoryDataEdit(unsafe.Pointer(&n), uint32(0x12345678))

	dataEdit.Edit()
	dataEdit.Revert()
	if n != 0xdeadbeef {
		t.Errorf("Expected 0xdeadbeef, got 0x%x", n)
	}

	dataEdit.Edit()
	dataEdit.Revert()
	if n != 0xdeadbeef {
		t.Errorf("Expected 0xdeadbeef, got 0x%x", n)
	}
}
