package simulator

import (
	"encoding/json"
	"fmt"
	"math/rand"
)

var (
	NoRootDriveError        = fmt.Errorf("No root drive found")
	MultipleRootDrivesError = fmt.Errorf("Multiple root drives found")
	CycleDetectedError      = fmt.Errorf("Detected cycle in config")
)

type driveNode struct {
	Kind   string
	Drives []string
	// Parity pool-specific parameters.
	Redundancy int
}

type driveGraph map[string]*driveNode
type driveNodeSet map[*driveNode]bool

// Parses JSON configs into an actual Drive configuration.
// TODO: Document this better.
func ParseConfig(config []byte, prng *rand.Rand) (Drive, error) {
	var nodes driveGraph
	if err := json.Unmarshal(config, &nodes); err != nil {
		return nil, err
	}

	root, err := findRoot(nodes)
	if err != nil {
		return nil, err
	}

	return generateDrive(nodes, root, driveNodeSet{}, prng)
}

func findRoot(nodes driveGraph) (*driveNode, error) {
	nonSourceNodes := map[string]bool{}
	for _, node := range nodes {
		for _, driveName := range node.Drives {
			nonSourceNodes[driveName] = true
		}
	}

	candidates := []*driveNode{}
	for name, node := range nodes {
		if _, present := nonSourceNodes[name]; !present {
			candidates = append(candidates, node)
		}
	}
	switch {
	case len(candidates) < 1:
		return nil, NoRootDriveError
	case len(candidates) > 1:
		return nil, MultipleRootDrivesError
	}

	return candidates[0], nil
}

func generateDrive(nodes driveGraph, node *driveNode, seen driveNodeSet, prng *rand.Rand) (Drive, error) {
	if _, present := seen[node]; present {
		return nil, CycleDetectedError
	}
	seen[node] = true
	defer delete(seen, node)

	drives := []Drive{}
	for _, driveName := range node.Drives {
		drive, err := generateDrive(nodes, nodes[driveName], seen, prng)
		if err != nil {
			return nil, err
		}
		drives = append(drives, drive)
	}
	switch node.Kind {
	case "hard_disk":
		return NewHardDiskDrive(prng), nil
	case "mirrored_pool":
		return NewMirroredPool(drives), nil
	case "parity_pool":
		// TODO: Return an error if redundancy is not set.
		return NewParityPool(drives, node.Redundancy), nil
	case "striped_pool":
		return NewStripedPool(drives), nil
	}
	return nil, fmt.Errorf("Invalid kind: %s", node.Kind)
}