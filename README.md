# Radix Tree

[![Test](https://github.com/racsoraul/radixtree/actions/workflows/go.yml/badge.svg?branch=master)](https://github.com/racsoraul/radixtree/actions/workflows/go.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/racsoraul/radixtree.svg)](https://pkg.go.dev/github.com/racsoraul/radixtree)

A fast, efficient Radix Tree implementation in Go. Provides a lexicographically ordered iteration and multiple lookup
methods. It leverages Go iterators for a more natural API when walking the tree.

## Features

- **Lexicographically Ordered Iteration**: Iterate over all keys in the tree in their natural order.
- **Iterators**: Walk the tree using the `for-range` loop.
- **Multiple Lookup Methods**: Supports exact match, longest prefix match, and prefix-based discovery (autocomplete).
- **High Performance**: Optimized for speed and low memory overhead.
- **Generics**: Store any data without losing type-safety.
- **Zero-Allocations**: No allocations during `Get` operations.

> 🚧 This project is still a WIP.

## Installation

```bash
go get github.com/racsoraul/radixtree
```

## API Reference

The full API documentation is available on [GoDoc](https://pkg.go.dev/github.com/racsoraul/radixtree).

## Usage

### Basic Operations

```go
package main

import (
	"fmt"

	"github.com/racsoraul/radixtree"
)

func main() {
	// Create a tree that holds integer values (can be any type).
	t := radixtree.New[int]()

	// Insert entries.
	t.Set("crash", 72)
	t.Set("ant", 20)
	t.Set("anagram", 40)
	t.Set("car", 30)
	t.Set("antihero", 100)
	t.Set("height", 11)
	t.Set("antares", 50)

	// Tree size (number of entries).
	fmt.Println(t.Len()) // Output: 7

	// Get entries.
	fmt.Println(t.Get("height")) // Output: 11 true
	fmt.Println(t.Get("care"))   // Output: 0 false

	// Walk the tree in lexicographical order.
	for k, v := range t.All() {
		fmt.Println("~>", k, v)
	}

	// Longest prefix with entries that are key nodes.
	fmt.Println(t.LongestPrefix("antagonist")) // Output: ant

	// Keys with prefix (autocompletion).
	fmt.Println(t.KeysWithPrefix("an", 5)) // Output: [anagram ant antares antihero]
}
```

### Set

The `Set` method inserts or updates a key-value pair in the tree.

```go
tree.Set("hello", 100)
```

### Get

The `Get` method retrieves the value associated with a key. It returns the value and a boolean indicating if the
key was found or not.

```go
if value, ok := tree.Get("hello"); ok {
	fmt.Println(value)
}
```

### Iteration

The `All()` method returns a Go iterator, allowing you to use `for ... range` loops. It allows you to walk the tree in
lexicographical order.

```go
for key, value := range tree.All() {
	fmt.Printf("%s: %v\n", key, value)
}
// Output:
// anagram: 40
// ant: 20
// antares: 50
// antihero: 100
// car: 30
// crash: 72
// height: 11
```

### Longest Prefix Match

Find the longest prefix of a given string that exists as a key in the tree.

```go
// Returns "ant" if we look for "antagonist".
fmt.Println(tree.LongestPrefix("antagonist"))
// Output: ant
```

### Keys With Prefix (Autocomplete)

Get all keys that start with a specific prefix. The `limit` parameter controls the maximum number of keys to return. Use
`-1` for no limit.

```go
// Get up to (autocomplete) 10 keys with prefix "ant".
fmt.Println(tree.KeysWithPrefix("ant", 10))
// Output: [anagram ant antares antihero]
```

### Visualization

The tree provides a `String()` method to visualize its internal structure. Handy for debugging small trees.

```go
fmt.Println(tree)
// Output:
// an(4)
//  |__agram(1)
//  |__t(3)
//     |__ares(1)
//     |__ihero(1)
// c(2)
// |__ar(1)
// |__rash(1)
// height(1)
```

## Benchmarks

Benchmarks performed on Apple M1 Max. The benchmarks use a dataset of 370,105 English words for "Big" tests and a small
set of manually defined entries for "Small" tests.

### Insertion (Set)

Measures the performance of adding or updating entries in the tree. The benchmark pre-fills the tree with 10,000 entries
and then performs repeated insertions.

| Benchmark           | Operations | Speed       | Memory   | Allocs      |
|:--------------------|:-----------|:------------|:---------|:------------|
| `BenchmarkTree_Set` | 5,201,592  | 232.9 ns/op | 241 B/op | 5 allocs/op |

### Lookup (Get)

Measures the performance of retrieving a value by its exact key. The "Big Tree" benchmark performs a lookup for the
key "hoar" in a tree containing 370,105 words.

| Benchmark          | Operations | Speed       | Memory | Allocs      |
|:-------------------|:-----------|:------------|:-------|:------------|
| `Get` (Small Tree) | 78,383,768 | 15.32 ns/op | 0 B/op | 0 allocs/op |
| `Get` (Big Tree)   | 32,526,567 | 36.76 ns/op | 0 B/op | 0 allocs/op |

### Keys With Prefix (KeysWithPrefix)

Measures the performance of finding all keys that share a common prefix.

- **Small**: Searches for "h" in a small tree, returns 4 results.
- **Medium**: Searches for "a" in the Big Tree (370,105 words), limited to 12,500 results.
- **Big**: Searches for "a" in the Big Tree, returns all 25,417 matches.

| Benchmark                                 | Operations | Speed         | Memory       | Allocs           |
|:------------------------------------------|:-----------|:--------------|:-------------|:-----------------|
| `KeysWithPrefix` (Small: 4 results)       | 9,890,601  | 121.6 ns/op   | 165 B/op     | 6 allocs/op      |
| `KeysWithPrefix` (Medium: 12,500 results) | 2,751      | 420,147 ns/op | 354,704 B/op | 12,501 allocs/op |
| `KeysWithPrefix` (Big: 25,417 results)    | 1,297      | 899,874 ns/op | 724,944 B/op | 25,418 allocs/op |

### Iteration (All)

Measures the performance of iterating over the tree's entries using Go iterators.

| Benchmark                      | Operations | Speed            | Memory         | Allocs            |
|:-------------------------------|:-----------|:-----------------|:---------------|:------------------|
| `All` (Small: ~15 entries)     | 4,433,824  | 271.5 ns/op      | 144 B/op       | 15 allocs/op      |
| `All` (Medium: 10,000 entries) | 3,798      | 313,525 ns/op    | 117,976 B/op   | 10,002 allocs/op  |
| `All` (Big: 370,105 entries)   | 92         | 13,233,376 ns/op | 4,582,438 B/op | 370,082 allocs/op |