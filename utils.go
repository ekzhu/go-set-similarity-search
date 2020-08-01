package SetSimilaritySearch

import (
	"bufio"
	"errors"
	"io"
	"sort"
	"strconv"
	"strings"
)

type flattenedRawSetEntry struct {
	setID    string
	rawToken string
}

// ReadFlattenedRawSets takes an input of a flattened set file,
// that contains unique lines in the format "<set ID> <token>", and returns
// the extracted set IDs and raw sets.
// Lines starting with "#" are ignored,
// If the input format is "<token> <set ID>" then set reversed to true.
func ReadFlattenedRawSets(file io.Reader,
	reversed bool) (setIDs []string, rawSets [][]string, err error) {
	// Read flattened raw set entries.
	entries := make([]flattenedRawSetEntry, 0)
	scanner := bufio.NewScanner(file)
	scanner.Buffer(nil, 1024*1024*1024*4)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "#") {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) != 2 {
			return nil, nil, errors.New("incorrect line detected")
		}
		var entry flattenedRawSetEntry
		if reversed {
			entry = flattenedRawSetEntry{fields[1], fields[0]}
		} else {
			entry = flattenedRawSetEntry{fields[0], fields[1]}
		}
		entries = append(entries, entry)
	}
	if err := scanner.Err(); err != nil {
		return nil, nil, err
	}
	// Sort entries by setID.
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].setID < entries[j].setID
	})
	// Create raw sets by merging flattened entries.
	setIDs = make([]string, 0)
	rawSets = make([][]string, 0)
	currSetID := entries[0].setID
	currSet := make([]string, 0)
	for _, entry := range entries {
		if entry.setID != currSetID {
			// Append the completed set.
			setIDs = append(setIDs, currSetID)
			rawSets = append(rawSets, currSet)
			// Create new set.
			currSetID = entry.setID
			currSet = make([]string, 0)
		}
		currSet = append(currSet, entry.rawToken)
	}
	// Append the last set.
	setIDs = append(setIDs, currSetID)
	rawSets = append(rawSets, currSet)
	return setIDs, rawSets, nil
}

// ReadFlattenedSortedRawSets takes an input of a flattened set file,
// that contains unique lines in the format "<set ID> <token>",
// sorted by <set ID>, and returns the extracted set IDs and raw sets.
// Lines starting with "#" are ignored,
// This function is more efficient than ReadFlattenedRawSets, but expects
// the input lines to be sorted.
func ReadFlattenedSortedRawSets(file io.Reader) (setIDs []string,
	rawSets [][]string, err error) {
	// Create raw sets by merging flattened entries.
	setIDs = make([]string, 0)
	rawSets = make([][]string, 0)
	var currSetID string
	firstLine := true
	currSet := make([]string, 0)
	scanner := bufio.NewScanner(file)
	scanner.Buffer(nil, 1024*1024*1024*4)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "#") {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) != 2 {
			return nil, nil, errors.New("incorrect line detected")
		}
		setID := fields[0]
		rawToken := fields[1]
		if firstLine {
			currSetID = setID
			firstLine = false
		}
		if setID != currSetID {
			// Append the completed set.
			setIDs = append(setIDs, currSetID)
			rawSets = append(rawSets, currSet)
			// Create new set.
			currSetID = setID
			currSet = make([]string, 0)
		}
		currSet = append(currSet, rawToken)
	}
	if err := scanner.Err(); err != nil {
		return nil, nil, err
	}
	// Append the last set.
	setIDs = append(setIDs, currSetID)
	rawSets = append(rawSets, currSet)
	return setIDs, rawSets, nil
}

// ReadFlattenedSortedTransformedSets takes an input of a flattened
// transformed set file,
// that contains unique lines in the format "<set ID:int> <token:int>",
// sorted by <set ID>, and returns the extracted set IDs and raw sets.
// Lines starting with "#" are ignored,
func ReadFlattenedSortedTransformedSets(file io.Reader) (setIDs []int,
	sets [][]int, err error) {
	// Create raw sets by merging flattened entries.
	setIDs = make([]int, 0)
	sets = make([][]int, 0)
	var currSetID int
	firstLine := true
	currSet := make([]int, 0)
	scanner := bufio.NewScanner(file)
	scanner.Buffer(nil, 1024*1024*1024*4)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "#") {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) != 2 {
			return nil, nil, errors.New("incorrect line detected")
		}
		setID, err := strconv.Atoi(fields[0])
		if err != nil {
			return nil, nil, err
		}
		token, err := strconv.Atoi(fields[1])
		if err != nil {
			return nil, nil, err
		}
		if firstLine {
			currSetID = setID
			firstLine = false
		}
		if setID != currSetID {
			// Append the completed set.
			setIDs = append(setIDs, currSetID)
			sets = append(sets, currSet)
			// Create new set.
			currSetID = setID
			currSet = make([]int, 0)
		}
		currSet = append(currSet, token)
	}
	if err := scanner.Err(); err != nil {
		return nil, nil, err
	}
	// Append the last set.
	setIDs = append(setIDs, currSetID)
	sets = append(sets, currSet)
	return setIDs, sets, nil
}
