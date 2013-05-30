package simulator

type storagePool struct {
	drives     []Drive
	state      DriveState
	redundancy int
	failures   int
}

type DriveCreator func() Drive

func newStoragePool(drives int, redundancy int, newDrive DriveCreator) Drive {
	pool := &storagePool{
		drives:     make([]Drive, drives),
		state:      OK,
		redundancy: redundancy,
		failures:   0,
	}
	for i, _ := range pool.drives {
		pool.drives[i] = newDrive()
	}
	return pool
}

func NewStripedPool(drives int, newDrive DriveCreator) Drive {
	return newStoragePool(drives, 0, newDrive)
}

func NewMirroredPool(drives int, newDrive DriveCreator) Drive {
	return newStoragePool(drives, drives-1, newDrive)
}

func NewParityPool(drives int, redundancy int, newDrive DriveCreator) Drive {
	return newStoragePool(drives, redundancy, newDrive)
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
