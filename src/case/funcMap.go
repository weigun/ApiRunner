package main

import (
	"strconv"
)

func randUser() string {

	x := strconv.FormatInt(random(99999999999), 10)
	return x

}

func randRange(min int64, max int64) string {
	if min >= max {
		return strconv.FormatInt(max, 10)
	}
	return strconv.FormatInt(random(max-min)+min, 10)
}
