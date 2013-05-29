package simulator

type DiskState int

const (
	OK = iota
	FAILED
)

type Disk interface {
	// Steps the simulation forward by an hour. A basic disk implementation
	// that represents a traditional spinning, magnetic hard disk might step
	// its age forward by an hour and roll a die to see if it enters a
	// failure state. A more complex disk implementation that represents a
	// RAID-5 storage volume would invoke Step() on each disk it owns.
	Step()

	// Returns the current disk state.
	State() DiskState
}
