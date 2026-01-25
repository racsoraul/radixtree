package radixtree

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"testing"
)

func TestLongestCommonPrefix(t *testing.T) {
	testCases := []struct {
		name                                                     string
		entry, label                                             string
		expectedPrefix, expectedEntrySuffix, expectedLabelSuffix string
	}{
		{
			name:                "Case one: entry is a prefix of label",
			entry:               "help",
			label:               "helping",
			expectedPrefix:      "help",
			expectedEntrySuffix: "",
			expectedLabelSuffix: "ing",
		},
		{
			name:                "Case two: label is prefix of entry",
			entry:               "helping",
			label:               "help",
			expectedPrefix:      "help",
			expectedEntrySuffix: "ing",
			expectedLabelSuffix: "",
		},
		{
			name:                "Case three: entry and label share a common prefix",
			entry:               "hello",
			label:               "helping",
			expectedPrefix:      "hel",
			expectedEntrySuffix: "lo",
			expectedLabelSuffix: "ping",
		},
		{
			name:                "Case four: entry and label have no common prefix",
			entry:               "hello",
			label:               "world",
			expectedPrefix:      "", // No edge.
			expectedEntrySuffix: "hello",
			expectedLabelSuffix: "world",
		},
		{
			name:                "Case five: entry and label are identical",
			entry:               "hello",
			label:               "hello",
			expectedPrefix:      "hello",
			expectedEntrySuffix: "",
			expectedLabelSuffix: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
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

func TestTree_SearchFromFile(t *testing.T) {
	tree, err := createTreeWithWordsFile()
	if err != nil {
		t.Fatal(err)
	}

	words, err := os.Open("./testdata/words_alpha.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer words.Close()

	scanner := bufio.NewScanner(words)
	for scanner.Scan() {
		word := scanner.Text()
		_, found := tree.Search(word)
		if !found {
			t.Fatalf("No match for word %q", word)
		}
	}
}

func TestTree_LongestPrefix(t *testing.T) {
	testCases := []struct {
		name           string
		setup          func(*Tree)
		key            string
		expectedPrefix string
	}{
		{
			name:  "empty tree",
			setup: func(t *Tree) {},
			key:   "hello",
		},
		{
			name:  "single entry, prefix of entry",
			setup: func(t *Tree) { t.Insert("hello", 50) },
			key:   "he",
		},
		{
			name:           "single entry, entry is prefix of key",
			setup:          func(t *Tree) { t.Insert("hello", 50) },
			key:            "hellothere",
			expectedPrefix: "hello",
		},
		{
			name:           "single entry, key is exact match of entry",
			setup:          func(t *Tree) { t.Insert("hello", 50) },
			key:            "hello",
			expectedPrefix: "hello",
		},
		{
			name:           "multiple entries, common prefix",
			setup:          setupMultipleEntries,
			key:            "antagonist",
			expectedPrefix: "ant",
		},
		{
			name:           "multiple entries, common prefix",
			setup:          setupMultipleEntries,
			key:            "antiacid",
			expectedPrefix: "ant",
		},
		{
			name:           "multiple entries, key is exact match of entry",
			setup:          setupMultipleEntries,
			key:            "ant",
			expectedPrefix: "ant",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tree := New()
			tc.setup(tree)
			prefix := tree.LongestPrefix(tc.key)
			if prefix != tc.expectedPrefix {
				t.Fatalf("want prefix=%q for key %q, got prefix=%q", tc.expectedPrefix, tc.key, prefix)
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
	t.Insert("antihero", 100)
	t.Insert("antecede", 101)
	t.Insert("antagony", 102)
}

// BenchmarkTree_Insert Measures the performance of the Insert method.
func BenchmarkTree_Insert(b *testing.B) {
	tree := New()
	for i := 0; i < 10_000; i++ {
		tree.Insert(fmt.Sprintf("%d", i), i)
	}

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		tree.Insert(strconv.Itoa(n), n)
	}
}

// BenchmarkTree_SearchBig Measures the performance of the Search method for a tree
// with more than ~370K entries in it.
func BenchmarkTree_SearchBig(b *testing.B) {
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

// BenchmarkTree_SearchSmall Measures the performance of the Search method for a small tree.
func BenchmarkTree_SearchSmall(b *testing.B) {
	tree := New()
	setupMultipleEntries(tree)

	for b.Loop() {
		_, found := tree.Search("height")
		if !found {
			b.Fatal("No match")
		}
	}
}

// BenchmarkTree_KeysWithPrefixBig Measures the performance of the KeysWithPrefix method by retrieving
// keys starting with the prefix "a", yielding 25,417 results.
func BenchmarkTree_KeysWithPrefixBig(b *testing.B) {
	tree, err := createTreeWithWordsFile()
	if err != nil {
		b.Fatal(err)
	}
	b.Logf("Loaded tree with %d entries", tree.Size())

	for b.Loop() {
		results := tree.KeysWithPrefix("a")
		if len(results) == 0 {
			b.Fatal("No results")
		}
	}
}

// BenchmarkTree_KeysWithPrefixSmall Measures the performance of the KeysWithPrefix method by retrieving
// keys starting with the prefix "h", yielding 25,417 results.
func BenchmarkTree_KeysWithPrefixSmall(b *testing.B) {
	tree := New()
	setupMultipleEntries(tree)

	for b.Loop() {
		results := tree.KeysWithPrefix("h")
		if len(results) == 0 {
			b.Fatal("No results")
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
