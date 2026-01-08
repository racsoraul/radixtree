package radixtree

import "testing"

func TestLongestCommonPrefix(t *testing.T) {
	testCases := []struct {
		text, label                                          string
		expectedPrefix, expectedSuffixOne, expectedSuffixTwo string
	}{
		// Case one: "text" is a prefix of "label".
		{"help", "helping", "help", "", "ing"},
		// Case two: "label" is a prefix of "text".
		{"helping", "help", "help", "ing", ""},
		// Case three: "text" and "label" share a common prefix.
		{"hello", "helping", "hel", "lo", "ping"},
		// Case four: "text" and "label" have no common prefix.
		{"hello", "world", "", "hello", "world"},
		// Case five: "text" and "label" are identical.
		{"hello", "hello", "hello", "", ""},
	}

	for _, tc := range testCases {
		t.Run(tc.text+"_"+tc.label, func(t *testing.T) {
			result, suffixOne, suffixTwo := longestCommonPrefix(tc.text, tc.label)
			if result != tc.expectedPrefix {
				t.Fatalf("Expected prefix %q, got %q", tc.expectedPrefix, result)
			}
			if suffixOne != tc.expectedSuffixOne {
				t.Fatalf("Expected suffixOne %q, got %q", tc.expectedSuffixOne, suffixOne)
			}
			if suffixTwo != tc.expectedSuffixTwo {
				t.Fatalf("Expected suffixTwo %q, got %q", tc.expectedSuffixTwo, suffixTwo)
			}
		})
	}
}
