package mnemosyne

import (
	"reflect"
	"testing"
)

func TestStringToBytes(t *testing.T) {
	expected := []uint8{0x12, 0x34, 0x56, 0x78, 0x90, 0xab, 0xcd, 0xef}
	result := StringToBytes("12 34 56 78 90 AB CD EF")

	if !reflect.DeepEqual(expected, result) {
		t.Errorf("Expected %v, but got %v", expected, result)
	}
}

func TestBytesToString(t *testing.T) {
	tests := []struct {
		bytes      []uint8
		letterCase Lettercase
		separator  string
		expected   string
	}{
		{[]uint8{0x12, 0x34, 0x56, 0x78, 0x90, 0xab, 0xcd, 0xef}, Uppercase, " ", "12 34 56 78 90 AB CD EF"},
		{[]uint8{0x12, 0x34, 0x56, 0x78, 0x90, 0xab, 0xcd, 0xef}, Lowercase, "", "1234567890abcdef"},
		{[]uint8{0x12, 0x34, 0x56, 0x78, 0x90, 0xab, 0xcd, 0xef}, Uppercase, "*", "12*34*56*78*90*AB*CD*EF"},
		{[]uint8{0x12, 0x34, 0x56, 0x78, 0x90, 0xab, 0xcd, 0xef}, Lowercase, "--", "12--34--56--78--90--ab--cd--ef"},
	}

	for _, test := range tests {
		result := BytesToString(test.bytes, test.letterCase, test.separator)
		if test.expected != result {
			t.Errorf("Expected %s, but got %s", test.expected, result)
		}
	}
}
