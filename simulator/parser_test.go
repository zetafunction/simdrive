package simulator

import (
	"fmt"
	"math/rand"
	"reflect"
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
	"sda": {"kind": "hard_disk"}}`,
	NewMirroredPool([]Drive{
		NewHardDiskDrive(testPrng),
		NewHardDiskDrive(testPrng),
	}),
	nil,
}, {
	`{"root": {"kind": "striped_pool", "drives": ["sda", "sda"]},
	"sda": {"kind": "hard_disk"}}`,
	NewStripedPool([]Drive{
		NewHardDiskDrive(testPrng),
		NewHardDiskDrive(testPrng),
	}),
	nil,
}}

func TestParseConfig(t *testing.T) {
	for _, test := range parseConfigTests {
		out, err := ParseConfig([]byte(test.in), testPrng)
		if !reflect.DeepEqual(test.out, out) {
			t.Errorf("Expected output: %+v, actual output: %+v", test.out, out)
		}
		if !reflect.DeepEqual(test.err, err) {
			t.Errorf("Expected error: %+v, actual error: %+v", test.err, err)
		}
	}
}
