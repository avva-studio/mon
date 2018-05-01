package date

import (
	"strings"
	"time"

	"github.com/pkg/errors"
)

const dateFormat = "2006-01-02"

type Value interface {
	String() string
	Set(string) error
	Type() string
}

type flag time.Time

func (f flag) String() string {
	return time.Time(f).Format(dateFormat)
}

func (flag) Type() string {
	return "date"
}

func (f *flag) Set(value string) error {
	val := strings.TrimSpace(value)
	if val == "" {
		return errors.New("no value given")
	}
	val = strings.ToLower(value)
	switch val {
	case "yesterday", "y":
		*f = flag(time.Now().Add(-time.Hour * 24))
		return nil
	}
	d, err := time.Parse(dateFormat, val)
	if err == nil {
		*f = flag(d)
	}
	return errors.Wrapf(err, "unsupported date value: %+v", value)
}
