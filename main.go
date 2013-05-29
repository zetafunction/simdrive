package main

import (
	"fmt"
	"github.com/zetafunction/simdrive/simulator"
	"math/rand"
	"time"
)

func runTrialAndReportStats(label string, runTrial func() int, numberOfTrials int) {
	totalHours := 0
	for i := 0; i < numberOfTrials; i++ {
		totalHours += runTrial()
	}
	fmt.Printf("%s ran %d times. Average lifetime: %f\n",
		label, numberOfTrials, float64(totalHours)/float64(numberOfTrials))
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	runTrialAndReportStats(
		"Single HDD",
		func() int {
			age := 0
			hdd := simulator.NewHardDiskDrive()
			for hdd.State() != simulator.FAILED {
				age++
				hdd.Step()
			}
			return age
		},
		1000)
	runTrialAndReportStats(
		"Striped storage pool with 2 disks",
		func() int {
			age := 0
			pool := simulator.NewStripedPool(2, simulator.NewHardDiskDrive)
			for pool.State() != simulator.FAILED {
				age++
				pool.Step()
			}
			return age
		},
		1000)
	runTrialAndReportStats(
		"Mirrored storage pool with 2 disks",
		func() int {
			age := 0
			pool := simulator.NewMirroredPool(2, simulator.NewHardDiskDrive)
			for pool.State() != simulator.FAILED {
				age++
				pool.Step()
			}
			return age
		},
		1000)
	runTrialAndReportStats(
		"RAIDZ-1 2+1 storage pool",
		func() int {
			age := 0
			pool := simulator.NewParityPool(3, 1, simulator.NewHardDiskDrive)
			for pool.State() != simulator.FAILED {
				age++
				pool.Step()
			}
			return age
		},
		1000)
	runTrialAndReportStats(
		"RAIDZ-3 8+3 storage pool",
		func() int {
			age := 0
			pool := simulator.NewParityPool(11, 3, simulator.NewHardDiskDrive)
			for pool.State() != simulator.FAILED {
				age++
				pool.Step()
			}
			return age
		},
		1000)
	runTrialAndReportStats(
		"Mirrored storage pool with 2x RAIDZ-2 4+2 storage pools",
		func() int {
			age := 0
			pool := simulator.NewStripedPool(2, func() simulator.Drive {
				return simulator.NewParityPool(4, 2, simulator.NewHardDiskDrive)
			})
			for pool.State() != simulator.FAILED {
				age++
				pool.Step()
			}
			return age
		},
		1000)
}
