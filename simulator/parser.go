package simulator

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"unicode"
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
	// Hard disk-specific parameters.
	Size       string
	Throughput string
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
		size, err := parseScaledUint(node.Size, "B")
		if err != nil {
			return nil, err
		}
		throughput, err := parseScaledUint(node.Throughput, "bps")
		drive = NewHardDisk(size, throughput, prng)
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

var siPrefixToScale = map[string]uint64{
	"":   1,
	"K":  1000,
	"M":  1000000,
	"G":  1000000000,
	"T":  1000000000000,
	"P":  1000000000000000,
	"E":  1000000000000000000,
	"Ki": 1 << 10,
	"Mi": 1 << 20,
	"Gi": 1 << 30,
	"Ti": 1 << 40,
	"Pi": 1 << 50,
	"Ei": 1 << 60,
}

// parseScaledUint parses strings that contain a number, an optional SI binary
// or metric prefix, and a unit specifier. If an SI prefix is found, the
// returned value n is scaled by the corresponding amount.
func parseScaledUint(input string, unit string) (n uint64, err error) {
	unitIndex := strings.Index(input, unit)
	if unitIndex < 0 || input[unitIndex:] != unit {
		return 0, fmt.Errorf("%s doesn't end with unit %s", input, unit)
	}

	afterDigitsIndex := strings.LastIndexFunc(input, unicode.IsDigit) + 1
	n, err = strconv.ParseUint(input[:afterDigitsIndex], 10, 64)
	if err != nil {
		return
	}

	if scale, present := siPrefixToScale[input[afterDigitsIndex:unitIndex]]; present {
		n *= scale
	} else {
		err = fmt.Errorf("invalid SI prefix in %s", input)
	}
	return
}
