package SetSimilaritySearch

import (
	"errors"
	"sort"
)

// SearchIndex is a data structure supports set similarity search queries.
// The algorithm is a combination of the prefix filter and position filter
// techniques.
type SearchIndex struct {
	threshold                 float64
	simFunc                   function
	overlapThresholdFunc      overlapThresholdFunction
	overlapIndexThresholdFunc overlapThresholdFunction
	positionFilterFunc        positionFilter
	sets                      [][]int
	postingLists              map[int][]postingListEntry
}

// NewSearchIndex builds a search index on the transformed sets given
// the similarity function and threshold.
// Currently supported similarity functions are "jaccard", "cosine"
// and "containment".
func NewSearchIndex(sets [][]int, similarityFunctionName string,
	similarityThreshold float64) (*SearchIndex, error) {
	if len(sets) == 0 {
		return nil, errors.New("input sets cannot be empty")
	}
	if similarityThreshold < 0 || similarityThreshold > 1.0 {
		return nil, errors.New("input similarityThreshold must be in the range [0, 1]")
	}
	si := SearchIndex{
		threshold:    similarityThreshold,
		sets:         sets,
		postingLists: make(map[int][]postingListEntry),
	}
	if f, exists := similarityFuncs[similarityFunctionName]; exists {
		si.simFunc = f
	} else {
		return nil, errors.New("input similarityFunctionName is not supported")
	}
	si.overlapThresholdFunc = overlapThresholdFuncs[similarityFunctionName]
	si.overlapIndexThresholdFunc = overlapIndexThresholdFuncs[similarityFunctionName]
	si.positionFilterFunc = positionFilterFuncs[similarityFunctionName]
	// Index transformed sets.
	for i, s := range sets {
		t := si.overlapIndexThresholdFunc(len(s), si.threshold)
		prefixSize := len(s) - t + 1
		prefix := s[:prefixSize]
		for j, token := range prefix {
			if _, exists := si.postingLists[token]; !exists {
				si.postingLists[token] = make([]postingListEntry, 0)
			}
			si.postingLists[token] = append(si.postingLists[token],
				postingListEntry{i, j, len(s)})
		}
	}
	// Sort each posting lists by set size for length filter.
	for token, postingList := range si.postingLists {
		sort.Slice(si.postingLists[token], func(i, j int) bool {
			return postingList[i].setSize < postingList[j].setSize
		})
	}
	return &si, nil
}

// SearchResult corresponding a set found from a query.
// It contains the index of the set found and the similarity to the query set.
type SearchResult struct {
	X          int
	Similarity float64
}

// Query probes the search index for sets whose similarity with the query set
// are above the given similarity threshold specified for the index.
// This function takes a transformed set and
// returns a slice of SearchResult that contain the indexes of the sets found.
func (si *SearchIndex) Query(s []int) []SearchResult {
	t := si.overlapThresholdFunc(len(s), si.threshold)
	prefixSize := len(s) - t + 1
	prefix := s[:prefixSize]
	// Find candidates using tokens in the prefix.
	candidates := make([]int, 0)
	for p1, token := range prefix {
		// TODO: use binary search to find starting position.
		// TODO: stops at an ending position for symmetric function.
		for _, entry := range si.postingLists[token] {
			if si.positionFilterFunc(s, si.sets[entry.setIndex], p1,
				entry.tokenPosition, si.threshold) {
				candidates = append(candidates, entry.setIndex)
			}
		}
	}
	// Sort and iterate through candidate indexes to verify
	// pairs.
	// TODO: optimize using partial overlaps.
	sort.Ints(candidates)
	results := make([]SearchResult, 0)
	prevCandidate := -1
	for _, x2 := range candidates {
		// Skip seen candidate.
		if x2 == prevCandidate {
			continue
		}
		prevCandidate = x2
		// Compute the exact similarity of this candidate
		sim := si.simFunc(s, si.sets[x2])
		if sim < si.threshold {
			continue
		}
		results = append(results, SearchResult{x2, sim})
	}
	return results
}
