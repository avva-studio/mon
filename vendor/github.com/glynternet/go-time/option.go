package time

import "time"

// Option is a function that acts on a pointer to a Range and returns an error
type Option func(*Range) error

// Start returns an Option that will change the start time of a Range to Valid
// and set it to the given time.
func Start(t time.Time) Option {
	return func(r *Range) error {
		if r.end.Valid && r.end.Time.Before(t) {
			return ErrEndTimeBeforeStartTime
		}
		r.start = NullTime{Valid:true, Time:t}
		return nil
	}
}

// End returns an Option that will change the end time of a Range to Valid
// and set it to the given time.
func End(t time.Time) Option {
	return func(r *Range) error {
		if r.start.Valid && r.start.Time.After(t) {
			return ErrEndTimeBeforeStartTime
		}
		r.end = NullTime{Valid:true, Time:t}
		return nil
	}
}

