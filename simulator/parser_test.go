package simulator

import (
	"fmt"
	"math/rand"
	"reflect"
	"strconv"
	"testing"
	"time"
)

var testPrng = rand.New(rand.NewSource(time.Now().UnixNano()))
var parseConfigTests = []struct {
	in  string
	out Drive
	err error
}{{"{}",
	nil,
	NoRootDriveError,
}, {
	`{"sda": {"kind": "hard_disk"}, "sdb": {"kind": "hard_disk"}}`,
	nil,
	MultipleRootDrivesError,
}, {
	`{"root": {"kind": "mirrored_pool", "drives": ["one", "one"]},
        "one": {"kind": "mirrored_pool", "drives": ["two", "two"]},
        "two": {"kind": "mirrored_pool", "drives": ["one", "one"]}}`,
	nil,
	CycleDetectedError,
}, {
	`{"invalid": {"kind": "invalid"}}`,
	nil,
	fmt.Errorf("Invalid kind: invalid"),
}, {
	`{"root": {"kind": "mirrored_pool", "drives": ["sda", "sda"]},
        "sda": {"kind": "hard_disk", "size": "1KB"}}`,
	NewMirroredPool([]Drive{
		NewHardDiskDrive(1000, testPrng),
		NewHardDiskDrive(1000, testPrng),
	}),
	nil,
}, {
	`{"root": {"kind": "parity_pool", "drives": ["sda", "sda", "sda"], "redundancy": 1},
        "sda": {"kind": "hard_disk", "size": "1KB"}}`,
	NewParityPool([]Drive{
		NewHardDiskDrive(1000, testPrng),
		NewHardDiskDrive(1000, testPrng),
		NewHardDiskDrive(1000, testPrng),
	}, 1),
	nil,
}, {
	`{"root": {"kind": "striped_pool", "drives": ["sda", "sda"]},
        "sda": {"kind": "hard_disk", "size": "1KB"}}`,
	NewStripedPool([]Drive{
		NewHardDiskDrive(1000, testPrng),
		NewHardDiskDrive(1000, testPrng),
	}),
	nil,
}}

func TestParseConfig(t *testing.T) {
	for _, test := range parseConfigTests {
		out, err := ParseConfig([]byte(test.in), testPrng)
		if !reflect.DeepEqual(test.out, out) {
			t.Errorf("Expected out: %+v, actual out: %+v", test.out, out)
		}
		if !reflect.DeepEqual(test.err, err) {
			t.Errorf("Expected err: %+v, actual err: %+v", test.err, err)
		}
	}
}

var parseScaledUintTests = []struct {
	input string
	unit  string
	n     uint64
	err   error
}{
	{"", "B", 0, fmt.Errorf(" doesn't end with unit B")},
	{"B", "B", 0, &strconv.NumError{"ParseUint", "", strconv.ErrSyntax}},
	{"-11B", "B", 0, &strconv.NumError{"ParseUint", "-11", strconv.ErrSyntax}},
	{"11ZB", "B", 11, fmt.Errorf("invalid SI prefix in 11ZB")},
	{"20B", "B", 20, nil},
	{"30KB", "B", 30000, nil},
	{"40MB", "B", 40000000, nil},
	{"50GB", "B", 50000000000, nil},
	{"60TB", "B", 60000000000000, nil},
	{"70PB", "B", 70000000000000000, nil},
	{"10EB", "B", 10000000000000000000, nil},
	{"25KiB", "B", 25 * (1 << 10), nil},
	{"35MiB", "B", 35 * (1 << 20), nil},
	{"45GiB", "B", 45 * (1 << 30), nil},
	{"55TiB", "B", 55 * (1 << 40), nil},
	{"65PiB", "B", 65 * (1 << 50), nil},
	{"15EiB", "B", 15 * (1 << 60), nil},
}

func TestParseScaledUint(t *testing.T) {
	for _, test := range parseScaledUintTests {
		n, err := parseScaledUint(test.input, test.unit)
		if test.n != n {
			t.Errorf("Expected n: %d, actual n: %d", test.n, n)
		}
		if !reflect.DeepEqual(test.err, err) {
			t.Errorf("Expected err: %v, actual err: %v", test.err, err)
		}
	}
}
