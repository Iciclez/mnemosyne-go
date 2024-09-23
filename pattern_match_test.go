package mnemosyne

import (
	"testing"
	"unsafe"
)

func TestPatternMatchNew(t *testing.T) {
	patternMatch := NewPatternMatch("0a 0b ?? ??   ?? 2e ??   ?? ", nil, 0)
	sanitizedPattern := "0a0b??????2e"
	if patternMatch.pattern != sanitizedPattern {
		t.Errorf("Expected pattern %s, got %s", sanitizedPattern, patternMatch.pattern)
	}
	if patternMatch.patternSize != len(sanitizedPattern)/2 {
		t.Errorf("Expected pattern size %d, got %d", len(sanitizedPattern)/2, patternMatch.patternSize)
	}
}

func TestPatternMatchFindAddress(t *testing.T) {
	haystack := []uint8{
		0xf1, 0x80, 0xd7, 0x50, 0x1a, 0x7b, 0x69, 0x57, 0x07, 0x80, 0xbc, 0x27, 0xc7, 0x5e,
		0x88, 0x0c, 0xac, 0x7f, 0xd8, 0xe0, 0x13, 0x7d, 0xf4, 0xfb, 0xf4, 0x91, 0x0b, 0x07,
		0xa6, 0xe1, 0x54, 0x22,
	}

	patternMatch := NewPatternMatch("7b ?? 57 07 ?? bc ?? c7", unsafe.Pointer(&haystack[0]), len(haystack))
	address := patternMatch.FindAddress()
	if address != unsafe.Pointer(&haystack[5]) {
		t.Errorf("Expected address %v, got %v", uintptr(unsafe.Pointer(&haystack[5])), address)
	}

	patternMatch = NewPatternMatch("7b ?? 57 33 07 ?? bc ?? c7", unsafe.Pointer(&haystack), len(haystack))
	address = patternMatch.FindAddress()
	if address != nil {
		t.Errorf("Expected null address, got %v", address)
	}
}

func TestPatternMatchFindNextAddress(t *testing.T) {
	haystack := []uint8{
		0xf1, 0x80, 0xd7, 0x50, 0x1a, 0x7b, 0x69, 0x57, 0x07, 0x80, 0xbc, 0x27, 0xc7, 0x5e,
		0x88, 0x0c, 0xac, 0x7f, 0xd8, 0xe0, 0x13, 0x7d, 0xf4, 0xfb, 0xf4, 0x91, 0x0b, 0x07,
		0xa6, 0xe1, 0x54, 0x22, 0x7b, 0xa0, 0x57, 0x07, 0x2b, 0xbc, 0xdd, 0xc7, 0x24, 0x53,
		0xb3, 0x3f, 0xf1, 0xd5, 0x67, 0x23,
	}

	patternMatch := NewPatternMatch("7b ?? 57 07 ?? bc ?? c7", unsafe.Pointer(&haystack[0]), len(haystack))
	address := patternMatch.FindAddress()
	if address != unsafe.Pointer(&haystack[5]) {
		t.Errorf("Expected address %v, got %v", uintptr(unsafe.Pointer(&haystack[5])), address)
	}

	address = patternMatch.FindNextAddress()
	if address != unsafe.Pointer(&haystack[32]) {
		t.Errorf("Expected address %v, got %v", uintptr(unsafe.Pointer(&haystack[32])), address)
	}
}
