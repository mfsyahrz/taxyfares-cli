package taxyfare

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

// elapsTime represents a time duration in hours, minutes, seconds and nanoseconds.
type ElapsTime struct {
	Hours     time.Duration
	Minutes   time.Duration
	Seconds   time.Duration
	Remaining time.Duration
}

// Before returns true if the elapsTime value is before the input value.
func (e *ElapsTime) Before(input ElapsTime) bool {
	hourDiff := e.Hours - input.Hours
	if hourDiff > time.Hour {
		return false
	}

	minDiff := e.Minutes - input.Minutes
	if minDiff > time.Minute {
		return false
	}

	secDiff := e.Seconds - input.Seconds
	if secDiff > 0 {
		return false
	}

	remDiff := time.Duration(e.Remaining.Nanoseconds()) - time.Duration(input.Remaining.Nanoseconds())
	return remDiff < 0
}

// Compare returns the time difference between two elapsTime values in hours.
func (e ElapsTime) Compare(input ElapsTime) time.Duration {
	hourDiff := e.Hours - input.Hours
	minDiff := e.Minutes - input.Minutes
	secDiff := e.Seconds - input.Seconds

	return hourDiff + minDiff + secDiff
}

// String returns a formatted string representing the elapsTime value.
func (e ElapsTime) String() string {
	return fmt.Sprintf("%s:%s:%s.%d",
		formatDuration(int(e.Hours.Hours())),
		formatDuration(int(e.Minutes.Minutes())),
		formatDuration(int(e.Seconds.Seconds())),
		e.Remaining,
	)
}

// parseElapsTime parses a string representation of elapsTime into an elapsTime value.
// The string format should be "HH:MM:SS.NNN".
func ParseElapsTime(elapstimeStr string) (*ElapsTime, error) {
	timeInfos := strings.Split(elapstimeStr, ":")

	if len(timeInfos) != 3 {
		return nil, errors.New("invalid elapsed time input format")
	}

	hour, err := time.ParseDuration(fmt.Sprintf("%sh", timeInfos[0]))
	if err != nil {
		return nil, errors.New("parsing hours for elapsed time error: " + err.Error())
	}

	min, err := time.ParseDuration(fmt.Sprintf("%sm", timeInfos[1]))
	if err != nil {
		return nil, errors.New("parsing minutes for elapsed time error: " + err.Error())
	}

	secAndRemaining := strings.Split(timeInfos[2], ".")

	if len(secAndRemaining) != 2 {
		return nil, errors.New("invalid elapsed time seconds input")
	}

	sec, err := time.ParseDuration(fmt.Sprintf("%ss", secAndRemaining[0]))
	if err != nil {
		return nil, errors.New("parsing seconds for elapsed time error: " + err.Error())
	}

	rem, err := time.ParseDuration(fmt.Sprintf("%sns", secAndRemaining[1]))
	if err != nil {
		return nil, errors.New("parsing remaaining for elapsed time error: " + err.Error())
	}

	elapsTime := &ElapsTime{
		Hours:     hour,
		Minutes:   min,
		Seconds:   sec,
		Remaining: rem,
	}

	if elapsTime.Hours > 99*time.Hour || elapsTime.Minutes > 99*time.Minute || elapsTime.Seconds > 99*time.Second || elapsTime.Remaining > 999 {
		return nil, fmt.Errorf("invalid record info, elapsed time out of range %s", elapsTime)
	}

	return elapsTime, nil
}

func formatDuration(dur int) string {
	if dur > 9 {
		return fmt.Sprintf("%d", dur)
	}

	return fmt.Sprintf("0%d", dur)
}
