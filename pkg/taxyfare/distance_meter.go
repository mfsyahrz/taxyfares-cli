package taxyfare

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/mfsyahrz/taxyfare-cli/pkg/utils/logger"
)

// DistanceMeter is a struct that holds the distance and elapsed time of a measurement.
type DistanceMeter struct {
	distance    float64
	elapsedTime ElapsTime
}

// CalcMileagePerMinute calculates and returns the mileage per minute based on the distance and elapsed time of a measurement.
// If the elapsed time is less than one minute, it returns 0.
func (d *DistanceMeter) CalcMileagePerMinute() float64 {
	if d.elapsedTime.Minutes.Minutes() < 1 {
		return 0
	}

	mileagePerMin := d.distance / d.elapsedTime.Minutes.Minutes()

	logger.Info(fmt.Sprintf("current mileage/min detail [distance = %.1f, elapsedtime = %.1f minute, mileagePerMin = %.1f m/min]",
		d.distance, float64(d.elapsedTime.Minutes.Minutes()), mileagePerMin))

	return mileagePerMin
}

// ParseDistanceMeter parses a string record and returns a DistanceMeter object.
// It extracts the elapsed time and distance from the record string and parses them into their respective data types.
func ParseDistanceMeter(record string) (*DistanceMeter, error) {
	var err error

	infos := strings.Fields(record)
	if len(infos) <= 1 {
		return nil, errors.New("invalid record info")
	}

	elapstimeStr, distanceStr := infos[0], infos[1]

	newElapsedTime, err := ParseElapsTime(elapstimeStr)
	if err != nil {
		return nil, errors.New("unable to parse elapsed time. reason: " + err.Error())
	}

	newDistance, err := parseDistance(distanceStr)
	if err != nil {
		return nil, errors.New("unable to parse distance. reason: " + err.Error())
	}

	return &DistanceMeter{
		elapsedTime: *newElapsedTime,
		distance:    newDistance,
	}, nil

}

// parseDistance parses a string representing distance and returns its value in float64 format.
func parseDistance(distStr string) (float64, error) {
	dist, err := strconv.ParseFloat(distStr, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid distance format: %s", distStr)
	}
	if dist < 0 || dist > 99999999.9 {
		return 0, fmt.Errorf("distance out of range: %s", distStr)
	}
	return dist, nil
}
