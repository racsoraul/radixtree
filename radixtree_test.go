package radixtree

import (
	"testing"
)

func TestLongestCommonPrefix(t *testing.T) {
	testCases := []struct {
		entry, label                                             string
		expectedPrefix, expectedEntrySuffix, expectedLabelSuffix string
	}{
		// Case one: "entry" is a prefix of "label".
		{"help", "helping", "help", "", "ing"},
		// Case two: "label" is a prefix of "entry".
		{"helping", "help", "help", "ing", ""},
		// Case three: "entry" and "label" share a common prefix.
		{"hello", "helping", "hel", "lo", "ping"},
		// Case four: "entry" and "label" have no common prefix.
		{"hello", "world", "", "hello", "world"},
		// Case five: "entry" and "label" are identical.
		{"hello", "hello", "hello", "", ""},
	}

	for _, tc := range testCases {
		t.Run(tc.entry+"_"+tc.label, func(t *testing.T) {
			result, entrySuffix, labelSuffix := longestCommonPrefix(tc.entry, tc.label)
			if result != tc.expectedPrefix {
				t.Fatalf("Expected prefix %q, got %q", tc.expectedPrefix, result)
			}
			if entrySuffix != tc.expectedEntrySuffix {
				t.Fatalf("Expected suffixOne %q, got %q", tc.expectedEntrySuffix, entrySuffix)
			}
			if labelSuffix != tc.expectedLabelSuffix {
				t.Fatalf("Expected suffixTwo %q, got %q", tc.expectedLabelSuffix, labelSuffix)
			}
		})
	}
}

func TestTree_Search(t *testing.T) {
	testCases := []struct {
		name          string
		setup         func(*Tree)
		key           string
		expectedFound bool
		expectedValue any
	}{
		{
			name:  "Search entry in empty tree",
			setup: func(t *Tree) {},
			key:   "hello",
		},
		{
			name: "Single entry - exact match",
			setup: func(t *Tree) {
				t.Insert("hello", 50)
			},
			key:           "hello",
			expectedFound: true,
			expectedValue: 50,
		},
		{
			name: "Single entry - key is prefix of entry",
			setup: func(t *Tree) {
				t.Insert("hello", 50)
			},
			key: "he",
		},
		{
			name: "Single entry - entry is prefix of key",
			setup: func(t *Tree) {
				t.Insert("hello", 50)
			},
			key: "helloworld",
		},
		{
			name: "Singly entry - exact match for empty key",
			setup: func(t *Tree) {
				t.Insert("", -9)
			},
			key:           "",
			expectedFound: true,
			expectedValue: -9,
		},
		{
			name: "Single entry - no match for empty key",
			setup: func(t *Tree) {
				t.Insert("hello", 50)
			},
			key: "",
		},
		{
			name: "Multiple entries - key is prefix of entry and exact match of another one",
			setup: func(t *Tree) {
				t.Insert("hello", 50)
				t.Insert("he", 25)
			},
			key:           "he",
			expectedFound: true,
			expectedValue: 25,
		},
		{
			name: "Multiple entries - one entry is prefix of key and key is exact match of another one",
			setup: func(t *Tree) {
				t.Insert("he", 25)
				t.Insert("hello", 50)
			},
			key:           "hello",
			expectedFound: true,
			expectedValue: 50,
		},
		{
			name: "Search entry in tree with common prefix",
			setup: func(t *Tree) {
				t.Insert("hello", 50)
				t.Insert("hella", 51)
			},
			key:           "hello",
			expectedFound: true,
			expectedValue: 50,
		},
		{
			name: "Search entry in tree with common prefix",
			setup: func(t *Tree) {
				t.Insert("hello", 50)
				t.Insert("hella", 51)
			},
			key:           "hella",
			expectedFound: true,
			expectedValue: 51,
		},
		{
			name: "Search entry in tree with common prefix - no match",
			setup: func(t *Tree) {
				t.Insert("hello", 50)
				t.Insert("hella", 51)
			},
			key: "hellb",
		},
		{
			name: "Exact match update",
			setup: func(t *Tree) {
				t.Insert("hello", 50)
				t.Insert("hello", 500)
			},
			key:           "hello",
			expectedFound: true,
			expectedValue: 500,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tree := New()
			tc.setup(tree)
			value, found := tree.Search(tc.key)
			if found != tc.expectedFound {
				t.Fatalf("want found=%v for key %q, got found=%v", tc.expectedFound, tc.key, found)
			}
			if tc.expectedFound && tc.expectedValue != value {
				t.Fatalf("want value=%v for key %q, got value=%v", tc.expectedValue, tc.key, value)
			}
		})
	}
}
