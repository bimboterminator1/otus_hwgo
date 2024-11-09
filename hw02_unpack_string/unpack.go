package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

var ErrInvalidString = errors.New("invalid string")

type runeWithStatus struct {
	obj       rune
	isEscaped bool
	isDigit   bool
	isSlash   bool
}

func decodeRune(p *[]byte) (runeWithStatus, int) {
	curRune, size := utf8.DecodeRune(*p)

	return runeWithStatus{obj: curRune, isDigit: unicode.IsDigit(curRune), isSlash: string(curRune) == `\`}, size
}

func Unpack(str string) (string, error) {
	unpackedStr := strings.Builder{}
	byteStr := []byte(str)

	if utf8.RuneCount(byteStr) == 0 {
		return "", nil
	}
	prevRune, size := decodeRune(&byteStr)

	// we consider strings starting with a digit as invalid
	if prevRune.isDigit {
		return "", ErrInvalidString
	}

	byteStr = byteStr[size:]
	for utf8.RuneCount(byteStr) > 0 {
		nextRune, size := decodeRune(&byteStr)
		byteStr = byteStr[size:]

		if nextRune.isDigit && prevRune.isDigit && !prevRune.isEscaped {
			// only single digits are allowed
			return "", ErrInvalidString
		}
		if prevRune.isSlash && !prevRune.isEscaped {
			// allow escaping only backslash and digits
			if nextRune.isDigit || nextRune.isSlash {
				nextRune.isEscaped = true
				prevRune = nextRune
				continue
			}

			return "", ErrInvalidString
		}
		if nextRune.isDigit {
			seqNum, err := strconv.Atoi(string(nextRune.obj))
			if err != nil {
				return "", err
			}
			_, err = unpackedStr.WriteString(strings.Repeat(string(prevRune.obj), seqNum))
			if err != nil {
				return "", err
			}
		} else if !(prevRune.isDigit && !prevRune.isEscaped) { // don't print expanded digit a5b -> aaaaab
			_, err := unpackedStr.WriteRune(prevRune.obj)
			if err != nil {
				return "", err
			}
		}
		prevRune = nextRune
	}

	if prevRune.isSlash && !prevRune.isEscaped {
		return "", ErrInvalidString
	} else if !(prevRune.isDigit && !prevRune.isEscaped) {
		_, err := unpackedStr.WriteRune(prevRune.obj)
		if err != nil {
			return "", err
		}
	}
	return unpackedStr.String(), nil
}
