package main

import "time"

const urlDateFormat string = `2006-01-02`

// urlFormatDateString formats a date into the accepted date for use for endpoints
func urlFormatDateString(date time.Time) string {
	return date.Format(urlDateFormat)
}

// parseDateString parses a date string in the expected format for use with this package
func parseDateString(dateString string) (time.Time, error) {
	return time.Parse(urlDateFormat, dateString)
}
