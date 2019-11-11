package main

import "testing"

type AddResult struct {
	x        int
	y        int
	expected int
}

var sampleTests = []AddResult{
	{1, 1, 2},
	{2, 1, 3},
}

func TestAdd(t *testing.T) {
	for _, test := range sampleTests {
		result := Add(test.x, test.y)
		if result != test.expected {
			t.Fatal("Expected Result Not Given")
		}
	}
}
