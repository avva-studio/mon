package time

import (
	"time"
)

// New return a pointer to a new Range.
// If any errors occur from the options, New will return a nil pointer along
// with the error.
func New(options ...Option) (*Range, error) {
	r := new(Range)
	for _, o := range options {
		err := o(r)
		if err != nil {
			return nil, err
		}
	}
	return r, nil
}

// Range represents a range of time that can be open ended at none, either or both ends.
type Range struct {
	start NullTime
	end   NullTime
}

// Start returns the start NullTime of the Range
func (r Range) Start() NullTime {
	return r.start
}

// End returns the end NullTime of the Range
func (r Range) End() NullTime {
	return r.end
}

// Equal returns true if two Range objects have matching start and end NullTimes
func (r Range) Equal(r2 Range) bool {
	if !r.start.Equal(r2.start) || !r.end.Equal(r2.end) {
		return false
	}
	return true
}

// Contains returns true if the Range contains given time.
// Contains will always return true when both the start time and end time are not Valid
// Contains returns true if the time is on or after the Range's start time and before the Range's end time.
func (r Range) Contains(time time.Time) bool {
	if r.start.Valid && time.Before(r.start.Time) {
		return false
	}
	if r.end.Valid && time.After(r.end.Time) {
		return false
	}
	return true
}

// rangeValidationError holds an error describing an issue with a Range object
type rangeValidationError string

// Error ensures that rangeValidationError adheres to the error interface
func (err rangeValidationError) Error() string {
	return string(err)
}

const (
	// ErrEndTimeBeforeStartTime is a rangeValidationError that is returned when a Balance's closing date is set to before the opening date.
	ErrEndTimeBeforeStartTime = rangeValidationError("end time is before start time.")
)
