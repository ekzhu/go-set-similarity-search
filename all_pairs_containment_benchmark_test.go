package SetSimilaritySearch

import (
	"bufio"
	"compress/gzip"
	"encoding/csv"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"
)

var (
	allPairsContainmentBenchmarkFilename  = "canada_us_uk_opendata.inp.gz"
	allPairsContainmentBenchmarkResult    = "canada_us_uk_opendata_all_pairs_containment.csv"
	allPairsContainmentBenchmarkThreshold = 0.9
	allPairsContainmentBenchmarkMinSize   = 10
)

func readGzippedTransformedSets(filename string,
	firstLineInfo bool, minSize int) (sets [][]int) {
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	gz, err := gzip.NewReader(file)
	if err != nil {
		panic(err)
	}
	defer gz.Close()
	sets = make([][]int, 0)
	scanner := bufio.NewScanner(gz)
	scanner.Buffer(nil, 1024*1024*1024*4)
	for scanner.Scan() {
		line := strings.Trim(scanner.Text(), "\n")
		if firstLineInfo && len(sets) == 0 {
			// Initialize the sets using the info given by the first line
			count, err := strconv.Atoi(strings.Split(line, " ")[0])
			if err != nil {
				panic(err)
			}
			sets = make([][]int, 0, count)
			firstLineInfo = false
			continue
		}
		raw := strings.Split(strings.Split(line, "\t")[1], ",")
		if len(raw) < minSize {
			continue
		}
		set := make([]int, len(raw))
		for i := range set {
			set[i], err = strconv.Atoi(raw[i])
			if err != nil {
				panic(err)
			}
		}
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}
	return sets
}

func BenchmarkAllPairContainment(b *testing.B) {
	b.Logf("Reading transformed sets from %s",
		allPairsContainmentBenchmarkFilename)
	start := time.Now()
	sets := readGzippedTransformedSets(allPairsContainmentBenchmarkFilename,
		/*firstLineInfo=*/ true,
		allPairsContainmentBenchmarkMinSize)
	b.Logf("Finished reading %d transformed sets in %s", len(sets),
		time.Now().Sub(start).String())
	b.Logf("Building search index")
	start = time.Now()
	searchIndex, err := NewSearchIndex(sets, "containment",
		allPairsContainmentBenchmarkThreshold)
	if err != nil {
		b.Fatal(err)
	}
	b.Logf("Finished building search index in %s",
		time.Now().Sub(start).String())
	out, err := os.Create(allPairsContainmentBenchmarkResult)
	if err != nil {
		b.Fatal(err)
	}
	defer out.Close()
	w := csv.NewWriter(out)
	b.Logf("Begin querying")
	start = time.Now()
	for i, set := range sets {
		results := searchIndex.Query(set)
		for _, result := range results {
			if result.X == i {
				continue
			}
			w.Write([]string{
				strconv.Itoa(i),
				strconv.Itoa(result.X),
				strconv.FormatFloat(result.Similarity, 'f', 4, 64),
			})
		}
	}
	b.Logf("Finished querying in %s", time.Now().Sub(start).String())
	w.Flush()
	if err := w.Error(); err != nil {
		b.Fatal(err)
	}
	b.Logf("Results written to %s", allPairsContainmentBenchmarkResult)
}
