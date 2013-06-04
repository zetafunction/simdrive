package simulator

import (
	"math"
)

type storagePool struct {
	drives     []Drive
	status     DriveStatus
	redundancy int
	failures   int
}

func newStoragePool(drives []Drive, redundancy int) Drive {
	// TODO: Error if redundancy is >= len(drives)
	return &storagePool{
		drives:     drives,
		status:     OK,
		redundancy: redundancy,
		failures:   0,
	}
}

// NewStripedPool returns a new Drive that stripes data across its member
// drives.
func NewStripedPool(drives []Drive) Drive {
	return newStoragePool(drives, 0)
}

// NewMirroredPool returns a new Drive that mirrors data across its member
// drives.
func NewMirroredPool(drives []Drive) Drive {
	// TODO: Error if the drives are not all the same capacity?
	return newStoragePool(drives, len(drives)-1)
}

// NewParityPool returns a new Drive that stripes data across its member drives
// using len(drives)-redundancy of the drives for data and the remainder
// for parity.
func NewParityPool(drives []Drive, redundancy int) Drive {
	// TODO: Error if the drives are not all the same capacity?
	return newStoragePool(drives, redundancy)
}

func (pool *storagePool) Step() {
	// TODO: Throw an error if disk.status == FAILED already.
	for _, drive := range pool.drives {
		if drive.Status() == FAILED {
			continue
		}
		drive.Step()
		if drive.Status() == FAILED {
			pool.failures++
			if pool.failures > pool.redundancy {
				pool.status = FAILED
			} else {
				// TODO: Attempt repair.
				pool.status = DEGRADED
			}
		}
	}
}

func (pool *storagePool) Status() DriveStatus {
	return pool.status
}

func (pool *storagePool) Size() (n uint64) {
	// TODO: Consider calculating this once at drive creation instead.
	if pool.redundancy == 0 {
		// A striped pool's size is simply the sum of all its member
		// drives.
		for _, drive := range pool.drives {
			n += drive.Size()
		}
	} else {
		// The actual size of a mirrored/parity pool is determined by
		// the size of its smallest member drive.
		var minMemberDriveSize uint64 = math.MaxUint64
		for _, drive := range pool.drives {
			if drive.Size() < minMemberDriveSize {
				minMemberDriveSize = drive.Size()
			}
		}
		n = minMemberDriveSize * uint64(len(pool.drives)-pool.redundancy)
	}
	return
}

func (pool *storagePool) Throughput() uint64 {
	// TODO: To be more realistic, this should take into account:
	// - Throughput might not scale 100% linearly, especially if parity
	//   drives are being used.
	// - In the real world, physical interface limits might cap the actual
	//   throughput for a storage pool with lots of drives.
	var minThroughput uint64 = math.MaxUint64
	for _, drive := range pool.drives {
		if drive.Status() == FAILED {
			continue
		}
		if drive.Throughput() < minThroughput {
			minThroughput = drive.Throughput()
		}
	}
	return minThroughput * uint64(len(pool.drives)-pool.redundancy)
}
