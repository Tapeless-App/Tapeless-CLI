package util

import (
	"time"
)

func FormatDate(input string) string {

	if len(input) == 0 {
		return ""
	}

	date, err := time.Parse("2006-01-02T15:04:05.000Z", input)

	if err != nil {
		return ""
	}

	return date.Format("2006-01-02")

}
