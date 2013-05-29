package simulator

type stripedDisk struct {
	disks []Disk
	state DiskState
}

func NewStripedDisk(numberOfDisks int, diskCreationFunction func() Disk) Disk {
	disk := new(stripedDisk)
	disk.disks = make([]Disk, numberOfDisks)
	for i, _ := range disk.disks {
		disk.disks[i] = diskCreationFunction()
	}
	disk.state = OK
	return disk
}

func (disk *stripedDisk) Step() {
	for i, _ := range disk.disks {
		disk.disks[i].Step()
		if disk.disks[i].State() == FAILED {
			disk.state = FAILED
		}
	}
}

func (disk *stripedDisk) State() DiskState {
	return disk.state
}
