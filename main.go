package main

import (
	"fmt"
	"github.com/zetafunction/disksim/simulator"
	"math/rand"
	"time"
)

func runTrialAndReportStats(
	label string, runTrial func() int, numberOfTrials int) {
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
		"Single hard disk",
		func() int {
			age := 0
			disk := simulator.NewHardDisk()
			for disk.State() != simulator.FAILED {
				age++
				disk.Step()
			}
			return age
		},
		1000)
	runTrialAndReportStats(
		"Striped disk (2)",
		func() int {
			age := 0
			disk := simulator.NewStripedDisk(2, simulator.NewHardDisk)
			for disk.State() != simulator.FAILED {
				age++
				disk.Step()
			}
			return age
		},
		1000)
	runTrialAndReportStats(
		"Mirrored disk (2)",
		func() int {
			age := 0
			disk := simulator.NewMirroredDisk(simulator.NewHardDisk)
			for disk.State() != simulator.FAILED {
				age++
				disk.Step()
			}
			return age
		},
		1000)
}
