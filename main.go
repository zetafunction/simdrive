package main

import (
	"fmt"
	"github.com/zetafunction/simdrive/simulator"
	"math"
	"runtime"
	"sort"
)

// TODO: This code should belong in a separate package.
type Task func(chan int)

func runTrialAndReportStats(label string, task Task, iterations int) {
	c := make(chan int)
	for i := 0; i < iterations; i++ {
		go task(c)
	}
	results := make([]int, iterations)
	for i := 0; i < iterations; i++ {
		results[i] = <-c
	}
	sort.IntSlice(results).Sort()
	fmt.Printf("Results for %d iterations of %s:\n", iterations, label)

	total := 0
	for _, result := range results {
		total += result
	}
	mean := float64(total) / float64(iterations)
	fmt.Printf("  mean: %f\n", mean)

	squared_difference_sum := 0.
	for _, result := range results {
		squared_difference_sum += math.Pow(float64(result)-mean, 2)
	}
	variance := squared_difference_sum / float64(iterations)
	fmt.Printf("  stddev: %f\n", math.Sqrt(variance))

	count := 0
	for i := 1; i <= 10; i++ {
		for results[count] < i*8760 {
			count++
		}
		fmt.Printf("  %d year survival rate: %f\n",
			i, float64(iterations-count)/float64(iterations))
	}
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	runTrialAndReportStats(
		"Single HDD",
		func(c chan int) {
			age := 0
			hdd := simulator.NewHardDiskDrive()
			for hdd.State() != simulator.FAILED {
				age++
				hdd.Step()
			}
			c <- age
		},
		1000)
	runTrialAndReportStats(
		"Striped storage pool with 2 disks",
		func(c chan int) {
			age := 0
			pool := simulator.NewStripedPool(2, simulator.NewHardDiskDrive)
			for pool.State() != simulator.FAILED {
				age++
				pool.Step()
			}
			c <- age
		},
		1000)
	runTrialAndReportStats(
		"Mirrored storage pool with 2 disks",
		func(c chan int) {
			age := 0
			pool := simulator.NewMirroredPool(2, simulator.NewHardDiskDrive)
			for pool.State() != simulator.FAILED {
				age++
				pool.Step()
			}
			c <- age
		},
		1000)
	runTrialAndReportStats(
		"RAIDZ-1 2+1 storage pool",
		func(c chan int) {
			age := 0
			pool := simulator.NewParityPool(3, 1, simulator.NewHardDiskDrive)
			for pool.State() != simulator.FAILED {
				age++
				pool.Step()
			}
			c <- age
		},
		1000)
	runTrialAndReportStats(
		"RAIDZ-3 8+3 storage pool",
		func(c chan int) {
			age := 0
			pool := simulator.NewParityPool(11, 3, simulator.NewHardDiskDrive)
			for pool.State() != simulator.FAILED {
				age++
				pool.Step()
			}
			c <- age
		},
		1000)
	runTrialAndReportStats(
		"Mirrored storage pool with 2x RAIDZ-2 4+2 storage pools",
		func(c chan int) {
			age := 0
			pool := simulator.NewStripedPool(2, func() simulator.Drive {
				return simulator.NewParityPool(4, 2, simulator.NewHardDiskDrive)
			})
			for pool.State() != simulator.FAILED {
				age++
				pool.Step()
			}
			c <- age
		},
		1000)
}
