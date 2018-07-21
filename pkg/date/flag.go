package date

import (
	"strings"
	"time"

	"github.com/pkg/errors"
)

const dateFormat = "2006-01-02"

type flag struct {
	*time.Time
}

// Flag will create a Flag for use as a pflag.Value
func Flag() *flag {
	return &flag{}
}

// String will return the date in yyyy-mm-dd format, or an empty string if one
// has not been set.
func (f flag) String() string {
	if f.Time == nil {
		return ""
	}
	return f.Time.Format(dateFormat)
}

// Type returns the string that represents the type of flag.
func (flag) Type() string {
	return "date"
}

// Set parses the given string, attempting to create a logical date from its
// content. Set will match any supported date format or case insensitively
// match 'y' or 'yesterday'.
func (f *flag) Set(value string) error {
	val := strings.TrimSpace(value)
	if val == "" {
		return errors.New("no value given")
	}
	val = strings.ToLower(value)
	switch val {
	case "yesterday", "y":
		y := time.Now().Add(-time.Hour * 24)
		*f = flag{Time: &y}
		return nil
	}
	d, err := time.Parse(dateFormat, val)
	if err == nil {
		*f = flag{&d}
	}
	return errors.Wrapf(err, "unsupported date value: %+v", value)
}
