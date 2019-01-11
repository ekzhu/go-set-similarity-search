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
	correctdict := map[string]int{
		"a": 4,
		"b": 3,
		"c": 2,
		"d": 1,
		"e": 0,
	}
	sets, dict := FrequencyOrderTransform(rawSets)
	for i := range sets {
		for j := range sets[i] {
			if sets[i][j] != correctSets[i][j] {
				t.Errorf("Expect transformed set %v got %v", correctSets[i],
					sets[i])
			}
		}
	}
	for rawToken := range dict {
		if dict[rawToken] != correctdict[rawToken] {
			t.Errorf("Expect %v's dict is %v, got %v", rawToken,
				correctdict[rawToken], dict[rawToken])
		}
	}

	rawSet := []string{"a", "e", "b", "f"}
	correctSet := []int{0, 3, 4}
	set := dict.Transform(rawSet)
	for i := range set {
		if set[i] != correctSet[i] {
			t.Errorf("Expect transformed set %v got %v", correctSet, set)
		}
	}
	if len(set) != len(correctSet) {
		t.Errorf("Expect transformed set %v got %v", correctSet, set)
	}
}
