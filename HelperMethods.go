package main

import "strconv"

func parseIDString(idString string) (uint, error) {
	id, err := strconv.ParseUint(idString, 10, 64)
	return uint(id), err
}
