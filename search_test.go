package SetSimilaritySearch

import "testing"

func resultExists(r SearchResult, results []SearchResult) bool {
	for i := range results {
		if r == results[i] {
			return true
		}
	}
	return false
}

func TestSearchIndexJaccard(t *testing.T) {
	sets := [][]int{
		[]int{1, 2, 3},
		[]int{3, 4, 5},
		[]int{2, 3, 4},
		[]int{5, 6, 7},
	}
	query := []int{3, 4, 5}
	correctResults := []SearchResult{
		SearchResult{1, 1.0},
		SearchResult{0, 0.2},
		SearchResult{2, 0.5},
		SearchResult{3, 0.2},
	}
	searchIndex, err := NewSearchIndex(sets, "jaccard", 0.1)
	if err != nil {
		t.Fatal(err)
	}
	results := searchIndex.Query(query)
	for _, r := range results {
		if !resultExists(r, correctResults) {
			t.Errorf("The result %v is not correct", r)
		}
	}
	if len(results) != len(correctResults) {
		t.Errorf("Expecting %d results got %d", len(correctResults),
			len(results))
	}
}

func TestSearchIndexContainment(t *testing.T) {
	sets := [][]int{
		[]int{1, 2, 3},
		[]int{3, 4, 5},
		[]int{2, 3, 4},
		[]int{5, 6, 7},
	}
	query := []int{3, 4, 5}
	correctResults := []SearchResult{
		SearchResult{1, 1.0},
		SearchResult{0, 1.0 / 3.0},
		SearchResult{2, 2.0 / 3.0},
		SearchResult{3, 1.0 / 3.0},
	}
	// Threshold 0.1
	searchIndex, err := NewSearchIndex(sets, "containment", 0.1)
	if err != nil {
		t.Fatal(err)
	}
	results := searchIndex.Query(query)
	for _, r := range results {
		if !resultExists(r, correctResults) {
			t.Errorf("The result %v is not correct", r)
		}
	}
	if len(results) != len(correctResults) {
		t.Errorf("Expecting %d results got %d", len(correctResults),
			len(results))
	}
	// Threshold 0.5
	correctResults = []SearchResult{
		SearchResult{1, 1.0},
		SearchResult{2, 2.0 / 3.0},
	}
	searchIndex, err = NewSearchIndex(sets, "containment", 0.5)
	if err != nil {
		t.Fatal(err)
	}
	results = searchIndex.Query(query)
	for _, r := range results {
		if !resultExists(r, correctResults) {
			t.Errorf("The result %v is not correct", r)
		}
	}
	if len(results) != len(correctResults) {
		t.Errorf("Expecting %d results got %d", len(correctResults),
			len(results))
	}
}
