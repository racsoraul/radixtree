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
			setup:         func(t *Tree) { t.Set("hello", 50) },
			key:           "hello",
			expectedFound: true,
			expectedValue: 50,
		},
		{
			name:  "single entry, key is prefix of entry",
			setup: func(t *Tree) { t.Set("hello", 50) },
			key:   "he",
		},
		{
			name:  "single entry, entry is prefix of key",
			setup: func(t *Tree) { t.Set("hello", 50) },
			key:   "helloworld",
		},
		{
			name:          "singly entry, exact match for empty key",
			setup:         func(t *Tree) { t.Set("", -9) },
			key:           "",
			expectedFound: true,
			expectedValue: -9,
		},
		{
			name:  "single entry, no match for empty key",
			setup: func(t *Tree) { t.Set("hello", 50) },
			key:   "",
		},
		{
			name: "key is prefix of entry and exact match of another one",
			setup: func(t *Tree) {
				t.Set("hello", 50)
				t.Set("he", 25)
			},
			key:           "he",
			expectedFound: true,
			expectedValue: 25,
		},
		{
			name: "one entry is prefix of key and key is exact match of another one",
			setup: func(t *Tree) {
				t.Set("he", 25)
				t.Set("hello", 50)
			},
			key:           "hello",
			expectedFound: true,
			expectedValue: 50,
		},
		{
			name: "search entry in tree with common prefix",
			setup: func(t *Tree) {
				t.Set("hello", 50)
				t.Set("hella", 51)
			},
			key:           "hello",
			expectedFound: true,
			expectedValue: 50,
		},
		{
			name: "search entry in tree with common prefix",
			setup: func(t *Tree) {
				t.Set("hello", 50)
				t.Set("hella", 51)
			},
			key:           "hella",
			expectedFound: true,
			expectedValue: 51,
		},
		{
			name: "search entry in tree with common prefix, no match",
			setup: func(t *Tree) {
				t.Set("hello", 50)
				t.Set("hella", 51)
			},
			key: "hellb",
		},
		{
			name: "exact match update",
			setup: func(t *Tree) {
				t.Set("hello", 50)
				t.Set("hello", 500)
			},
			key:           "hello",
			expectedFound: true,
			expectedValue: 500,
		},
		{
			name: "exact match, no common prefix",
			setup: func(t *Tree) {
				t.Set("hello", 50)
				t.Set("world", 75)
			},
			key:           "hello",
			expectedFound: true,
			expectedValue: 50,
		},
		{
			name: "exact match, no common prefix",
			setup: func(t *Tree) {
				t.Set("hello", 50)
				t.Set("world", 75)
			},
			key:           "world",
			expectedFound: true,
			expectedValue: 75,
		},
		{
			name: "no match",
			setup: func(t *Tree) {
				t.Set("hello", 50)
				t.Set("world", 75)
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
				t.Set("a", 1)
				t.Set("b", 2)
				t.Set("c", 3)
			},
			key:           "b",
			expectedFound: true,
			expectedValue: 2,
		},
		{
			name: "long entry match",
			setup: func(t *Tree) {
				t.Set("/", 1)
				t.Set("/awesomedomain/api/v1/usermanagement/update", 10000000)
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
			value, found := tree.Get(tc.key)
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
		_, found := tree.Get(word)
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
			setup: func(t *Tree) { t.Set("hello", 50) },
			key:   "he",
		},
		{
			name:           "single entry, entry is prefix of key",
			setup:          func(t *Tree) { t.Set("hello", 50) },
			key:            "hellothere",
			expectedPrefix: "hello",
		},
		{
			name:           "single entry, key is exact match of entry",
			setup:          func(t *Tree) { t.Set("hello", 50) },
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

func TestTree_KeysWithPrefix(t *testing.T) {
	tree := New()
	setupMultipleEntries(tree)

	testCases := []struct {
		prefix         string
		expectedResult []string
	}{
		{
			prefix:         "h",
			expectedResult: []string{"he", "height", "hella", "hello"},
		},
		{
			prefix:         "he",
			expectedResult: []string{"he", "height", "hella", "hello"},
		},
		{
			prefix:         "hel",
			expectedResult: []string{"hella", "hello"},
		},
		{
			prefix:         "hell",
			expectedResult: []string{"hella", "hello"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.prefix, func(t *testing.T) {
			results := tree.KeysWithPrefix(tc.prefix, -1)
			if len(results) != len(tc.expectedResult) {
				t.Fatalf("want %d results, got %d", len(tc.expectedResult), len(results))
			}
			for i, result := range results {
				if result != tc.expectedResult[i] {
					t.Fatalf("want result %q, got %q", tc.expectedResult[i], result)
				}
			}
		})
	}

}

func setupDeepTree(t *Tree) {
	t.Set("a", 1)
	t.Set("ab", 12)
	t.Set("abc", 123)
	t.Set("abcd", 1234)
}

func setupMultipleEntries(t *Tree) {
	t.Set("he", 25)
	t.Set("hello", 50)
	t.Set("hella", 51)
	t.Set("height", 60)
	t.Set("ant", 20)
	t.Set("anagram", 40)
	t.Set("car", 30)
	t.Set("crazy", 70)
	t.Set("crash", 72)
	t.Set("antihero", 100)
	t.Set("antecede", 101)
	t.Set("antagony", 102)
}

// BenchmarkTree_Set Measures the performance of the Set method.
func BenchmarkTree_Set(b *testing.B) {
	tree := New()
	for i := range 10_000 {
		tree.Set(fmt.Sprintf("%d", i), i)
	}

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		tree.Set(strconv.Itoa(n), n)
	}
}

// BenchmarkTree_GetSmall Measures the performance of the Get method for a small tree.
func BenchmarkTree_GetSmall(b *testing.B) {
	tree := New()
	setupMultipleEntries(tree)

	for b.Loop() {
		_, found := tree.Get("height")
		if !found {
			b.Fatal("No match")
		}
	}
}

// BenchmarkTree_GetBig Measures the performance of the Get method for a tree
// with more than ~370K entries in it.
func BenchmarkTree_GetBig(b *testing.B) {
	tree, err := createTreeWithWordsFile()
	if err != nil {
		b.Fatal(err)
	}
	b.Logf("Loaded tree with %d entries", tree.Len())

	for b.Loop() {
		_, found := tree.Get("hoar")
		if !found {
			b.Fatal("No match")
		}
	}
}

// BenchmarkTree_KeysWithPrefixSmall Measures the performance of the KeysWithPrefix method by retrieving
// keys starting with the prefix "h", 4 results.
func BenchmarkTree_KeysWithPrefixSmall(b *testing.B) {
	tree := New()
	setupMultipleEntries(tree)

	for b.Loop() {
		results := tree.KeysWithPrefix("h", -1)
		if len(results) != 4 {
			b.Fatalf("wrong number of results: %d", len(results))
		}
	}
}

// BenchmarkTree_KeysWithPrefixMedium Measures the performance of the KeysWithPrefix method by retrieving
// keys starting with the prefix "a", with a limit of 12,500.
func BenchmarkTree_KeysWithPrefixMedium(b *testing.B) {
	tree, err := createTreeWithWordsFile()
	if err != nil {
		b.Fatal(err)
	}
	b.Logf("Loaded tree with %d entries", tree.Len())

	for b.Loop() {
		results := tree.KeysWithPrefix("a", 12_500)
		if len(results) != 12_500 {
			b.Fatalf("wrong number of results: %d", len(results))
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
	b.Logf("Loaded tree with %d entries", tree.Len())

	for b.Loop() {
		results := tree.KeysWithPrefix("a", -1)
		if len(results) != 25_417 {
			b.Fatalf("wrong number of results: %d", len(results))
		}
	}
}

// BenchmarkTree_AllSmall Measures the performance of the All method by iterating over a small set of entries.
func BenchmarkTree_AllSmall(b *testing.B) {
	tree := New()
	setupMultipleEntries(tree)

	for b.Loop() {
		count := 0
		for k, v := range tree.All() {
			_ = k
			_ = v
			count++
		}
		if count != tree.Len() {
			b.Fatalf("wrong number of results: %d", count)
		}
	}
}

// BenchmarkTree_AllMedium Measures the performance of the All method by iterating over a big set of entries
// but breaking the iteration after a certain number of entries (10,000).
func BenchmarkTree_AllMedium(b *testing.B) {
	tree, err := createTreeWithWordsFile()
	if err != nil {
		b.Fatal(err)
	}
	b.Logf("Loaded tree with %d entries", tree.Len())

	for b.Loop() {
		count := 0
		for k, v := range tree.All() {
			count++
			_ = k
			_ = v
			if count == 10_000 {
				break
			}
		}
		if count != 10_000 {
			b.Fatalf("wrong number of results: %d", count)
		}
	}
}

// BenchmarkTree_AllBig Measures the performance of the All method by iterating over a big set of entries.
func BenchmarkTree_AllBig(b *testing.B) {
	tree, err := createTreeWithWordsFile()
	if err != nil {
		b.Fatal(err)
	}
	b.Logf("Loaded tree with %d entries", tree.Len())

	for b.Loop() {
		count := 0
		for k, v := range tree.All() {
			_ = k
			_ = v
			count++
		}
		if count != tree.Len() {
			b.Fatalf("wrong number of results: %d", count)
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
		tree.Set(word, wordCounter)
		wordCounter++
	}
	return tree, scanner.Err()
}
