package SetSimilaritySearch

import "testing"

func pairExists(p Pair, pairs []Pair) bool {
	for i := range pairs {
		if p == pairs[i] {
			return true
		}
	}
	return false
}

func TestAllPairJaccard(t *testing.T) {
	sets := [][]int{
		[]int{1, 2, 3},
		[]int{3, 4, 5},
		[]int{2, 3, 4},
		[]int{5, 6, 7},
	}
	correctPairs := []Pair{
		Pair{1, 0, 0.2},
		Pair{2, 0, 0.5},
		Pair{2, 1, 0.5},
		Pair{3, 1, 0.2},
	}
	pairs, err := AllPairs(sets, "jaccard", 0.1)
	if err != nil {
		t.Fatal(err)
	}
	count := 0
	for p := range pairs {
		if !pairExists(p, correctPairs) {
			t.Errorf("The pair %v is not correct", p)
		}
		count++
	}
	if count != len(correctPairs) {
		t.Errorf("Expecting %d pairs but found %d", len(correctPairs), count)
	}
}
