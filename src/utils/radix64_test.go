package utils

import "testing"

func TestRadix64Encode(t *testing.T) {
	var tests = []struct {
		input uint64
		want  string
	}{
		{0, "0"},
		{1, "1"},
		{10, "A"},
		{36, "a"},
		{62, "-"},
		{63, "_"},
		{64, "10"},
		{65, "11"},
		{128, "20"},
		{4095, "__"},
		{4096, "100"},
		{262144, "1000"},
		{1152921504606846975, "__________"},
	}
	for _, test := range tests {
		if got, _ := Radix64Encode(test.input); got != test.want {
			t.Errorf("Radix64Encode(%v) = %v", test.input, got)
		}
	}
}

func TestRadix64Decode(t *testing.T) {
	var tests = []struct {
		input string
		want  uint64
	}{
		{"0", 0},
		{"1", 1},
		{"A", 10},
		{"a", 36},
		{"-", 62},
		{"_", 63},
		{"10", 64},
		{"11", 65},
		{"20", 128},
		{"__", 4095},
		{"100", 4096},
		{"1000", 262144},
		{"__________", 1152921504606846975},
	}
	for _, test := range tests {
		if got, _ := Radix64Decode(test.input); got != test.want {
			t.Errorf("Radix64Decode(%v) = %v", test.input, got)
		}
	}
}
