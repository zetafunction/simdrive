package simulator

type mirroredDisk struct {
	disks        [2]Disk
	state        DiskState
	failureCount int
}

func NewMirroredDisk(diskCreationFunction func() Disk) Disk {
	disk := new(mirroredDisk)
	for i, _ := range disk.disks {
		disk.disks[i] = diskCreationFunction()
	}
	disk.state = OK
	return disk
}

func (disk *mirroredDisk) Step() {
	for i, _ := range disk.disks {
		if disk.disks[i].State() == FAILED {
			continue
		}
		disk.disks[i].Step()
		if disk.disks[i].State() == FAILED {
			disk.failureCount++
			if disk.failureCount > 1 {
				disk.state = FAILED
			}
		}
	}
}

func (disk *mirroredDisk) State() DiskState {
	return disk.state
}
