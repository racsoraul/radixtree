# Radix Tree

[![Test](https://github.com/racsoraul/radixtree/actions/workflows/go.yml/badge.svg?branch=master)](https://github.com/racsoraul/radixtree/actions/workflows/go.yml)

A fast, efficient Radix Tree implementation in Go. Provides a lexicographically ordered iteration and multiple lookup
methods. It leverages Go iterators for a more natural API.

## Features

- **Lexicographically Ordered Iteration**: Iterate over all keys in the tree in their natural order.
- **Iterators**: Walk the tree using the `for-range` loop.
- **Multiple Lookup Methods**: Supports exact match, longest prefix match, and prefix-based discovery.
- **High Performance**: Optimized for speed and low memory overhead.
- **Generics**: Tree uses generics, so can store any data without losing type-safety.
- **Zero-Allocations**: No allocations during `Get` operations.

> 🚧 This project is still a WIP.

## Installation

```bash
go get github.com/racsoraul/radixtree
```

## Usage

### Basic Operations

```go
package main

import (
	"fmt"

	"github.com/racsoraul/radixtree"
)

func main() {
	tree := radixtree.New[string]()

	// Insert entries
	tree.Set("apple", "A sweet red fruit")
	tree.Set("app", "A small application")
	tree.Set("banana", "A long yellow fruit")

	// Get an entry.
	if val, ok := tree.Get("apple"); ok {
		fmt.Printf("apple: %v\n", val)
	}

	// Check tree size.
	fmt.Printf("Tree size: %d\n", tree.Len())
	// Output:
	// apple: A sweet red fruit
	// Tree size: 3
}
```

### Iteration

The `All()` method returns a Go iterator, allowing you to use `for ... range` loops.

```go
for key, value := range tree.All() {
	fmt.Printf("%s: %v\n", key, value)
}
// Output:
// app: A small application
// apple: A sweet red fruit
// banana: A long yellow fruit
```

### Longest Prefix Match

Find the longest prefix of a given string that exists as a key in the tree.

```go
// Returns "app" if only "app" and "apple" are in the tree and we look for "application".
prefix := tree.LongestPrefix("application")
fmt.Printf("Longest prefix: %s\n", prefix)
// Output: Longest prefix: app
```

### Prefix Search

Get all keys that start with a specific prefix. The `limit` parameter controls the maximum number of keys to return. Use
`-1` for no limit.

```go
// Get up to 10 keys with prefix "ap".
keys := tree.KeysWithPrefix("ap", 10)
fmt.Println("Keys with prefix 'ap':", keys)
// Output: Keys with prefix 'ap': [app apple]
```

### Visualization

The tree provides a `String()` method to visualize its internal structure. Handy for debugging small trees.

```go
fmt.Println(tree)
// Output:
// app(2)
//   |__le(1)
// banana(1)
```

## API Reference

The full API documentation is available on [GoDoc](https://pkg.go.dev/github.com/racsoraul/radixtree).

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

### Prefix Search (KeysWithPrefix)

Measures the performance of finding all keys that share a common prefix.

- **Small**: Searches for "h" in a small tree, returns 4 results.
- **Medium**: Searches for "a" in the Big Tree (370,105 words), limited to 12,500 results.
- **Big**: Searches for "a" in the Big Tree, returns all 25,417 matches.

| Benchmark                                 | Operations | Speed           | Memory       | Allocs           |
|:------------------------------------------|:-----------|:----------------|:-------------|:-----------------|
| `KeysWithPrefix` (Small: 4 results)       | 9,890,601  | 121.6 ns/op     | 165 B/op     | 6 allocs/op      |
| `KeysWithPrefix` (Medium: 12,500 results) | 2,751      | 420,147 ns/op   | 354,704 B/op | 12,501 allocs/op |
| `KeysWithPrefix` (Big: 25,417 results)    | 1,297      | 899,874 ns/op   | 724,944 B/op | 25,418 allocs/op |

### Iteration (All)

Measures the performance of iterating over the tree's entries using Go iterators.

| Benchmark                      | Operations | Speed            | Memory         | Allocs            |
|:-------------------------------|:-----------|:-----------------|:---------------|:------------------|
| `All` (Small: ~15 entries)     | 4,433,824  | 271.5 ns/op      | 144 B/op       | 15 allocs/op      |
| `All` (Medium: 10,000 entries) | 3,798      | 313,525 ns/op    | 117,976 B/op   | 10,002 allocs/op  |
| `All` (Big: 370,105 entries)   | 92         | 13,233,376 ns/op | 4,582,438 B/op | 370,082 allocs/op |