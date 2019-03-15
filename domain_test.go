package domain

import (
	"testing"
	"fmt"
)

func TestIsStringInStringSlice_True(t *testing.T) {
	type TestCase struct {
		needle string
		haystack []string
	}

	testData := []TestCase{
		{"a", []string{"a"}},
		{"a", []string{"b", "c", "a"}},
	}

	for index, data := range testData {
		results := IsStringInStringSlice(data.needle, data.haystack)

		if !results {
			errStart := fmt.Sprintf("At testData index %d, the string was not found in the slice but should have been.\n", index)
			t.Errorf("%s   String: %s.\n   Slice: %v", errStart, data.needle, data.haystack)
			return
		}
	}
}

func TestIsStringInStringSlice_False(t *testing.T) {
	type TestCase struct {
		needle string
		haystack []string
	}

	testData := []TestCase{
		{"a", []string{}},
		{"a", []string{"b", "c"}},
	}

	for index, data := range testData {
		results := IsStringInStringSlice(data.needle, data.haystack)

		if results {
			errStart := fmt.Sprintf("At testData index %d, the string was found in the slice but should NOT have been.\n", index)
			t.Errorf("%s   String: %s.\n   Slice: %v", errStart, data.needle, data.haystack)
			return
		}
	}
}