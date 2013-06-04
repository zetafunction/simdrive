package simulator

// DriveStatus represents the status of a drive.
type DriveStatus int

const (
	// OK means that the drive is healthy.
	OK DriveStatus = iota
	// DEGRADED means that the drive has experienced failures in a manner
	// that are recoverable.
	DEGRADED
	// FAILED means that the drive has lost data.
	FAILED
)

// A Drive represents a simulation of an abstract storage mechanism for data.
type Drive interface {
	// Steps the simulation forward by an hour. A basic drive implementation
	// that represents a traditional spinning, magnetic HDD might step its
	// age forward by an hour and roll a die to see if it enters a failure
	// state. A more complex drive implementation that represents a RAID-5
	// storage volume would invoke Step() on each disk it owns.
	Step()

	// Returns the current drive status.
	Status() DriveStatus
}
