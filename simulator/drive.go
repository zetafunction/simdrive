package simulator

type DriveState int

const (
	OK DriveState = iota
	FAILED
)

type Drive interface {
	// Steps the simulation forward by an hour. A basic drive implementation
	// that represents a traditional spinning, magnetic HDD might step its
	// age forward by an hour and roll a die to see if it enters a failure
	// state. A more complex drive implementation that represents a RAID-5
	// storage volume would invoke Step() on each disk it owns.
	Step()

	// Returns the current drive state.
	State() DriveState
}
