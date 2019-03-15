package testutils

import "fmt"

func AreStringSlicesEqual(expected, results []string) (bool, string) {
	if len(expected) != len(results) {
		return false, fmt.Sprintf("Slices are not of equal length. Expected length: %d. But got length: %d",
			len(expected), len(results))
	}

	for index, nextExpected := range expected {
		nextResult := results[index]

		if nextExpected != nextResult {
			return false, fmt.Sprintf("Mismatch at index: %d. Expected: %s. But got: %s",
				index, nextExpected, nextResult)
		}
	}

	return true, ""
}
