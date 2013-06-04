package simulator

import (
	"encoding/json"
	"fmt"
	"math/rand"
)

var (
	// NoRootDriveError means that a candidate drive to use as the root of
	// the drive hierarchy was not found in the JSON config.
	NoRootDriveError = fmt.Errorf("No root drive found")
	// MultipleRootDrivesError means that multiple candidate drives for the
	// root of the drive hierarchy were found in the JSON config.
	MultipleRootDrivesError = fmt.Errorf("Multiple root drives found")
	// CycleDetectedError means that a cycle was found in the dependency
	// tree when creating the drive hierarchy.
	CycleDetectedError = fmt.Errorf("Detected cycle in config")
)

type driveNode struct {
	Kind   string
	Drives []string
	// Parity pool-specific parameters.
	Redundancy int
}

type driveGraph map[string]*driveNode
type driveNodeSet map[*driveNode]bool

// ParseConfig parses the JSON string in config into a Drive object.
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

func generateDrive(nodes driveGraph, node *driveNode, seen driveNodeSet, prng *rand.Rand) (drive Drive, err error) {
	if _, present := seen[node]; present {
		return nil, CycleDetectedError
	}
	seen[node] = true
	defer delete(seen, node)

	drives := []Drive{}
	for _, driveName := range node.Drives {
		d, err := generateDrive(nodes, nodes[driveName], seen, prng)
		if err != nil {
			return nil, err
		}
		drives = append(drives, d)
	}
	switch node.Kind {
	case "hard_disk":
		drive = NewHardDiskDrive(prng)
	case "mirrored_pool":
		drive = NewMirroredPool(drives)
	case "parity_pool":
		// TODO: Return an error if redundancy is not set.
		drive = NewParityPool(drives, node.Redundancy)
	case "striped_pool":
		drive = NewStripedPool(drives)
	default:
		err = fmt.Errorf("Invalid kind: %s", node.Kind)
	}
	return drive, err
}
