package SetSimilaritySearch

import (
	"encoding/csv"
	"log"
	"os"
	"strconv"
	"testing"
	"time"
)

var (
	allPairBenchmarkFilename  = "canada_us_uk_opendata.inp.gz"
	allPairBenchmarkResult    = "canada_us_uk_opendata_all_pairs.csv"
	allPairBenchmarkThreshold = 0.9
	allPairBenchmarkMinSize   = 10
	allPairBenchmarkFunction  = "jaccard"
)

func BenchmarkAllPair(b *testing.B) {
	log.Printf("Reading transformed sets from %s",
		allPairBenchmarkFilename)
	start := time.Now()
	sets := readGzippedTransformedSets(allPairBenchmarkFilename,
		/*firstLineInfo=*/ true,
		allPairBenchmarkMinSize)
	log.Printf("Finished reading %d transformed sets in %s", len(sets),
		time.Now().Sub(start).String())
	log.Printf("Running AllPairs algorithm")
	out, err := os.Create(allPairBenchmarkResult)
	if err != nil {
		b.Fatal(err)
	}
	defer out.Close()
	w := csv.NewWriter(out)
	start = time.Now()
	pairs, err := AllPairs(sets, allPairBenchmarkFunction,
		allPairBenchmarkThreshold)
	for pair := range pairs {
		w.Write([]string{
			strconv.Itoa(pair.X),
			strconv.Itoa(pair.Y),
			strconv.FormatFloat(pair.Similarity, 'f', 4, 64),
		})
	}
	log.Printf("Finished AllPairs in %s", time.Now().Sub(start).String())
	w.Flush()
	if err := w.Error(); err != nil {
		b.Fatal(err)
	}
	if err != nil {
		b.Fatal(err)
	}
	log.Printf("Results written to %s", allPairBenchmarkResult)
}
