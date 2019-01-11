package SetSimilaritySearch

import "sort"

// Dictionary maps raw token to an integer token in the global order.
type Dictionary map[string]int

// FrequencyOrderTransform transforms tokens to integers according to global
// frequency order. This step replaces all original tokens in the sets with
// integers, and helps to speed up subsequent prefix filtering and similarity
// computation.  See Section 4.3.2 in the paper "A Primitive Operator for
// Similarity Joins in Data Cleaning" by Chaudhuri et al..
func FrequencyOrderTransform(rawSets [][]string) (sets [][]int,
	dict Dictionary) {
	// Count token frequencies.
	counts := make(map[string]int)
	for _, rawSet := range rawSets {
		for _, rawToken := range rawSet {
			if _, exists := counts[rawToken]; !exists {
				counts[rawToken] = 0
			}
			counts[rawToken]++
		}
	}
	// Create token order based on global frequency.
	type entry struct {
		rawToken string
		freq     int
	}
	entries := make([]entry, 0, len(counts))
	for rawToken, freq := range counts {
		entries = append(entries, entry{rawToken, freq})
	}
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].freq < entries[j].freq
	})
	dict = make(Dictionary)
	for i, entry := range entries {
		dict[entry.rawToken] = i
	}
	// Convert raw tokens into integer tokens.
	sets = make([][]int, len(rawSets))
	for i, rawSet := range rawSets {
		sets[i] = make([]int, len(rawSet))
		for j, rawToken := range rawSet {
			sets[i][j] = dict[rawToken]
		}
		sort.Ints(sets[i])
	}
	return sets, dict
}

// Transform takes a set of raw tokens and returns a set of integer tokens
// based on a global order.
func (dict Dictionary) Transform(rawSet []string) (set []int) {
	set = make([]int, 0, len(rawSet))
	for _, rawToken := range rawSet {
		if token, exists := dict[rawToken]; exists {
			set = append(set, token)
		}
	}
	sort.Ints(set)
	return set
}
