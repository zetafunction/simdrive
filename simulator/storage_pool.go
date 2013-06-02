package simulator

type storagePool struct {
	drives     []Drive
	state      DriveState
	redundancy int
	failures   int
}

type DriveCreator func() Drive

func newStoragePool(drives []Drive, redundancy int) Drive {
	// TODO: Error if redundancy is >= len(drives)
	pool := &storagePool{
		drives:     drives,
		state:      OK,
		redundancy: redundancy,
		failures:   0,
	}
	return pool
}

func NewStripedPool(drives []Drive) Drive {
	return newStoragePool(drives, 0)
}

func NewMirroredPool(drives []Drive) Drive {
	// TODO: Error if the drives are not all the same capacity?
	return newStoragePool(drives, len(drives)-1)
}

func NewParityPool(drives []Drive, redundancy int) Drive {
	// TODO: Error if the drives are not all the same capacity?
	return newStoragePool(drives, redundancy)
}

func (pool *storagePool) Step() {
	// TODO: Throw an error if disk.state == FAILED already.
	for _, drive := range pool.drives {
		if drive.State() == FAILED {
			continue
		}
		drive.Step()
		if drive.State() == FAILED {
			pool.failures++
			if pool.failures > pool.redundancy {
				pool.state = FAILED
			} else {
				// TODO: Attempt repair.
				pool.state = DEGRADED
			}
		}
	}
}

func (pool *storagePool) State() DriveState {
	return pool.state
}
