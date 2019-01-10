package SetSimilaritySearch

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// IntersectionSize computes the number of overlaps of two transformed sets
// (sorted integers).
func intersectionSize(s1, s2 []int) int {
	var i, j int
	var overlap int
	for i < len(s1) && j < len(s2) {
		switch d := s1[i] - s2[j]; {
		case d == 0:
			overlap++
			i++
			j++
		case d < 0:
			i++
		case d > 0:
			j++
		}
	}
	return overlap
}

type function func([]int, []int) float64

// Jaccard computes the Jaccard similarity of two transformed sets.
func jaccard(s1, s2 []int) float64 {
	if len(s1) == 0 && len(s2) == 0 {
		return 0.0
	}
	intersectionSize := intersectionSize(s1, s2)
	return float64(intersectionSize) / float64(len(s1)+len(s2)-intersectionSize)
}

// Containment computes the Containment of s1 in s2 -- the fraction of s1
// being found in s2.
func containment(s1, s2 []int) float64 {
	if len(s1) == 0 {
		return 0.0
	}
	intersectionSize := intersectionSize(s1, s2)
	return float64(intersectionSize) / float64(len(s1))
}

type overlapThresholdFunction func(int, float64) int

// x is the set size
// t is the Jaccard threshold
func jaccardOverlapThresholdFunc(x int, t float64) int {
	return max(1, int(float64(x)*t))
}

var jaccardOverlapIndexThresholdFunc = jaccardOverlapThresholdFunc

// This is used for query only.
func containmentOverlapThresholdFunc(x int, t float64) int {
	return max(1, int(float64(x)*t))
}

func containmentOverlapIndexThresholdFunc(x int, t float64) int {
	return 1
}

type positionFilter func([]int, []int, int, int, float64) bool

func jaccardPositionFilter(s1, s2 []int, p1, p2 int, t float64) bool {
	l1, l2 := len(s1), len(s2)
	return float64(min(l1-p1, l2-p2))/float64(max(l1, l2)) >= t
}

func containmentPositionFilter(s1, s2 []int, p1, p2 int, t float64) bool {
	l1, l2 := len(s1), len(s2)
	return float64(min(l1-p1, l2-p2))/float64(l1) >= t
}

var similarityFuncs = map[string]function{
	"jaccard":     jaccard,
	"containment": containment,
}

var overlapThresholdFuncs = map[string]overlapThresholdFunction{
	"jaccard":     jaccardOverlapThresholdFunc,
	"containment": containmentOverlapThresholdFunc,
}

var overlapIndexThresholdFuncs = map[string]overlapThresholdFunction{
	"jaccard":     jaccardOverlapIndexThresholdFunc,
	"containment": containmentOverlapIndexThresholdFunc,
}

var positionFilterFuncs = map[string]positionFilter{
	"jaccard":     jaccardPositionFilter,
	"containment": containmentPositionFilter,
}
