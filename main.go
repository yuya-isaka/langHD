package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/yuya-isaka/langHD/hd"
	"github.com/yuya-isaka/langHD/myutil"
)

func main() {
	// Command line arguments
	var (
		cache = flag.Bool("cache", true, "cache flag")
		train = flag.String("train", "data/train", "training data path")
		test  = flag.String("test", "data/test", "test data path")
		ngram = flag.Int("ngram", 3, "n consective letters")
		dim   = flag.Int("dim", 10000, "dimension of hypervector")
	)
	flag.Parse()

	t := time.Now()

	// Making cache
	fmt.Println("\nMaking cache ...")
	myutil.MakeCacheDirectry(cache)

	hdc := hd.NewLangHD(dim, ngram)

	// Encoding
	fmt.Println("Encoding Asciis ...")
	hdc.EncodeAsciis(cache)

	fmt.Println("Encoding training data ...")
	hdc.EncodeTrainingData(cache, train)

	fmt.Println("Encoding testing data ...")
	hdc.EncodeTestingData(test)

	diff := time.Since(t)
	fmt.Println("\nFinish ", diff.Seconds(), " seconds")

	// Testing
	fmt.Println("\nTesting")
	hdc.Testing()
}
