package testutils

import (
	"fmt"
	"testing"
)

func TestAreStringSlicesEqual_True(t *testing.T) {
	testData := [][2][]string{
		{{}, {}},
		{{"a"}, {"a"}},
		{{"a", "b", "z"}, {"a", "b", "z"}},
	}

	for index, data := range testData {
		s1 := data[0]
		s2 := data[1]

		results, _ := AreStringSlicesEqual(s1, s2)

		if !results {
			errStart := fmt.Sprintf("String slices at index %d were not seen as equal when they should have been.\n", index)
			t.Errorf("%s   Slice 1: %v.\n   Slice 2: %v", errStart, s1, s2)
			return
		}
	}
}

func TestAreStringSlicesEqual_False(t *testing.T) {
	testData := [][2][]string{
		{{"a"}, {}},
		{{}, {"a"}},
		{{"a"}, {"b"}},
		{{"a", "b", "c"}, {"a", "b", "c", "z"}},
		{{"a", "b", "c", "z"}, {"b", "a", "c"}},
	}

	for index, data := range testData {
		s1 := data[0]
		s2 := data[1]

		results, _ := AreStringSlicesEqual(s1, s2)

		if results {
			errStart := fmt.Sprintf("String slices at index %d were seen as equal when they should NOT have been.\n", index)
			t.Errorf("%s   Slice 1: %v.\n   Slice 2: %v", errStart, s1, s2)
			return
		}
	}
}
