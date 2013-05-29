package main

import (
	"fmt"
	"github.com/zetafunction/disksim/simulator"
	"math/rand"
	"time"
)

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	disk := simulator.NewHardDisk()
	for disk.State() != simulator.FAILED {
		disk.Step()
	}
	fmt.Printf("Disk failed at the age of %d hours.\n", disk.Age)
}
