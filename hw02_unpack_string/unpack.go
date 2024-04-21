package hw02unpackstring

import (
	"errors"
	"strings"
	"unicode"

	"golang.org/x/example/hello/reverse"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(str string) (string, error) {
	var result strings.Builder
	multiple, numbersInRow := 1, 0

	if len(str) == 0 {
		return "", nil
	}

	if unicode.IsDigit(rune(str[0])) {
		return "", ErrInvalidString
	}
	for _, char := range reverse.String(str) {
		if char >= '0' && char <= '9' {
			multiple = int(char - '0')
			numbersInRow++
		} else {
			result.WriteString(strings.Repeat(string(char), multiple))
			multiple, numbersInRow = 1, 0
		}
		if numbersInRow > 1 {
			return "", ErrInvalidString
		}
	}
	return reverse.String(result.String()), nil
}
