package main

import (
	"fmt"
	"github.com/zetafunction/simdrive/simulator"
	"runtime"
)

type Task func(chan int)

func runTrialAndReportStats(label string, task Task, iterations int) {
	c := make(chan int)
	for i := 0; i < iterations; i++ {
		go task(c)
	}
	totalHours := 0
	for i := 0; i < iterations; i++ {
		totalHours += <-c
	}
	fmt.Printf("%s ran %d times. Average lifetime: %f\n",
		label, iterations, float64(totalHours)/float64(iterations))
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
