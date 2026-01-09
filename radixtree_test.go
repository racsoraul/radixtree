package radixtree

import (
	"bufio"
	"os"
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
			name:  "search entry in empty tree",
			setup: func(t *Tree) {},
			key:   "hello",
		},
		{
			name:          "single entry, exact match",
			setup:         func(t *Tree) { t.Insert("hello", 50) },
			key:           "hello",
			expectedFound: true,
			expectedValue: 50,
		},
		{
			name:  "single entry, key is prefix of entry",
			setup: func(t *Tree) { t.Insert("hello", 50) },
			key:   "he",
		},
		{
			name:  "single entry, entry is prefix of key",
			setup: func(t *Tree) { t.Insert("hello", 50) },
			key:   "helloworld",
		},
		{
			name:          "singly entry, exact match for empty key",
			setup:         func(t *Tree) { t.Insert("", -9) },
			key:           "",
			expectedFound: true,
			expectedValue: -9,
		},
		{
			name:  "single entry, no match for empty key",
			setup: func(t *Tree) { t.Insert("hello", 50) },
			key:   "",
		},
		{
			name: "key is prefix of entry and exact match of another one",
			setup: func(t *Tree) {
				t.Insert("hello", 50)
				t.Insert("he", 25)
			},
			key:           "he",
			expectedFound: true,
			expectedValue: 25,
		},
		{
			name: "one entry is prefix of key and key is exact match of another one",
			setup: func(t *Tree) {
				t.Insert("he", 25)
				t.Insert("hello", 50)
			},
			key:           "hello",
			expectedFound: true,
			expectedValue: 50,
		},
		{
			name: "search entry in tree with common prefix",
			setup: func(t *Tree) {
				t.Insert("hello", 50)
				t.Insert("hella", 51)
			},
			key:           "hello",
			expectedFound: true,
			expectedValue: 50,
		},
		{
			name: "search entry in tree with common prefix",
			setup: func(t *Tree) {
				t.Insert("hello", 50)
				t.Insert("hella", 51)
			},
			key:           "hella",
			expectedFound: true,
			expectedValue: 51,
		},
		{
			name: "search entry in tree with common prefix, no match",
			setup: func(t *Tree) {
				t.Insert("hello", 50)
				t.Insert("hella", 51)
			},
			key: "hellb",
		},
		{
			name: "exact match update",
			setup: func(t *Tree) {
				t.Insert("hello", 50)
				t.Insert("hello", 500)
			},
			key:           "hello",
			expectedFound: true,
			expectedValue: 500,
		},
		{
			name: "exact match, no common prefix",
			setup: func(t *Tree) {
				t.Insert("hello", 50)
				t.Insert("world", 75)
			},
			key:           "hello",
			expectedFound: true,
			expectedValue: 50,
		},
		{
			name: "exact match, no common prefix",
			setup: func(t *Tree) {
				t.Insert("hello", 50)
				t.Insert("world", 75)
			},
			key:           "world",
			expectedFound: true,
			expectedValue: 75,
		},
		{
			name: "no match",
			setup: func(t *Tree) {
				t.Insert("hello", 50)
				t.Insert("world", 75)
			},
			key: "something",
		},
		{
			name:          "match deep key node",
			setup:         setupDeepTree,
			key:           "abcd",
			expectedFound: true,
			expectedValue: 1234,
		},
		{
			name:          "match intermediate key node",
			setup:         setupDeepTree,
			key:           "ab",
			expectedFound: true,
			expectedValue: 12,
		},
		{
			name:          "match intermediate key node",
			setup:         setupDeepTree,
			key:           "abc",
			expectedFound: true,
			expectedValue: 123,
		},
		{
			name:          "match first key node",
			setup:         setupDeepTree,
			key:           "a",
			expectedFound: true,
			expectedValue: 1,
		},
		{
			name:  "deep tree, no match",
			setup: setupDeepTree,
			key:   "abcde",
		},
		{
			name:          "multiple entries, match",
			setup:         setupMultipleEntries,
			key:           "crazy",
			expectedFound: true,
			expectedValue: 70,
		},
		{
			name:          "multiple entries, match",
			setup:         setupMultipleEntries,
			key:           "anagram",
			expectedFound: true,
			expectedValue: 40,
		},
		{
			name:  "multiple entries, no match",
			setup: setupMultipleEntries,
			key:   "ana",
		},
		{
			name: "single character match",
			setup: func(t *Tree) {
				t.Insert("a", 1)
				t.Insert("b", 2)
				t.Insert("c", 3)
			},
			key:           "b",
			expectedFound: true,
			expectedValue: 2,
		},
		{
			name: "long entry match",
			setup: func(t *Tree) {
				t.Insert("/", 1)
				t.Insert("/awesomedomain/api/v1/usermanagement/update", 10000000)
			},
			key:           "/awesomedomain/api/v1/usermanagement/update",
			expectedFound: true,
			expectedValue: 10000000,
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

func setupDeepTree(t *Tree) {
	t.Insert("a", 1)
	t.Insert("ab", 12)
	t.Insert("abc", 123)
	t.Insert("abcd", 1234)
}

func setupMultipleEntries(t *Tree) {
	t.Insert("he", 25)
	t.Insert("hello", 50)
	t.Insert("hella", 51)
	t.Insert("height", 60)
	t.Insert("ant", 20)
	t.Insert("anagram", 40)
	t.Insert("car", 30)
	t.Insert("crazy", 70)
	t.Insert("crash", 72)
}

// BenchmarkTree_Search Measures the performance of the Search method for a tree
// with more than ~370K entries in it.
func BenchmarkTree_Search(b *testing.B) {
	tree, err := createTreeWithWordsFile()
	if err != nil {
		b.Fatal(err)
	}
	b.Logf("Loaded tree with %d entries", tree.Size())

	for b.Loop() {
		_, found := tree.Search("hoar")
		if !found {
			b.Fatal("No match")
		}
	}
}

func createTreeWithWordsFile() (*Tree, error) {
	words, err := os.Open("./testdata/words_alpha.txt")
	if err != nil {
		return nil, err
	}
	defer words.Close()

	scanner := bufio.NewScanner(words)
	tree := New()
	wordCounter := 1
	for scanner.Scan() {
		word := scanner.Text()
		tree.Insert(word, wordCounter)
		wordCounter++
	}
	return tree, scanner.Err()
}
