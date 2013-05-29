package simulator

type storagePool struct {
	drives     []Drive
	state      DriveState
	redundancy int
	failures   int
}

type DriveCreator func() Drive

func newStoragePool(drives int, redundancy int, newDrive DriveCreator) Drive {
	pool := new(storagePool)
	pool.drives = make([]Drive, drives)
	for i, _ := range pool.drives {
		pool.drives[i] = newDrive()
	}
	pool.state = OK
	pool.redundancy = redundancy
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
	for i, _ := range pool.drives {
		if pool.drives[i].State() == FAILED {
			continue
		}
		pool.drives[i].Step()
		if pool.drives[i].State() == FAILED {
			pool.failures++
			if pool.failures > pool.redundancy {
				pool.state = FAILED
			}
		}
	}
}

func (pool *storagePool) State() DriveState {
	return pool.state
}
