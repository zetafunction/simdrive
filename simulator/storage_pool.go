package simulator

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

// NewStripedPool returns a new Drive that stripes data across the provided
// drives.
func NewStripedPool(drives []Drive) Drive {
	return newStoragePool(drives, 0)
}

// NewMirroredPool returns a new Drive that mirrors data across the provided
// drives.
func NewMirroredPool(drives []Drive) Drive {
	// TODO: Error if the drives are not all the same capacity?
	return newStoragePool(drives, len(drives)-1)
}

// NewParityPool returns a new Drive that stripes data across the provided
// drives, using len(drives)-redundancy of the drives for data and the remainder
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
