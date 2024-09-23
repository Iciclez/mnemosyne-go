package mnemosyne

import (
	"encoding/hex"
	"strings"
	"unsafe"
)

type PatternMatch struct {
	pattern        string
	patternSize    int
	memoryStart    unsafe.Pointer
	memorySize     int
	currentAddress uintptr
	byteArray      []uint8
	mask           []uint8
}

func NewPatternMatch(pattern string, memoryStart unsafe.Pointer, memorySize int) *PatternMatch {
	patternMatch := &PatternMatch{
		pattern:        pattern,
		memoryStart:    memoryStart,
		memorySize:     memorySize,
		currentAddress: uintptr(memoryStart),
		byteArray:      []uint8{},
		mask:           []uint8{},
	}

	patternMatch.pattern = strings.TrimRightFunc(patternMatch.pattern, func(r rune) bool {
		return r == '?' || r == ' ' || r == '\n' || r == '\t' || r == '\r'
	})

	patternMatch.pattern = strings.ReplaceAll(patternMatch.pattern, " ", "")

	if patternMatch.pattern == "" || len(patternMatch.pattern)%2 != 0 {
		return nil
	}

	patternMatch.patternSize = len(patternMatch.pattern) / 2

	patternMatch.byteArray = make([]uint8, 0, patternMatch.patternSize)
	patternMatch.mask = make([]uint8, 0, patternMatch.patternSize)

	for i := 0; i < len(patternMatch.pattern); i += 2 {
		if patternMatch.pattern[i] == '?' || patternMatch.pattern[i+1] == '?' {
			patternMatch.mask = append(patternMatch.mask, 1)
			patternMatch.byteArray = append(patternMatch.byteArray, 0)
		} else {
			patternMatch.mask = append(patternMatch.mask, 0)
			byteVal, err := hex.DecodeString(patternMatch.pattern[i : i+2])
			if err != nil {
				panic(err)
			}
			patternMatch.byteArray = append(patternMatch.byteArray, byteVal[0])
		}
	}

	return patternMatch
}

func (p *PatternMatch) FindAddress() unsafe.Pointer {
	return p.findAddressFrom(uintptr(p.memoryStart))
}

func (p *PatternMatch) FindNextAddress() unsafe.Pointer {
	return p.findAddressFrom(p.currentAddress + 1)
}

func (p *PatternMatch) findAddressFrom(addressFrom uintptr) unsafe.Pointer {
	memoryEnd := uintptr(unsafe.Add(p.memoryStart, p.memorySize))

	for p.currentAddress = addressFrom; p.currentAddress < memoryEnd; p.currentAddress++ {
		if p.tryMatchAtCurrentAddress() {
			return unsafe.Pointer(p.currentAddress)
		}
	}

	return nil
}

func (p *PatternMatch) tryMatchAtCurrentAddress() bool {
	for j := 0; j < p.patternSize; j++ {
		if p.mask[j] == 0 && *(*uint8)(unsafe.Add(unsafe.Pointer(p.currentAddress), j)) != p.byteArray[j] {
			return false
		}
	}
	return true
}
