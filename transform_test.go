package SetSimilaritySearch

import "testing"

func TestTransform(t *testing.T) {
	rawSets := [][]string{
		[]string{"a"},
		[]string{"a", "b"},
		[]string{"a", "b", "c"},
		[]string{"a", "b", "c", "d"},
		[]string{"a", "b", "c", "d", "e"},
	}
	correctSets := [][]int{
		[]int{4},
		[]int{3, 4},
		[]int{2, 3, 4},
		[]int{1, 2, 3, 4},
		[]int{0, 1, 2, 3, 4},
	}
	correctOrder := map[string]int{
		"a": 4,
		"b": 3,
		"c": 2,
		"d": 1,
		"e": 0,
	}
	sets, order := FrequencyOrderTransform(rawSets)
	for i := range sets {
		for j := range sets[i] {
			if sets[i][j] != correctSets[i][j] {
				t.Errorf("Expect transformed set %v got %v", correctSets[i],
					sets[i])
			}
		}
	}
	for rawToken := range order {
		if order[rawToken] != correctOrder[rawToken] {
			t.Errorf("Expect %v's order is %v, got %v", rawToken,
				correctOrder[rawToken], order[rawToken])
		}
	}

}
