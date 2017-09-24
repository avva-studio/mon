package main

import "strconv"

func parseIdString(idString string) (uint, error) {
	id, err := strconv.ParseUint(idString, 10, 64)
	return uint(id), err
}