package SetSimilaritySearch

import (
	"bytes"
	"sort"
	"testing"
)

func TestReadFlattenedRawSets(t *testing.T) {
	testInput := `# Test input
1 a
1 b
1 c
1 d
2 a
2 b
2 c
3 a
3 b
3 c
3 f
4 f
4 g
4 h
`
	// Test forward.
	correctSetIDs := []string{"1", "2", "3", "4"}
	correctRawSets := [][]string{
		[]string{"a", "b", "c", "d"},
		[]string{"a", "b", "c"},
		[]string{"a", "b", "c", "f"},
		[]string{"f", "g", "h"},
	}
	file := bytes.NewBufferString(testInput)
	setIDs, rawSets, err := ReadFlattenedRawSets(file, false)
	if err != nil {
		t.Fatal(err)
	}
	for i := range setIDs {
		if setIDs[i] != correctSetIDs[i] {
			t.Errorf("Incorrect set ID %v", setIDs[i])
		}
	}
	for i, rawSet := range rawSets {
		for j := range rawSet {
			if rawSet[j] != correctRawSets[i][j] {
				t.Errorf("Incorrect raw set %v", rawSet)
			}
		}
	}

	// Test reverse.
	correctSetIDs = []string{"a", "b", "c", "d", "f", "g", "h"}
	correctRawSets = [][]string{
		[]string{"1", "2", "3"},
		[]string{"1", "2", "3"},
		[]string{"1", "2", "3"},
		[]string{"1"},
		[]string{"3", "4"},
		[]string{"4"},
		[]string{"4"},
	}
	file = bytes.NewBufferString(testInput)
	setIDs, rawSets, err = ReadFlattenedRawSets(file, true)
	if err != nil {
		t.Fatal(err)
	}
	for i := range setIDs {
		if setIDs[i] != correctSetIDs[i] {
			t.Errorf("Incorrect set ID %v", setIDs[i])
		}
	}
	for i, rawSet := range rawSets {
		sort.Strings(rawSet)
		for j := range rawSet {
			if rawSet[j] != correctRawSets[i][j] {
				t.Errorf("Incorrect raw set %v expecting %v", rawSet,
					correctRawSets[i])
			}
		}
	}
}

func TestReadFlattenedSortedRawSets(t *testing.T) {
	testInput := `# Test input
1 a
1 b
1 c
1 d
2 a
2 b
2 c
3 a
3 b
3 c
3 f
4 f
4 g
4 h
`
	// Test forward.
	correctSetIDs := []string{"1", "2", "3", "4"}
	correctRawSets := [][]string{
		[]string{"a", "b", "c", "d"},
		[]string{"a", "b", "c"},
		[]string{"a", "b", "c", "f"},
		[]string{"f", "g", "h"},
	}
	file := bytes.NewBufferString(testInput)
	setIDs, rawSets, err := ReadFlattenedSortedRawSets(file)
	if err != nil {
		t.Fatal(err)
	}
	for i := range setIDs {
		if setIDs[i] != correctSetIDs[i] {
			t.Errorf("Incorrect set ID %v", setIDs[i])
		}
	}
	for i, rawSet := range rawSets {
		for j := range rawSet {
			if rawSet[j] != correctRawSets[i][j] {
				t.Errorf("Incorrect raw set %v", rawSet)
			}
		}
	}
}
