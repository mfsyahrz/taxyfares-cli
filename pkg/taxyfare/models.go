package taxyfare

import (
	"fmt"
	"time"
)

const (
	baseFare           float64 = 400
	fareUpto10KM       float64 = 40
	fareAbove10KM      float64 = 40
	distCountUpTo10KM  float64 = 400
	distCountAbove10KM float64 = 350

	distance1KM = 1000
	distance9KM = 9000

	maxTimeBetweenRecords = 5 * time.Minute
	minTripRecords        = 2
	distance0             = 0.0
)

// Record represents a taxi trip record.
type Record struct {
	DistanceMeter         // distance meter information
	mileAgeDiff   float64 // mileage difference
}

// String returns a formatted string representing the record value.
func (t Record) String() string {
	return fmt.Sprintf("%s %.1f %.1f", t.elapsedTime, t.distance, t.mileAgeDiff)
}
