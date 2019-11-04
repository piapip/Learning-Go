package main

import (
	"fmt"
	"testing"
)

func BenchmarkAbs(b *testing.B) {
	for i := 0; i < b.N; i++ {
		fmt.Sprintf("%d", Abs(-1))
	}
}
func TestAbs(t *testing.T) {
	if Abs(-1) != 1 {
		t.Error("Expected Abs(-1) to be -1")
	}
}

func TestTableAbs(t *testing.T) {
	var tests = []struct {
		input    int
		expected int
	}{
		{2, 2},
		{-1, 1},
		{3, 3},
		{4, 4},
		{-5, 5},
	}

	for _, test := range tests {
		if output := Abs(test.input); output != test.expected {
			t.Errorf("Ab(%d) should be %d", test.input, test.expected)
		}
	}
}
