package main

import (
	"flag"
	"github.com/zetafunction/simdrive/simulator"
	"io/ioutil"
	"log"
	"math"
	"math/rand"
	"os"
	"runtime"
	"sort"
	"time"
)

var iterations = flag.Int("n", 1000, "number of iterations")

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	flag.Parse()

	args := flag.Args()
	if len(args) == 0 {
		log.Fatal("no input files")
	}
	for _, file := range args {
		log.Printf("Processing %s", file)

		f, err := os.Open(file)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()

		bytes, err := ioutil.ReadAll(f)
		if err != nil {
			log.Fatal(err)
		}

		start := time.Now()

		runAndReportResults(func(c chan<- int) {
			rng := rand.New(rand.NewSource(time.Now().UnixNano()))
			// TODO: Ideally,  there should be a way to clone a
			// Drive hierarchy rather than reparsing the input in
			// each goroutine.
			drive, err := simulator.ParseConfig(bytes, rng)
			if err != nil {
				log.Fatal(err)
			}
			age := 0
			for drive.Status() != simulator.FAILED {
				age++
				drive.Step()
			}
			c <- age
		}, *iterations)

		end := time.Now()
		elapsed := end.Sub(start)
		log.Printf("%s elapsed", elapsed)
	}
}

// TODO: This code should belong in a separate package.
func runAndReportResults(task func(chan<- int), iterations int) {
	c := make(chan int)
	for i := 0; i < iterations; i++ {
		go task(c)
	}
	results := make([]int, iterations)
	for i := 0; i < len(results); i++ {
		results[i] = <-c
	}
	sort.Ints(results)
	log.Printf("Results of %d iterations:", len(results))

	total := 0
	for _, result := range results {
		total += result
	}
	mean := float64(total) / float64(len(results))
	log.Printf("  mean: %f", mean)

	sumOfResidualsSquared := 0.
	for _, result := range results {
		sumOfResidualsSquared += math.Pow(float64(result)-mean, 2)
	}
	variance := sumOfResidualsSquared / float64(len(results))
	log.Printf("  stddev: %f", math.Sqrt(variance))

	count := 0
	for i := 1; i <= 10; i++ {
		for count < len(results) && results[count] < i*8760 {
			count++
		}
		log.Printf("  %d year survival rate: %f",
			i, float64(len(results)-count)/float64(len(results)))
	}
}
