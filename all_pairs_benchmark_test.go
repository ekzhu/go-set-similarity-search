package SetSimilaritySearch

import (
	"compress/gzip"
	"encoding/csv"
	"log"
	"os"
	"strconv"
	"testing"
	"time"
)

var (
	// Download from https://github.com/ekzhu/set-similarity-search-benchmarks
	allPairOpenDataBenchmarkFilename  = "canada_us_uk_opendata.inp.gz"
	allPairOpenDataBenchmarkResult    = "canada_us_uk_opendata_all_pairs.csv"
	allPairOpenDataBenchmarkThreshold = 0.2
	allPairOpenDataBenchmarkMinSize   = 10
	allPairOpenDataBenchmarkFunction  = "jaccard"
)

var (
	// Download from https://github.com/ekzhu/set-similarity-search-benchmarks
	allPairBmsPosBenchmarkFilename  = "BMS-POS_dup_dr.inp.gz"
	allPairBmsPosBenchmarkResult    = "BMS-POS_dup_dr_all_pairs.csv"
	allPairBmsPosBenchmarkThreshold = 0.2
	allPairBmsPosBenchmarkMinSize   = 1
	allPairBmsPosBenchmarkFunction  = "jaccard"
)

var (
	// Download from https://snap.stanford.edu/data/soc-Pokec.html
	allPairsPokecBenchmarkFilename      = "soc-pokec-relationships.txt.gz"
	allPairsPokecBenchmarkResult        = "soc-pokec-relationships-all-pairs.csv"
	allPairsPokecBenchmarkThreshold     = 0.5
	allPairsPokecBenchmarkFunction      = "jaccard"
	allPairsPokecBenchmarkInputReversed = false
)

var (
	// Download from https://snap.stanford.edu/data/soc-LiveJournal1.html
	allPairsLiveJournalBenchmarkFilename      = "soc-LiveJournal1.txt.gz"
	allPairsLiveJournalBenchmarkResult        = "soc-LiveJournal1-all-pairs.csv"
	allPairsLiveJournalBenchmarkFunction      = "jaccard"
	allPairsLiveJournalBenchmarkThreshold     = 0.5
	allPairsLiveJournalBenchmarkInputReversed = false
)

func BenchmarkOpenDataAllPair(b *testing.B) {
	benchmarkAllPairRowFile(b, allPairOpenDataBenchmarkFilename,
		allPairOpenDataBenchmarkResult, allPairOpenDataBenchmarkFunction,
		allPairOpenDataBenchmarkThreshold, allPairOpenDataBenchmarkMinSize)
}

func BenchmarkBmsPosAllPair(b *testing.B) {
	benchmarkAllPairRowFile(b, allPairBmsPosBenchmarkFilename,
		allPairBmsPosBenchmarkResult, allPairBmsPosBenchmarkFunction,
		allPairBmsPosBenchmarkThreshold, allPairBmsPosBenchmarkMinSize)
}

func benchmarkAllPairRowFile(b *testing.B, benchmarkFile, resultFile, function string, threshold float64, minSize int) {
	log.Printf("Reading transformed sets from %s",
		benchmarkFile)
	start := time.Now()
	sets := readGzippedTransformedSets(benchmarkFile,
		/*firstLineInfo=*/ true,
		minSize)
	log.Printf("Finished reading %d transformed sets in %s", len(sets),
		time.Now().Sub(start).String())
	log.Printf("Running AllPairs algorithm")
	out, err := os.Create(resultFile)
	if err != nil {
		b.Fatal(err)
	}
	defer out.Close()
	w := csv.NewWriter(out)
	start = time.Now()
	pairs, err := AllPairs(sets, function,
		threshold)
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
	log.Printf("Results written to %s", resultFile)
}

func BenchmarkPokecAllPair(b *testing.B) {
	benchmarkUseGzippedFlattendRawSets(
		allPairsPokecBenchmarkFilename,
		allPairsPokecBenchmarkResult,
		allPairsPokecBenchmarkFunction,
		allPairsPokecBenchmarkThreshold,
		allPairsPokecBenchmarkInputReversed,
		b,
	)
}

func BenchmarkLiveJournalAllPair(b *testing.B) {
	benchmarkUseGzippedFlattendRawSets(
		allPairsLiveJournalBenchmarkFilename,
		allPairsLiveJournalBenchmarkResult,
		allPairsLiveJournalBenchmarkFunction,
		allPairsLiveJournalBenchmarkThreshold,
		allPairsLiveJournalBenchmarkInputReversed,
		b,
	)
}

func benchmarkUseGzippedFlattendRawSets(input string,
	output string, function string, threshold float64, inputReversed bool,
	b *testing.B) {

	log.Printf("Reading raw sets from %s", input)
	start := time.Now()
	file, err := os.Open(input)
	if err != nil {
		b.Fatal(err)
	}
	defer file.Close()
	gz, err := gzip.NewReader(file)
	if err != nil {
		panic(err)
	}
	defer gz.Close()
	setIDs, rawSets, err := ReadFlattenedRawSets(gz, inputReversed)
	if err != nil {
		b.Fatal(err)
	}
	log.Printf("Finished reading %d raw sets in %s", len(setIDs),
		time.Now().Sub(start).String())

	log.Printf("Transforming raw sets")
	start = time.Now()
	sets, _ := FrequencyOrderTransform(rawSets)
	// Remove unused data.
	rawSets = nil
	log.Printf("Finished transforming sets in %s",
		time.Now().Sub(start).String())

	log.Printf("Running AllPairs algorithm")
	out, err := os.Create(output)
	if err != nil {
		b.Fatal(err)
	}
	defer out.Close()
	w := csv.NewWriter(out)
	start = time.Now()
	pairs, err := AllPairs(sets, function, threshold)
	for pair := range pairs {
		w.Write([]string{
			setIDs[pair.X],
			setIDs[pair.Y],
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
	log.Printf("Results written to %s", output)
}
