# Set Similarity Search in Go

[![Build Status](https://travis-ci.org/ekzhu/go-set-similarity-search.svg?branch=master)](https://travis-ci.org/ekzhu/go-set-similarity-search)
[![GoDoc](https://godoc.org/github.com/ekzhu/go-set-similarity-search?status.svg)](https://godoc.org/github.com/ekzhu/go-set-similarity-search)

This is a mirror implementation of the 
Python [SetSimilaritySearch](https://github.com/ekzhu/SetSimilaritySearch)
library in Go, with better performance.

## Benchmarks

Run `AllPairs` algorithm on 3.5 GHz Intel Core i7, 
using similarity function `jaccard` and similarity threshold 0.5.

| Dataset | Input Sets | Avg. Size | `go-set-similarity-search` Runtime | `SetSimilaritySearch` Runtime |
|---------|------------|-----------|---|---|
| [Pokec social network (relationships)](https://snap.stanford.edu/data/soc-Pokec.html): from-nodes are set IDs; to-nodes are elements | 1432693 | 27.31 | 1m25s | 10m49s |
| [LiveJournal](https://snap.stanford.edu/data/soc-LiveJournal1.html): from-nodes are set IDs; to-nodes are elements | 4308452 | 16.01 | 4m11s | 28m51s |

## Library Usage

For *All-Pairs*, 
it takes an input of a list of sets, and output pairs that meet the 
similarity threshold.

```go
import (
    "fmt"
    "go-set-similarity-search"
)


func main() {
    // Each raw set must be a slice of unique string tokens.
    rawSets := [][]string{
        []string{"a"},
        []string{"a", "b"},
        []string{"a", "b", "c"},
        []string{"a", "b", "c", "d"},
        []string{"a", "b", "c", "d", "e"},
    }
    // Use frequency order transformation to replace the string tokens
    // with integers.
    sets, _ := SetSimilaritySearch.FrequencyOrderTransform(rawSets)
    // Run all-pairs algorithm, get a channel of pairs.
    pairs, _ := SetSimilaritySearch.AllPairs(sets,    
        /*similarityFunctionName=*/"jaccard", 
        /*similarityThreshold=*/0.1)
    // The pairs contain the indexes of sets to the original
    // rawSets and sets slices.
    for pair := range pairs {
        fmt.Println(pair)
    }
}
```

For *Query*, it takes an input of a list of sets, and builds a search 
index that can compute any number of queries. Currently the search index 
only supports a static collection of sets with no updates.

```go
import (
    "fmt"
    "go-set-similarity-search"
)

func main() {
    // Each raw set must be a slice of unique string tokens.
    rawSets := [][]string{
        []string{"a"},
        []string{"a", "b"},
        []string{"a", "b", "c"},
        []string{"a", "b", "c", "d"},
        []string{"a", "b", "c", "d", "e"},
    }
    // Use frequency order transformation to replace the string tokens
    // with integers.
    sets, dict := SetSimilaritySearch.FrequencyOrderTransform(rawSets)
    // Build a search index.
    searchIndex, err := SetSimilaritySearch.NewSearchIndex(sets,
        /*similarityFunctionName=*/"jaccard", 
        /*similarityThreshold=*/0.1)
    // Use dictionary to transform a query set.
    querySet := dict.Transform([]string{"a", "c", "d"})
    // Query the search index.
    searchResults := searchIndex.Query(querySet)
    // The results contain the indexes of sets to the original
    // rawSets and sets slices.
    for _, result := range searchResults {
        fmt.Println(result)
    }
}
```

Supported similarity functions (more to come):
* [Jaccard](https://en.wikipedia.org/wiki/Jaccard_index): intersection size divided by union size; set `similarityFunctionName="jaccard"`.
* [Cosine](https://en.wikipedia.org/wiki/Cosine_similarity): intersection size divided by square root of the product of sizes; set `similarityFunctionName="cosine"`.
* [Containment](https://ekzhu.github.io/datasketch/lshensemble.html#containment): intersection size divided by the size of the first set (or query set); set `similarityFunctionName="containment"`.
