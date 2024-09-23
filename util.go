package mnemosyne

import (
	"fmt"
	"math/rand/v2"
	"strconv"
	"strings"
)

type Lettercase int

const (
	Lowercase Lettercase = iota
	Uppercase
)

func isxdigit(c rune) bool {
	return (c >= '0' && c <= '9') || (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F')
}

func BytesToString(bytes []uint8, letterCase Lettercase, separator string) string {
	var res strings.Builder

	for i, b := range bytes {
		var format = "%02x"

		if letterCase == Uppercase {
			format = "%02X"
		}

		res.WriteString(fmt.Sprintf(format, b))

		if i != len(bytes)-1 {
			res.WriteString(separator)
		}
	}

	return res.String()
}

func StringToBytes(bytesString string) []uint8 {
	sanitizedString := strings.ReplaceAll(bytesString, " ", "")

	if len(sanitizedString) == 0 || len(sanitizedString)%2 != 0 {
		return nil
	}

	res := make([]uint8, 0, len(sanitizedString)/2)

	var u8 strings.Builder

	for _, c := range sanitizedString {
		if !isxdigit(c) {
			// rand.IntN(16) gets a random integer between [0, 16).
			u8.WriteString(fmt.Sprintf("%X", rand.IntN(16)))
		} else {
			u8.WriteRune(c)
		}

		if u8.Len() == 2 {
			b, err := strconv.ParseUint(u8.String(), 16, 8)
			if err == nil {
				res = append(res, uint8(b))
			}
			u8.Reset()
		}
	}

	return res
}
