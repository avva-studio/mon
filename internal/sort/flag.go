package sort

import (
	"fmt"
	"strings"
)

type key struct {
	*string
}

// NewKey creates a key for use as a pflag.Value
func NewKey() *key {
	return &key{}
}

// String provides the string that represents the value the flag has been set
// to
func (f key) String() string {
	if f.string == nil {
		return ""
	}
	return *f.string
}

// Type returns the string that represents the type of key.
func (key) Type() string {
	return "sort-key"
}

// Set returns an error if the given value is not a suppported sort key and
// sets the value of the sort key
func (f *key) Set(value string) error {
	val := strings.TrimSpace(strings.ToLower(value))
	if !keyExists(AllKeys(), val) {
		return fmt.Errorf("unsupported sort key: %+v", value)
	}
	*f = key{string: &val}
	return nil
}

func keyExists(valids []string, key string) bool {
	for _, valid := range valids {
		if valid == key {
			return true
		}
	}
	return false
}
