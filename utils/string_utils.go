package utils

import (
	"errors"
	"unicode"
)

func IsASCII(s string) (bool, error) {
	for i := 0; i < len(s); i++ {
		if s[i] > unicode.MaxASCII {
			return false, errors.New("Error string " + s + "is not ascii")
		}
	}
	return true, nil
}
