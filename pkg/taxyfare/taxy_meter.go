package taxyfare

import (
	"fmt"
	"sort"
	"time"

	"github.com/mfsyahrz/taxyfare-cli/pkg/utils/logger"
)

// TaxyMeter is an interface for calculating taxi fare based on recorded data
type TaxyMeter interface {
	//adds a new record to the meter with the given string input
	AddRecord(record string) error

	// GetCurrentFare calculates and returns the current fare based on the meter's records
	GetCurrentFare() (int, error)

	// GetRecordHistory returns a string representation of the meter's recorded history
	GetRecordHistory() string
}

type manager struct {
	DistanceMeter
	fare    float64   // The current fare of the trip
	records []*Record // trip records history
}

func New() TaxyMeter {
	return &manager{}
}

// AddRecord takes a string as input, parses it into a DistanceMeter,
// validates it, and adds it to the list of records.
// it will update the fare based on the current distance.
func (t *manager) AddRecord(record string) error {
	logger.Info("Start Adding new trip with record: ", record)

	dm, err := ParseDistanceMeter(record)
	if err != nil {
		return err
	}

	if len(t.records) == 0 {
		t.processRecord(dm, 0)
		return nil
	}

	if err := t.validateToLastRecord(dm); err != nil {
		return err
	}

	t.processRecord(dm, dm.distance-t.distance)
	return nil
}

// GetCurrentFare returns the current fare in integer based on the recorded trips.
func (t *manager) GetCurrentFare() (int, error) {
	if len(t.records) < minTripRecords {
		return 0, fmt.Errorf("unable to get fare, trip records is lower than %d", minTripRecords)
	}

	if t.distance == distance0 {
		return 0, fmt.Errorf("unable to get fare, total mileage is %.1f m", distance0)
	}

	return int(t.fare), nil
}

// GetRecordHistory returns a string containing a history of all the recorded trips
// sorted in descending order by the difference in mileage between the trip.
func (t *manager) GetRecordHistory() string {
	trips := t.records

	// sort the trip records by the difference in mileage since the last record in descending order
	sort.Slice(trips, func(i, j int) bool {
		return trips[i].mileAgeDiff > trips[j].mileAgeDiff
	})

	var output string
	for _, trip := range trips {
		output += fmt.Sprintf("%s\n", trip)
	}

	return output
}

func (t *manager) validateToLastRecord(dm *DistanceMeter) error {
	if dm.elapsedTime.Before(t.elapsedTime) {
		return fmt.Errorf("elapsed time [%s] is earlier than last elapsed time [%s]", dm.elapsedTime, t.elapsedTime)
	}

	if dm.elapsedTime.Compare(t.elapsedTime) > 5*time.Minute {
		return fmt.Errorf("elapsed time is aparted more than %s from last elapsed time", maxTimeBetweenRecords)
	}

	if dm.distance < t.distance {
		return fmt.Errorf("distance [%.1f] is smaller than last distance [%.1f]", dm.distance, t.distance)
	}

	return nil
}

func (t *manager) processRecord(dm *DistanceMeter, mileAgeDiff float64) {
	t.DistanceMeter = *dm
	t.updateFare()

	t.records = append(t.records, &Record{
		DistanceMeter: *dm,
		mileAgeDiff:   mileAgeDiff,
	})

	logger.Info(fmt.Sprintf("Current fare : %d ", int(t.fare)))

	fmt.Printf("Current Mileage PerMinute: %.1f KM/Minute\n", t.CalcMileagePerMinute()/1000)
}

// updateFare updates current fare following current distance value.
// this function works idempotently, it means the fare will not change
// if no change in the distance value prior calling this function
// fare rules are:
// 1. The base fare is 400 yen for up to 1 km.
// 2. Up to 10 km, 40 yen is added every 400 meters.
// 3. Over 10km, 40 yen is added every 350 meters.
func (t *manager) updateFare() {
	t.fare = baseFare

	distanceDiffFrom1km := t.distance - distance1KM
	if distanceDiffFrom1km <= 0 {
		return
	}

	distanceDiffFrom9km := distanceDiffFrom1km - distance9KM
	if distanceDiffFrom9km <= 0 {
		t.fare += (distanceDiffFrom1km / distCountUpTo10KM) * fareUpto10KM
		return
	}

	t.fare += ((distance9KM / distCountUpTo10KM) * fareUpto10KM) + ((distanceDiffFrom9km / distCountAbove10KM) * fareAbove10KM)
}
