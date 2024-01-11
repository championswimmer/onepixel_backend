package utils

import (
	"github.com/samber/lo"
)

const Alphabet = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz-_"

var AlphabetIndex = map[rune]uint64{}

const MaxSafeStringLength = 10

// MaxSafeNumber is the largest 64-bit number that can be represented in radix64
// 64^10 < 2^64, but 64^11 > 2^64
const MaxSafeNumber = 1152921504606846976 // 64^10

type Radix64Error struct {
	message string
}

func (e *Radix64Error) Error() string {
	return e.message
}

var Radix64TooLargeNumberError = &Radix64Error{"Number greater than 1152921504606846976"}
var Radix64StringTooLongError = &Radix64Error{"String longer than 10 characters"}

func init() {
	for i, c := range Alphabet {
		AlphabetIndex[c] = uint64(i)
	}
}

func Radix64Encode(number uint64) (string, error) {
	if number > MaxSafeNumber {
		return "", Radix64TooLargeNumberError
	}
	var result []byte
	if number == 0 { // special case
		return "0", nil
	}
	for number > 0 {
		result = append(result, Alphabet[number%64])
		number /= 64
	}
	return string(lo.Reverse(result)), nil
}

func Radix64Decode(str string) (uint64, error) {
	if len(str) > MaxSafeStringLength {
		return 0, Radix64StringTooLongError
	}

	var result uint64
	for _, c := range str {
		result *= 64
		result += uint64(AlphabetIndex[c])
	}
	return result, nil
}
