package simulator

import (
	"math"
	"math/rand"
)

// This HDD model is based on the annualized failure rates (AFR) from Google's
// 2007 study on hard disks (see http://goo.gl/sPqui). The classical formula for
// determining AFR from mean time between failures (MTBF) is:
// AFR = 1 - exp(-8760 / mean time before failure)
// The study provides AFRs for several generation buckets (this is estimated
// from the chart in figure 2), as actual numbers are not provided):
// 3 months: 0.03
// 6 months: 0.018
// 1 year: 0.017
// 2 years: 0.08
// 3 years: 0.86
// The paper notes that there may be some overlap between the 3 month, 6 month,
// and 1 year bucket, but that is currently ignored. These AFRs are used to
// calculate the MTBF (which are a bit easier to read and represent). A modified
// version of the AFR formula substituting 1 for 8760 is used to determine the
// hourly failure rate for actual simulation purposes.
type hardDiskDrive struct {
	// Age of the HDD, in hours.
	// TODO: Convert the representation to use go's time.Duration.
	age    int
	status DriveStatus
	prng   *rand.Rand
}

// NewHardDiskDrive returns a new Drive representing a tradtional disk drive
// using spinning magnetic media.
func NewHardDiskDrive(prng *rand.Rand) Drive {
	return &hardDiskDrive{
		age:    0,
		status: OK,
		prng:   prng,
	}
}

// Several helpers to help convert the data from annualized failure rates to
// hourly failure rates.
func annualizedFailureRateToMeanTimeBetweenFailures(x float64) float64 {
	return -8760 / math.Log(1-x)
}

func meanTimeBetweenFailuresToHourlyFailureRate(x float64) float64 {
	return 1 - math.Exp(-1/x)
}

func annualizedFailureRateToHourlyFailureRate(x float64) float64 {
	return meanTimeBetweenFailuresToHourlyFailureRate(
		annualizedFailureRateToMeanTimeBetweenFailures(x))

}

// TODO: These are really constants, so having to mark them as var is ugly.
var (
	threeMonthHourlyFailureRate = annualizedFailureRateToHourlyFailureRate(0.03)
	sixMonthHourlyFailureRate   = annualizedFailureRateToHourlyFailureRate(0.018)
	oneYearHourlyFailureRate    = annualizedFailureRateToHourlyFailureRate(0.017)
	twoYearHourlyFailureRate    = annualizedFailureRateToHourlyFailureRate(0.08)
	threeYearHourlyFailureRate  = annualizedFailureRateToHourlyFailureRate(0.086)
)

func hourlyFailureRateForAge(age int) float64 {
	switch {
	case age < 8760/4:
		return threeMonthHourlyFailureRate
	case age < 8760/2:
		return sixMonthHourlyFailureRate
	case age < 8760:
		return oneYearHourlyFailureRate
	case age < 8760*2:
		return twoYearHourlyFailureRate
	}
	return threeYearHourlyFailureRate
}

func (hdd *hardDiskDrive) Step() {
	hdd.age++
	chance := hourlyFailureRateForAge(hdd.age)
	if hdd.prng.Float64() < chance {
		hdd.status = FAILED
	}
}

func (hdd *hardDiskDrive) Status() DriveStatus {
	return hdd.status
}
