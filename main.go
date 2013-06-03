package main

import (
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

// TODO: This code should belong in a separate package.
type Task func(chan<- int)

func runAndReportResults(task Task, iterations int) {
	c := make(chan int)
	for i := 0; i < iterations; i++ {
		go task(c)
	}
	results := make([]int, iterations)
	for i := 0; i < iterations; i++ {
		results[i] = <-c
	}
	sort.IntSlice(results).Sort()
	log.Printf("Results of %d iterations:", iterations)

	total := 0
	for _, result := range results {
		total += result
	}
	mean := float64(total) / float64(iterations)
	log.Printf("  mean: %f", mean)

	squared_difference_sum := 0.
	for _, result := range results {
		squared_difference_sum += math.Pow(float64(result)-mean, 2)
	}
	variance := squared_difference_sum / float64(iterations)
	log.Printf("  stddev: %f", math.Sqrt(variance))

	count := 0
	for i := 1; i <= 10; i++ {
		for results[count] < i*8760 {
			count++
		}
		log.Printf("  %d year survival rate: %f",
			i, float64(iterations-count)/float64(iterations))
	}
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	bytes, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}

	start := time.Now()

	// TODO: Configure the number of iterations with a command-line flag.
	runAndReportResults(func(c chan<- int) {
		rng := rand.New(rand.NewSource(time.Now().UnixNano()))
		// TODO: Ideally,  there should be a way to clone a Drive
		// hierarchy rather than reparsing the input in each goroutine.
		drive, err := simulator.ParseConfig(bytes, rng)
		if err != nil {
			log.Fatal(err)
		}
		age := 0
		for drive.State() != simulator.FAILED {
			age++
			drive.Step()
		}
		c <- age
	}, 1000)

	end := time.Now()
	elapsed := end.Sub(start)
	log.Printf("%s elapsed", elapsed)
}
